package tcln

import (
	"bytes"
	"errors"
	"testing"
	"time"

	bmock "github.com/cmd-stream/base-go/testdata/mock"
	"github.com/cmd-stream/delegate-go"
)

const Delta = 100 * time.Millisecond

func TestTransport(t *testing.T) {

	t.Run("ReceiveServerInfo should decode info from MUS encoding",
		func(t *testing.T) {
			var (
				wantInfo = []byte("info")
				bs       = func() []byte {
					bs := make([]byte, 0, delegate.ServerInfoMUS.Size(wantInfo))
					buf := bytes.NewBuffer(bs)
					delegate.ServerInfoMUS.Marshal(wantInfo, buf)
					return buf.Bytes()
				}()
				conn = bmock.NewConn().RegisterRead(
					func(b []byte) (n int, err error) {
						n = copy(b, bs)
						return
					},
				)
				transport = New[any](conn, nil)
				info, err = transport.ReceiveServerInfo()
			)
			if !bytes.Equal(info, wantInfo) {
				t.Errorf("unexpected info, want '%v' actual '%v'", wantInfo, info)
			}
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
		})

	t.Run("If decoding fails with an error, ReceiveServerInfo should return this error",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Read error")
				conn    = bmock.NewConn().RegisterRead(
					func(b []byte) (n int, err error) {
						return 0, wantErr
					},
				)
				transport = New[any](conn, nil)
				info, err = transport.ReceiveServerInfo()
			)
			if info != nil {
				t.Errorf("unexpected info, want '%v' actual '%v'", nil, info)
			}
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
		})

}

func SameTime(t1, t2 time.Time) bool {
	return !(t1.Before(t2.Truncate(Delta)) || t1.After(t2.Add(Delta)))
}
