package server

import (
	"bufio"
	"net"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/delegate-go"
	"github.com/cmd-stream/transport-go"
)

// New creates a new Transport.
func New[T any](conn net.Conn, codec transport.Codec[core.Result, core.Cmd[T]],
	opts ...transport.SetOption,
) *Transport[T] {
	o := transport.Options{}
	transport.Apply(opts, &o)
	var (
		w = bufio.NewWriterSize(conn, o.WriterBufSize)
		r = bufio.NewReaderSize(conn, o.ReaderBufSize)
	)
	return &Transport[T]{transport.New(conn, w, r, codec), w}
}

// Transport implements the delegate.ServerTransport interface.
type Transport[T any] struct {
	*transport.Transport[core.Result, core.Cmd[T]]
	w transport.Writer
}

func (t *Transport[T]) SendServerInfo(info delegate.ServerInfo) (
	err error,
) {
	_, err = delegate.ServerInfoMUS.Marshal(info, t.w)
	if err != nil {
		return
	}
	return t.Flush()
}

func (t *Transport[T]) WriterBufSize() int {
	return t.Transport.W.(*bufio.Writer).Size()
}

func (t *Transport[T]) ReaderBufSize() int {
	return t.Transport.R.(*bufio.Reader).Size()
}
