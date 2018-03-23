package main

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func getNextFreePort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(fmt.Errorf("cannot obtain free port: %v", err))
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func takePort(ctx context.Context, port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	go func() {
		defer listener.Close()
		<-ctx.Done()
	}()
	return nil
}

func takePortAfter(ctx context.Context, t time.Duration, port int, test *testing.T) error {
	go func() {
		select {
		case <-ctx.Done():
			test.Errorf("context expired before taking the port")
			return
		case <-time.After(t):
		}
		if err := takePort(ctx, port); err != nil {
			test.Errorf("Failed to take port")
		}
	}()
	return nil
}

func getPreambleFunc(state string, duration time.Duration) (preFn func(context.Context, int, *testing.T) error) {
	if state == "inuse" {
		preFn = func(ctx context.Context, port int, t *testing.T) error {
			return takePortAfter(ctx, duration, port, t)
		}
	} else {
		preFn = func(ctx context.Context, port int, t *testing.T) error {
			// go vet doesn't let me not call cancel and let the timeout expire
			// because if just looks for all paths calling cancel, not that WithTimeout
			// will do it for me, so with reimplement it
			// innerCtx, cancel := context.WithTimeout(ctx, duration)
			innerCtx, cancel := context.WithCancel(ctx)
			if err := takePort(innerCtx, port); err != nil {
				cancel()
				return err
			}
			time.AfterFunc(duration, cancel)
			return nil
		}
	}
	return preFn
}

func getCliArguments(port int, state string, host string, timeout int) []string {
	return []string{
		"--state", state,
		"--host", host,
		"--timeout", fmt.Sprintf("%d", timeout),
		fmt.Sprintf("%d", port),
	}
}

func testViaStruct(port int, state string, host string, timeout int, t *testing.T) error {
	cmd := NewWaitForPortCmd()
	cmd.State = state
	cmd.Host = host
	cmd.Timeout = timeout
	cmd.Args.Port = port

	return cmd.Execute([]string{})
}

func testViaCli(port int, state string, host string, timeout int, t *testing.T) error {
	cliArgs := getCliArguments(port, state, host, timeout)
	res := RunTool(cliArgs...)
	if !res.Success() {
		return fmt.Errorf("%s", res.Stderr())
	}
	return nil
}
