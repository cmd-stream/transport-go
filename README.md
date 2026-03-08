# transport-go

[![Go Reference](https://pkg.go.dev/badge/github.com/cmd-stream/transport-go.svg)](https://pkg.go.dev/github.com/cmd-stream/transport-go)
[![GoReportCard](https://goreportcard.com/badge/cmd-stream/transport-go)](https://goreportcard.com/report/github.com/cmd-stream/transport-go)
[![codecov](https://codecov.io/gh/cmd-stream/transport-go/graph/badge.svg?token=6JVVHR8QHF)](https://codecov.io/gh/cmd-stream/transport-go)

**transport-go** provides transport abstractions required to deliver Commands 
and Results between a `cmd-stream` client and server.

It includes implementations of the `client.Transport` and `server.Transport` 
interfaces defined in the `delegate-go` module.

The package uses `bufio.Writer` and `bufio.Reader` for efficient buffered I/O,
and relies on a user-defined codec to encode/decode `core.Cmd` or `core.Result`.
