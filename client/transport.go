package client

import (
	"bufio"
	"net"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/delegate-go"
	"github.com/cmd-stream/transport-go"
	"github.com/cmd-stream/transport-go/common"
)

// New creates a new Transport.
func New[T any](conf common.Conf, conn net.Conn,
	codec transport.Codec[base.Cmd[T], base.Result]) *Transport[T] {
	var (
		w = bufio.NewWriterSize(conn, conf.WriterBufSize)
		r = bufio.NewReaderSize(conn, conf.ReaderBufSize)
	)
	return &Transport[T]{
		r:         r,
		Transport: common.New(conn, w, r, codec),
	}
}

// Transport is an implementation of the delegate.ClientTransport interface.
//
// It will not send a Command if it is too large for the server.
type Transport[T any] struct {
	r transport.Reader
	*common.Transport[base.Cmd[T], base.Result]
}

func (t *Transport[T]) ReceiveServerInfo() (info delegate.ServerInfo, err error) {
	info, _, err = delegate.UnmarshalServerInfoMUS(t.r)
	return
}
