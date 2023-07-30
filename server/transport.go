package server

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
	codec transport.Codec[base.Result, base.Cmd[T]]) *Transport[T] {
	var (
		w = bufio.NewWriterSize(conn, conf.WriterBufSize)
		r = bufio.NewReaderSize(conn, conf.ReaderBufSize)
	)
	return &Transport[T]{w, common.New(conn, w, r, codec)}
}

// Transport is an implementation of the delegate.ServerTransport.
type Transport[T any] struct {
	w transport.Writer
	*common.Transport[base.Result, base.Cmd[T]]
}

func (t *Transport[T]) SendServerInfo(info delegate.ServerInfo) (
	err error) {
	_, err = delegate.MarshalServerInfoMUS(info, t.w)
	if err != nil {
		return
	}
	return t.Flush()
}

func (t *Transport[T]) SendServerSettings(settings delegate.ServerSettings) (
	err error) {
	_, err = delegate.MarshalServerSettingsMUS(settings, t.w)
	if err != nil {
		return
	}
	return t.Flush()
}
