package server

import (
	"bytes"
	"errors"
	"testing"

	cmock "github.com/cmd-stream/core-go/test/mock"
	"github.com/cmd-stream/delegate-go"
	mock "github.com/cmd-stream/transport-go/test/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

func TestTransport(t *testing.T) {
	t.Run("SendServerInfo should encode info to MUS encoding",
		func(t *testing.T) {
			var (
				wantInfo delegate.ServerInfo = []byte("info")
				wantBs                       = infoToBs(wantInfo)
				wantErr  error               = nil
				conn                         = cmock.NewConn().RegisterWrite(
					func(bs []byte) (n int, err error) {
						asserterror.EqualDeep(t, bs, wantBs)
						n = len(bs)
						return
					},
				)
				transport = New[any](conn, nil)
				err       = transport.SendServerInfo(wantInfo)
			)
			asserterror.EqualError(t, err, wantErr)
		})

	t.Run("If Conn.Write fails with an error, SendServerInfo should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Conn.Write error")
				conn    = cmock.NewConn().RegisterWrite(
					func(b []byte) (n int, err error) {
						err = wantErr
						return
					},
				)
				transport = New[any](conn, nil)
				err       = transport.SendServerInfo(nil)
			)
			asserterror.EqualError(t, err, wantErr)
		})

	t.Run("If MarshalServerInfo fails with an error, SendServerInfo should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("WriteByte error")
				writer  = mock.NewWriter().RegisterWriteByte(
					func(b byte) error { return wantErr },
				)
				mocks     = []*mok.Mock{writer.Mock}
				transport = &Transport[any]{w: writer}
				err       = transport.SendServerInfo((delegate.ServerInfo([]byte{})))
			)
			asserterror.EqualError(t, err, wantErr)
			asserterror.EqualDeep(t, mok.CheckCalls(mocks), mok.EmptyInfomap)
		})
}

func infoToBs(info delegate.ServerInfo) []byte {
	var (
		size = delegate.ServerInfoMUS.Size(info)
		bs   = make([]byte, 0, size)
		buf  = bytes.NewBuffer(bs)
		_, _ = delegate.ServerInfoMUS.Marshal(info, buf)
	)
	return buf.Bytes()
}
