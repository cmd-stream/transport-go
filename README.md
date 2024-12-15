# transport-go

[![Go Reference](https://pkg.go.dev/badge/github.com/cmd-stream/transport-go.svg)](https://pkg.go.dev/github.com/cmd-stream/transport-go)
[![GoReportCard](https://goreportcard.com/badge/cmd-stream/transport-go)](https://goreportcard.com/report/github.com/cmd-stream/transport-go)
[![codecov](https://codecov.io/gh/cmd-stream/transport-go/graph/badge.svg?token=6JVVHR8QHF)](https://codecov.io/gh/cmd-stream/transport-go)

transport-go provides commands/results delivery for the cmd-stream-go client 
and server.

It contains implementations of the `delegate.ClientTransport` and 
`delegate.ServerTransport` interfaces.

A key feature of this module is the use of `bufio.Writer` and `bufio.Reader` 
along with a user-defined Codec to convert raw bytes into `base.Cmd` or 
`base.Result`.
