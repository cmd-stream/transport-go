package tser

import (
	"bufio"
	"net"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/delegate-go"
	"github.com/cmd-stream/transport-go"
	tcom "github.com/cmd-stream/transport-go/common"
)

// New creates a new Transport.
func New[T any](conn net.Conn, codec transport.Codec[base.Result, base.Cmd[T]],
	ops ...tcom.SetOption) *Transport[T] {
	options := tcom.Options{}
	tcom.Apply(ops, &options)
	var (
		w = bufio.NewWriterSize(conn, options.WriterBufSize)
		r = bufio.NewReaderSize(conn, options.ReaderBufSize)
	)
	return &Transport[T]{tcom.New(conn, w, r, codec), w}
}

// Transport implements the delegate.ServerTransport interface.
type Transport[T any] struct {
	*tcom.Transport[base.Result, base.Cmd[T]]
	w transport.Writer
}

func (t *Transport[T]) SendServerInfo(info delegate.ServerInfo) (
	err error) {
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
