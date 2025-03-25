package tcom

import (
	"net"
	"time"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
)

// New creates a new Transport.
func New[T, V any](conn net.Conn, w transport.Writer, r transport.Reader,
	codec transport.Codec[T, V]) *Transport[T, V] {
	return &Transport[T, V]{w, r, conn, codec}
}

// Transport is a common transport for both client and server.
//
// It uses a user-defined codec to encode and decode data over the connection.
type Transport[T, V any] struct {
	W     transport.Writer
	R     transport.Reader
	conn  net.Conn
	codec transport.Codec[T, V]
}

// LocalAddr returns the connection local network address.
func (tn *Transport[T, V]) LocalAddr() net.Addr {
	return tn.conn.LocalAddr()
}

// RemoteAddr returns the connection remote network address.
func (tn *Transport[T, V]) RemoteAddr() net.Addr {
	return tn.conn.RemoteAddr()
}

// SetSendDeadline sets a send deadline.
func (tn *Transport[T, V]) SetSendDeadline(deadline time.Time) error {
	return tn.conn.SetWriteDeadline(deadline)
}

// Send sends data using the codec.
func (tn *Transport[T, V]) Send(seq base.Seq, t T) (err error) {
	return tn.codec.Encode(seq, t, tn.W)
}

// Flush flushes any buffered data.
func (tn *Transport[T, V]) Flush() (err error) {
	return tn.W.Flush()
}

// SetReceiveDeadline sets a receive deadline.
func (tn *Transport[T, V]) SetReceiveDeadline(deadline time.Time) error {
	return tn.conn.SetReadDeadline(deadline)
}

// Receive receives data using the codec.
func (tn *Transport[T, V]) Receive() (seq base.Seq, v V, err error) {
	return tn.codec.Decode(tn.R)
}

// Close closes the underlying connection.
func (tn *Transport[T, V]) Close() error {
	return tn.conn.Close()
}
