package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

const (
	// PortInUse defines the state of a port in use
	PortInUse = "inuse"
	// PortFree defines the state of a free port
	PortFree = "free"
)

// WaitForPortCmd allows checking a port state
type WaitForPortCmd struct {
	Host    string `short:"h" long:"host" description:"Host where to check for the port" default:"" value-name:"HOST"`
	State   string `short:"s" long:"state" choice:"inuse" choice:"free" description:"State to wait for" default:"inuse"`
	Timeout int    `short:"t" long:"timeout" default:"30" description:"Timeout in seconds to wait for the port" value-name:"SECONDS"`
	Args    struct {
		Port int `positional-arg-name:"port"`
	} `positional-args:"yes" required:"yes"`
}

// NewWaitForPortCmd returns a WaitForPortCmd with configured defaults
func NewWaitForPortCmd() *WaitForPortCmd {
	return &WaitForPortCmd{
		State:   PortFree,
		Host:    "",
		Timeout: 30,
	}
}

// Execute performs the port check
func (c *WaitForPortCmd) Execute(args []string) error {
	var checkPortState func(ctx context.Context, host string, port int) bool
	switch c.State {
	case PortInUse:
		checkPortState = portIsInUse
	case PortFree:
		checkPortState = func(ctx context.Context, host string, port int) bool {
			return !portIsInUse(ctx, host, port)
		}
	default:
		return fmt.Errorf("unknown state %q", c.State)
	}
	if err := validatePort(c.Args.Port); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
	defer cancel()
	if err := validateHost(ctx, c.Host); err != nil {
		return err
	}

	for !checkPortState(ctx, c.Host, c.Args.Port) {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout reached before the port went into state %q", c.State)
		case <-time.After(500 * time.Millisecond):
		}
	}
	return nil
}

func validatePort(port int) error {
	if port <= 0 {
		return fmt.Errorf("port out of range: port must be greater than zero")
	} else if port > 65535 {
		return fmt.Errorf("port out of range: port must be <= 65535")
	}
	return nil
}

func validateHost(ctx context.Context, host string) error {
	// An empty host is perfectly fine for us but net.LookupHost will fail
	if host == "" {
		return nil
	}
	if _, err := net.DefaultResolver.LookupHost(ctx, host); err != nil {
		return fmt.Errorf("cannot resolve host %q: %v", host, err)
	}
	return nil
}

func isAddrInUseError(err error) bool {
	if err, ok := err.(*net.OpError); ok {
		if err, ok := err.Err.(*os.SyscallError); ok {
			return err.Err == syscall.EADDRINUSE
		}
	}
	return false
}

func canConnectToPort(ctx context.Context, host string, port int) bool {
	d := net.Dialer{Timeout: 60 * time.Second}
	conn, err := d.DialContext(ctx, "tcp", net.JoinHostPort(host, fmt.Sprintf("%d", port)))
	if err == nil {
		defer conn.Close()
		return true
	}
	return false
}

// portIsInUse allows checking if a port is in use in the specified host.
func portIsInUse(ctx context.Context, host string, port int) bool {
	// If we can connect, is in use
	if canConnectToPort(ctx, host, port) {
		return true
	}

	// If we are trying to check a remote host, we cannot do more, so we consider it not in use
	if host != "" {
		return false
	}

	// If we are checking locally, try to listen
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err == nil {
		listener.Close()
		return false
	} else if isAddrInUseError(err) {
		return true
	}
	// We could not connect to the port, and we cannot listen on it, the safest thing
	// we can assume in localhost is that is not in use (binding to a privileged port, for example)
	return false
}
