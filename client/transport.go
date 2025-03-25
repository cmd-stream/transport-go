package tcln

import (
	"bufio"
	"net"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/delegate-go"
	"github.com/cmd-stream/transport-go"
	tcom "github.com/cmd-stream/transport-go/common"
)

// New creates a new Transport.
func New[T any](conn net.Conn, codec transport.Codec[base.Cmd[T], base.Result],
	ops ...tcom.SetOption) *Transport[T] {
	options := tcom.Options{}
	tcom.Apply(ops, &options)
	var (
		w = bufio.NewWriterSize(conn, options.WriterBufSize)
		r = bufio.NewReaderSize(conn, options.ReaderBufSize)
	)
	return &Transport[T]{tcom.New(conn, w, r, codec)}
}

// Transport implements the delegate.ClientTransport interface.
type Transport[T any] struct {
	*tcom.Transport[base.Cmd[T], base.Result]
}

func (t *Transport[T]) ReceiveServerInfo() (info delegate.ServerInfo,
	err error) {
	info, _, err = delegate.ServerInfoMUS.Unmarshal(t.Transport.R)
	return
}

func (t *Transport[T]) WriterBufSize() int {
	return t.Transport.W.(*bufio.Writer).Size()
}

func (t *Transport[T]) ReaderBufSize() int {
	return t.Transport.R.(*bufio.Reader).Size()
}
