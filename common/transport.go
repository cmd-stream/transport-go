package common

import (
	"net"
	"time"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
)

// New creates a new Transport.
func New[T, V any](conn net.Conn, w transport.Writer, r transport.Reader,
	codec transport.Codec[T, V]) *Transport[T, V] {
	return &Transport[T, V]{conn, codec, w, r}
}

// Transport is a common transport for the client and server delegates.
//
// It uses a user-defined codec to encode/decode data to/from the connection.
type Transport[T, V any] struct {
	conn  net.Conn
	codec transport.Codec[T, V]
	w     transport.Writer
	r     transport.Reader
}

// LocalAddr returns the connection local network address.
func (tr *Transport[T, V]) LocalAddr() net.Addr {
	return tr.conn.LocalAddr()
}

// RemoteAddr returns the connection remote network address.
func (tr *Transport[T, V]) RemoteAddr() net.Addr {
	return tr.conn.RemoteAddr()
}

// SetSendDeadline sets a send deadline.
func (tr *Transport[T, V]) SetSendDeadline(deadline time.Time) error {
	return tr.conn.SetWriteDeadline(deadline)
}

// Send sends data with the associated sequence number using the codec.
func (tr *Transport[T, V]) Send(seq base.Seq, t T) (err error) {
	return tr.codec.Encode(seq, t, tr.w)
}

// Flush flushes any buffered data.
func (tr *Transport[T, V]) Flush() (err error) {
	return tr.w.Flush()
}

// SetReceiveDeadline sets a receive deadline.
func (tr *Transport[T, V]) SetReceiveDeadline(deadline time.Time) error {
	return tr.conn.SetReadDeadline(deadline)
}

// Receives receives data with the associated sequence number using the codec.
func (tr *Transport[T, V]) Receive() (seq base.Seq, v V, err error) {
	return tr.codec.Decode(tr.r)
}

// Close closes the underlying connection.
func (tr *Transport[T, V]) Close() error {
	return tr.conn.Close()
}
