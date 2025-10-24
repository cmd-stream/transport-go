# transport-go

[![Go Reference](https://pkg.go.dev/badge/github.com/cmd-stream/transport-go.svg)](https://pkg.go.dev/github.com/cmd-stream/transport-go)
[![GoReportCard](https://goreportcard.com/badge/cmd-stream/transport-go)](https://goreportcard.com/report/github.com/cmd-stream/transport-go)
[![codecov](https://codecov.io/gh/cmd-stream/transport-go/graph/badge.svg?token=6JVVHR8QHF)](https://codecov.io/gh/cmd-stream/transport-go)

**transport-go** provides transport abstractions for delivering Commands and
Results between `cmd-stream-go` clients and servers.

It includes implementations of the `delegate.ClientTransport` and
`delegate.ServerTransport` interfaces.

The package uses `bufio.Writer` and `bufio.Reader` for efficient buffered I/O,
and relies on a user-defined codec to serialize and deserialize values into
`core.Cmd` or `core.Result`.
