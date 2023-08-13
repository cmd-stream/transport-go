# transport-go
transport-go provides commands/results delivery for the cmd-stream-go client 
and server.

It contains implementations of the `delegate.ClientTransport` and 
`delegate.ServerTransport` interfaces (they are located in the corresponding 
packages).

A feature of this module is the use of `bufio.Writer`, `bufio.Reader`, and 
user-defined `Codec` to convert raw bytes to `base.Cmd` or `base.Result`.

# Tests
Test coverage is about 95%.