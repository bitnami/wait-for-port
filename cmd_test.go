package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testFn func(port int, state string, host string, timeout int, t *testing.T) error

func TestWaitForPortCmd_ExecuteValidations(t *testing.T) {
	validPort := 1000
	for _, tt := range []struct {
		name        string
		port        int
		state       string
		host        string
		timeout     int
		expectedErr interface{}
	}{
		{
			name: "Detects invalid port (too low)", port: -1000, state: PortFree, timeout: 10,
			expectedErr: "port out of range: port must be greater than zero",
		},
		{
			name: "Detects invalid port (too low)", port: 0, state: PortFree, timeout: 10,
			expectedErr: "port out of range: port must be greater than zero",
		},
		{
			name: "Detects invalid port (too high)", port: 65536, state: PortFree, timeout: 10,
			expectedErr: "port out of range: port must be <= 65535",
		},
		{
			name: "Detects invalid host (nonexistent)", host: "someveryrandomhostfrombitnami", port: validPort, state: PortFree, timeout: 10,
			expectedErr: `cannot resolve host "someveryrandomhostfrombitnami"`,
		},
		{
			name: "Detects invalid state", port: validPort, state: "", timeout: 10,
			expectedErr: "unknown state",
		},
		{
			name: "Detects invalid state", port: validPort, state: "foobar",
			expectedErr: "unknown state",
		},
	} {

		t.Run(tt.name, func(t *testing.T) {
			err := testViaStruct(tt.port, tt.state, tt.host, tt.timeout, t)
			if err == nil {
				t.Errorf("WaitForPortCmd.Execute() was expected to fail but succeeded")
			} else if tt.expectedErr != nil {
				assert.Regexp(t, tt.expectedErr, err, "Expected error %v to match %v", err, tt.expectedErr)
			}
		})
	}
}

func TestWaitForPortCmd_Execute_Timeout(t *testing.T) {

	freePort := getNextFreePort()
	takenPort := getNextFreePort()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := takePort(ctx, takenPort); err != nil {
		t.Fatalf("cannot take port %d: %v", takenPort, err)
	}

	maxTimeDelta := 1 * time.Second
	timeout := 2
	duration := time.Duration(timeout) * time.Second

	for _, spec := range []struct {
		testFn testFn
		suffix string
	}{
		{testFn: testViaStruct, suffix: "using struct"},
		{testFn: testViaCli, suffix: "using cli"},
	} {
		t.Run("Timeout waiting for port to be in use "+spec.suffix, func(t *testing.T) {
			start := time.Now()
			err := spec.testFn(freePort, PortInUse, "", timeout, t)
			end := time.Now()

			assert.WithinDuration(t, end, start.Add(duration), maxTimeDelta)

			if err == nil {
				t.Errorf("expected check to fail with timeout, but it succeeded")
			} else {
				assert.Regexp(t, `timeout reached before the port went into state "inuse"`, err)
			}
		})
		t.Run("Timeout waiting for port to be free "+spec.suffix, func(t *testing.T) {
			start := time.Now()
			err := spec.testFn(takenPort, PortFree, "", timeout, t)
			end := time.Now()

			assert.WithinDuration(t, end, start.Add(duration), maxTimeDelta)

			if err == nil {
				t.Errorf("expected check to fail with timeout, but it succeeded")
			} else {
				assert.Regexp(t, `timeout reached before the port went into state "free"`, err)
			}
		})
	}
}
func TestWaitForPortCmd_Execute(t *testing.T) {
	timeout := 10

	for _, sleep := range []int{0, 5} {
		for _, tt := range []struct {
			name  string
			sleep int
			state string
			host  string
		}{
			{
				name:  "Detect port was taken",
				state: "inuse",
			},
			{
				name:  "Detect port is freed",
				state: "free",
			},
		} {
			duration := time.Duration(sleep) * time.Second
			maxTimeDelta := time.Second
			preFn := getPreambleFunc(tt.state, duration)

			for _, spec := range []struct {
				testFn testFn
				suffix string
			}{
				{testFn: testViaStruct, suffix: "using struct"},
				{testFn: testViaCli, suffix: "using cli"},
			} {
				testFn := spec.testFn
				newTitle := fmt.Sprintf("%s %s", tt.name, spec.suffix)

				if sleep != 0 {
					newTitle = fmt.Sprintf("%s (after %d seconds)", newTitle, sleep)
				}

				t.Run(newTitle, func(t *testing.T) {
					port := getNextFreePort()
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					if err := preFn(ctx, port, t); err != nil {
						t.Fatalf("Cannot execute test preamble: %v", err)
						return
					}

					start := time.Now()
					if err := testFn(port, tt.state, tt.host, timeout, t); err != nil {
						t.Errorf("Expected command to succeed but got: %v", err)
					}
					end := time.Now()
					assert.WithinDuration(t, end, start.Add(duration), maxTimeDelta)
				})
			}
		}
	}
}

func Test_isInUse(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	freePort := getNextFreePort()
	takenPort := getNextFreePort()
	if err := takePort(ctx, takenPort); err != nil {
		t.Fatalf("cannot take port %d: %v", takenPort, err)
	}
	type args struct {
		host string

		port int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Free port detected (no host)",
			want: false,
			args: args{
				port: freePort,
			},
		},
		{
			name: "Free port detected (localhost)",
			want: false,
			args: args{
				port: freePort,
				host: "localhost",
			},
		},
		{
			name: "Taken port detected (no host)",
			want: true,
			args: args{
				port: takenPort,
			},
		},

		{
			name: "Taken port detected (localhost)",
			want: true,
			args: args{
				port: takenPort,
				host: "localhost",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := portIsInUse(context.Background(), tt.args.host, tt.args.port); got != tt.want {
				t.Errorf("portIsInUse() = %v, want %v", got, tt.want)
			}
		})
	}
}
