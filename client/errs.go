package client

import "errors"

// ErrTooBigCmd happens when the client tries to send a command that is too
// big for the server.
var ErrTooBigCmd = errors.New("too big command")
