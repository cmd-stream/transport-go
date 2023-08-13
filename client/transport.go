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
		codec:     codec,
	}
}

// Transport is an implementation of the delegate.ClientTransport interface.
//
// It will not send a command if it is too big for the server.
type Transport[T any] struct {
	r transport.Reader
	*common.Transport[base.Cmd[T], base.Result]
	codec    transport.Codec[base.Cmd[T], base.Result]
	settings delegate.ServerSettings
}

func (t *Transport[T]) ReceiveServerInfo() (info delegate.ServerInfo, err error) {
	info, _, err = delegate.UnmarshalServerInfoMUS(t.r)
	return
}

func (t *Transport[T]) ReceiveServerSettings() (settings delegate.ServerSettings,
	err error) {
	settings, _, err = delegate.UnmarshalServerSettingsMUS(t.r)
	return
}

func (t *Transport[T]) ApplyServerSettings(settings delegate.ServerSettings) {
	t.settings = settings
}

// Send sends a command.
//
// Returns ErrTooBigCmd, if the command size is bigger than
// ServerSettings.MaxCmdSize.
func (t *Transport[T]) Send(seq base.Seq, cmd base.Cmd[T]) error {
	if t.settings.MaxCmdSize > 0 && t.codec.Size(cmd) > t.settings.MaxCmdSize {
		return ErrTooBigCmd
	}
	return t.Transport.Send(seq, cmd)
}
