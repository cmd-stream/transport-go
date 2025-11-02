package client

import (
	"bytes"
	"errors"
	"testing"

	"github.com/cmd-stream/delegate-go"
	cmocks "github.com/cmd-stream/testkit-go/mocks/core"

	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestTransport(t *testing.T) {
	t.Run("ReceiveServerInfo should decode info from MUS encoding",
		func(t *testing.T) {
			var (
				wantInfo delegate.ServerInfo = []byte("info")
				wantErr  error               = nil
				bs                           = func() []byte {
					bs := make([]byte, 0, delegate.ServerInfoMUS.Size(wantInfo))
					buf := bytes.NewBuffer(bs)
					delegate.ServerInfoMUS.Marshal(wantInfo, buf)
					return buf.Bytes()
				}()
				conn = cmocks.NewConn().RegisterRead(
					func(b []byte) (n int, err error) {
						n = copy(b, bs)
						return
					},
				)
				transport = New[any](conn, nil)
				info, err = transport.ReceiveServerInfo()
			)
			asserterror.EqualDeep(info, wantInfo, t)
			asserterror.EqualError(err, wantErr, t)
		})

	t.Run("If decoding fails with an error, ReceiveServerInfo should return this error",
		func(t *testing.T) {
			var (
				wantInfo delegate.ServerInfo = nil
				wantErr                      = errors.New("Read error")
				conn                         = cmocks.NewConn().RegisterRead(
					func(b []byte) (n int, err error) {
						return 0, wantErr
					},
				)
				transport = New[any](conn, nil)
				info, err = transport.ReceiveServerInfo()
			)
			asserterror.EqualDeep(info, wantInfo, t)
			asserterror.EqualError(err, wantErr, t)
		})
}
