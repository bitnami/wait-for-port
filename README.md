[![Go Report Card](https://goreportcard.com/badge/github.com/bitnami/wait-for-port)](https://goreportcard.com/report/github.com/bitnami/wait-for-port)
[![CI](https://github.com/bitnami/wait-for-port/actions/workflows/main.yml/badge.svg)](https://github.com/bitnami/wait-for-port/actions/workflows/main.yml)

# wait-for-port

This tool allows waiting for a port to enter into the requested state (free or in use), with a customizable timeout

# Basic usage

```console
$ wait-for-port --help
Usage:
  wait-for-port [OPTIONS] port

Application Options:
  -h, --host=HOST             Host where to check for the port
  -s, --state=[inuse|free]    State to wait for (default: inuse)
  -t, --timeout=SECONDS       Timeout in seconds to wait for the port (default: 30)

Help Options:
  -h, --help                  Show this help message
```

# Examples

## Wait for a port to be in use

```console
$ wait-for-port --state=inuse 12345
$ echo $?
0
```

Or in a remote server:

```console
$ wait-for-port --host=myhost.example.com --state=inuse 12345
$ echo $?
0
```

## Wait for a port to be free

```console
$ wait-for-port --state=free 12345
$ echo $?
0
```

Or in a remote server:

```console
$ wait-for-port --host=myhost.example.com --state=free 12345
$ echo $?
0
```

## The tool times out before the port goes into the required state

If the port does not go into the required state under the provided timeout time, the process will return a non-zero exit code
so it is easily recognizable from a parent process:

```console
$ wait-for-port --timeout=10 --state=inuse 13456
timeout reached before the port went into state "inuse"
$ echo $?
1
```

```console
$ wait-for-port --timeout=10 --state=free 8080
timeout reached before the port went into state "free"
$ echo $?
1
```

## License

Copyright &copy; 2025 Broadcom. The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.

You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and limitations under the License.
