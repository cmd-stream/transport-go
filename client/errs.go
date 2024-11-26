package client

import "errors"

// ErrTooLargeCmd happens when the client tries to send a command that is too
// large for the server.
var ErrTooLargeCmd = errors.New("too large command")
