package client

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/cmd-stream/base-go"
	cs_mock "github.com/cmd-stream/base-go/testdata/mock"
	"github.com/cmd-stream/delegate-go"
	"github.com/cmd-stream/transport-go/common"
	transport_mock "github.com/cmd-stream/transport-go/testdata/mock"
)

const Delta = 100 * time.Millisecond

func TestTransport(t *testing.T) {

	t.Run("ReceiveServerInfo should decode info from MUS encoding",
		func(t *testing.T) {
			var (
				wantInfo = []byte("info")
				bs       = func() []byte {
					bs := make([]byte, 0, delegate.SizeServerInfoMUS(wantInfo))
					buf := bytes.NewBuffer(bs)
					delegate.MarshalServerInfoMUS(wantInfo, buf)
					return buf.Bytes()
				}()
				conn = cs_mock.NewConn().RegisterRead(
					func(b []byte) (n int, err error) {
						n = copy(b, bs)
						return
					},
				)
				transport = New[any](common.Conf{}, conn, nil)
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
				conn    = cs_mock.NewConn().RegisterRead(
					func(b []byte) (n int, err error) {
						return 0, wantErr
					},
				)
				transport = New[any](common.Conf{}, conn, nil)
				info, err = transport.ReceiveServerInfo()
			)
			if info != nil {
				t.Errorf("unexpected info, want '%v' actual '%v'", nil, info)
			}
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
		})

	t.Run("ReceiveServerSettings should decode settings from MUS encoding",
		func(t *testing.T) {
			var (
				wantSettings = delegate.ServerSettings{MaxCmdSize: 100}
				bs           = func() []byte {
					bs := make([]byte, 0, delegate.SizeServerSettingsMUS(wantSettings))
					buf := bytes.NewBuffer(bs)
					delegate.MarshalServerSettingsMUS(wantSettings, buf)
					return buf.Bytes()
				}()
				conn = cs_mock.NewConn().RegisterRead(
					func(b []byte) (n int, err error) {
						n = copy(b, bs)
						return
					},
				)
				transport     = New[any](common.Conf{}, conn, nil)
				settings, err = transport.ReceiveServerSettings()
			)
			if !reflect.DeepEqual(settings, wantSettings) {
				t.Errorf("unexpected info, want '%v' actual '%v'", wantSettings, settings)
			}
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
		})

	t.Run("If decoding fails with an error, ReceiveServerSettings should return this error",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Read error")
				conn    = cs_mock.NewConn().RegisterRead(
					func(b []byte) (n int, err error) {
						return 0, wantErr
					},
				)
				transport     = New[any](common.Conf{}, conn, nil)
				settings, err = transport.ReceiveServerSettings()
			)
			if !reflect.DeepEqual(settings, delegate.ServerSettings{}) {
				t.Errorf("unexpected info, want '%v' actual '%v'", nil, settings)
			}
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
		})

	t.Run("ApplyServerSettings should affect Send method",
		func(t *testing.T) {
			var (
				codec = transport_mock.NewClientCodec().RegisterSize(
					func(cmd base.Cmd[any]) (size int) { return 5 },
				)
				transport = New[any](common.Conf{}, nil, codec)
			)
			transport.ApplyServerSettings(delegate.ServerSettings{MaxCmdSize: 3})
			err := transport.Send(1, nil)
			if err != ErrTooLargeCmd {
				t.Errorf("unexpected error, want '%v' actual '%v'", ErrTooLargeCmd, err)
			}
		})

}

func SameTime(t1, t2 time.Time) bool {
	return !(t1.Before(t2.Truncate(Delta)) || t1.After(t2.Add(Delta)))
}
