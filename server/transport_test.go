package server

import (
	"bytes"
	"errors"
	"testing"
	"time"

	bmock "github.com/cmd-stream/base-go/testdata/mock"
	"github.com/cmd-stream/delegate-go"
	"github.com/cmd-stream/transport-go/common"
	"github.com/cmd-stream/transport-go/testdata/mock"
	"github.com/ymz-ncnk/mok"
)

const Delta = 100 * time.Millisecond

func TestTransport(t *testing.T) {

	t.Run("SendServerInfo should encode info to MUS encoding",
		func(t *testing.T) {
			var (
				wantInfo = []byte("info")
				wantBs   = func() []byte {
					bs := make([]byte, 0, delegate.SizeServerInfoMUS(wantInfo))
					buf := bytes.NewBuffer(bs)
					delegate.MarshalServerInfoMUS(wantInfo, buf)
					return buf.Bytes()
				}()
				conn = bmock.NewConn().RegisterWrite(
					func(bs []byte) (n int, err error) {
						if !bytes.Equal(bs, wantBs) {
							t.Errorf("unexpected bs, want '%v' actual '%v'", wantBs, bs)
						}
						n = len(bs)
						return
					},
				)
				transport = New[any](common.Conf{}, conn, nil)
				err       = transport.SendServerInfo(wantInfo)
			)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
		})

	t.Run("If Conn.Write fails with an error, SendServerInfo should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Conn.Write error")
				conn    = bmock.NewConn().RegisterWrite(
					func(b []byte) (n int, err error) {
						err = wantErr
						return
					},
				)
				transport = New[any](common.Conf{}, conn, nil)
				err       = transport.SendServerInfo(nil)
			)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
		})

	t.Run("SendServerSettings should encode settings to MUS encoding",
		func(t *testing.T) {
			var (
				wantSettings = delegate.ServerSettings{MaxCmdSize: 3}
				wantBs       = func() []byte {
					bs := make([]byte, 0, delegate.SizeServerSettingsMUS(wantSettings))
					buf := bytes.NewBuffer(bs)
					delegate.MarshalServerSettingsMUS(wantSettings, buf)
					return buf.Bytes()
				}()
				conn = bmock.NewConn().RegisterWrite(
					func(bs []byte) (n int, err error) {
						if !bytes.Equal(bs, wantBs) {
							t.Errorf("unexpected bs, want '%v' actual '%v'", wantBs, bs)
						}
						n = len(bs)
						return
					},
				)
				transport = New[any](common.Conf{}, conn, nil)
				err       = transport.SendServerSettings(wantSettings)
			)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
		})

	t.Run("If Conn.Write fails with an error, SendServerSettings should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Conn.Write error")
				conn    = bmock.NewConn().RegisterWrite(
					func(b []byte) (n int, err error) {
						err = wantErr
						return
					},
				)
				transport = New[any](common.Conf{}, conn, nil)
				err       = transport.SendServerSettings(delegate.ServerSettings{})
			)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
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
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
			if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
				t.Error(infomap)
			}
		})

	t.Run("If MarshalServerSettings fails with an error, SendServerSettings should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("WriteByte error")
				writer  = mock.NewWriter().RegisterWriteByte(
					func(b byte) error { return wantErr },
				)
				transport = &Transport[any]{w: writer}
				err       = transport.SendServerSettings(delegate.ServerSettings{})
			)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
		})

}
