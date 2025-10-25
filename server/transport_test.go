package server

import (
	"bytes"
	"errors"
	"testing"

	cmock "github.com/cmd-stream/core-go/testdata/mock"
	"github.com/cmd-stream/delegate-go"
	"github.com/cmd-stream/transport-go/testdata/mock"
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
						asserterror.EqualDeep(bs, wantBs, t)
						n = len(bs)
						return
					},
				)
				transport = New[any](conn, nil)
				err       = transport.SendServerInfo(wantInfo)
			)
			asserterror.EqualError(err, wantErr, t)
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
			asserterror.EqualError(err, wantErr, t)
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
			asserterror.EqualError(err, wantErr, t)
			asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
		})
}

func infoToBs(info delegate.ServerInfo) []byte {
	bs := make([]byte, 0, delegate.ServerInfoMUS.Size(info))
	buf := bytes.NewBuffer(bs)
	delegate.ServerInfoMUS.Marshal(info, buf)
	return buf.Bytes()
}
