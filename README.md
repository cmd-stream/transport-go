# transport-go
transport-go is an implementation of the transport for the `delegate-go` module.
That is, it contains implementations of the `delegate.ClientTransport` and 
`delegate.ServerTransport` interfaces (they are located in the corresponding 
packages).

A feature of this module is the use of the `bufio.Writer`, `bufio.Reader`, and a 
user-defined `Codec` to convert raw bytes to the `base.Cmd` or `base.Result`.

# Tests
Test coverage is about 95%.