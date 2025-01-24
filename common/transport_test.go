package common

import (
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/cmd-stream/base-go"
	bmock "github.com/cmd-stream/base-go/testdata/mock"
	"github.com/cmd-stream/transport-go"
	tmock "github.com/cmd-stream/transport-go/testdata/mock"
	"github.com/ymz-ncnk/mok"
)

const Delta = 100 * time.Millisecond

func TestTransport(t *testing.T) {

	t.Run("LocalAddr should return local address of the conn",
		func(t *testing.T) {
			var (
				wantAddr = &net.TCPAddr{IP: net.ParseIP("127.0.0.1:9000")}
				conn     = bmock.NewConn().RegisterLocalAddr(
					func() (addr net.Addr) {
						return wantAddr
					},
				)
				mocks     = []*mok.Mock{conn.Mock}
				transport = New[base.Cmd[any], base.Result](conn, nil, nil, nil)
			)
			addr := transport.LocalAddr()
			if addr != wantAddr {
				t.Errorf("unexpected addr, want '%v' actual '%v'", wantAddr, addr)
			}
			if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
				t.Error(infomap)
			}
		})

	t.Run("RemoteAddr should return remote address of the conn",
		func(t *testing.T) {
			var (
				wantAddr = &net.TCPAddr{IP: net.ParseIP("127.0.0.1:9000")}
				conn     = bmock.NewConn().RegisterRemoteAddr(
					func() (addr net.Addr) {
						return wantAddr
					},
				)
				mocks     = []*mok.Mock{conn.Mock}
				transport = New[base.Cmd[any], base.Result](conn, nil, nil, nil)
			)
			addr := transport.RemoteAddr()
			if addr != wantAddr {
				t.Errorf("unexpected addr, want '%v' actual '%v'", wantAddr, addr)
			}
			if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
				t.Error(infomap)
			}
		})

	t.Run("Conn.SetWriteDeadline should receive same deadline as SetSendDeadline",
		func(t *testing.T) {
			var (
				wantDeadline = time.Now()
				conn         = bmock.NewConn().RegisterSetWriteDeadline(
					func(deadline time.Time) (err error) {
						if deadline != wantDeadline {
							err = fmt.Errorf("unexpected deadline, want '%v' actual '%v'",
								wantDeadline,
								deadline)
						}
						return
					},
				)
				mocks     = []*mok.Mock{conn.Mock}
				transport = New[base.Cmd[any], base.Result](conn, nil, nil, nil)
			)
			transport.SetSendDeadline(wantDeadline)
			if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
				t.Error(infomap)
			}
		})

	t.Run("If Conn.SetWriteDeadline fails with an error, SetSendDeadline should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Conn.SetWriteDeadline error")
				conn    = bmock.NewConn().RegisterSetWriteDeadline(
					func(deadline time.Time) (err error) { return wantErr },
				)
				mocks     = []*mok.Mock{conn.Mock}
				transport = New[base.Cmd[any], base.Result](conn, nil, nil, nil)
				err       = transport.SetSendDeadline(time.Time{})
			)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
			if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
				t.Error(infomap)
			}
		})

	t.Run("Send should encode data with help of the Codec", func(t *testing.T) {
		var (
			wantSeq base.Seq = 1
			wantCmd          = bmock.NewCmd()
			writer           = tmock.NewWriter()
			codec            = tmock.NewClientCodec().RegisterEncode(
				func(seq base.Seq, cmd base.Cmd[any], w transport.Writer) (err error) {
					if seq != wantSeq {
						t.Errorf("unexpected seq, want '%v' actual '%v'", wantSeq, seq)
					}
					if cmd != wantCmd {
						t.Errorf("unexpected cmd, want '%v' actual '%v'", wantCmd, cmd)
					}
					return nil
				},
			)
			mocks     = []*mok.Mock{writer.Mock, codec.Mock}
			transport = New[base.Cmd[any], base.Result](nil, writer, nil,
				codec)
			err = transport.Send(wantSeq, wantCmd)
		)
		if err != nil {
			t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
		}
		if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
			t.Error(infomap)
		}
	})

	t.Run("If Codec.Encode fails with an error, Send should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Codec.Encode error")
				codec   = tmock.NewClientCodec().RegisterEncode(
					func(seq base.Seq, cmd base.Cmd[any], w transport.Writer) (err error) {
						return wantErr
					},
				)
				mocks     = []*mok.Mock{codec.Mock}
				transport = New[base.Cmd[any], base.Result](nil, nil, nil, codec)
				err       = transport.Send(1, nil)
			)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
			if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
				t.Error(infomap)
			}
		})

	t.Run("Conn.SetReadDeadline should receive same deadline as SetReceiveDeadline",
		func(t *testing.T) {
			var (
				wantDeadline = time.Now()
				conn         = bmock.NewConn().RegisterSetReadDeadline(
					func(deadline time.Time) (err error) {
						if deadline != wantDeadline {
							err = fmt.Errorf("unexpected deadline, want '%v' actual '%v'",
								wantDeadline,
								deadline)
						}
						return
					},
				)
				mocks     = []*mok.Mock{conn.Mock}
				transport = New[base.Cmd[any], base.Result](conn, nil, nil, nil)
			)
			transport.SetReceiveDeadline(wantDeadline)
			if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
				t.Error(infomap)
			}
		})

	t.Run("If Conn.SetReadDeadline fails with an error, SetReceiveDeadline should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Conn.SetReadDeadline error")
				conn    = bmock.NewConn().RegisterSetReadDeadline(
					func(deadline time.Time) (err error) { return wantErr },
				)
				mocks     = []*mok.Mock{conn.Mock}
				transport = New[base.Cmd[any], base.Result](conn, nil, nil, nil)
				err       = transport.SetReceiveDeadline(time.Time{})
			)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
			if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
				t.Error(infomap)
			}
		})

	t.Run("Receive should decode data with help of the Codec", func(t *testing.T) {
		var (
			wantSeq    base.Seq = 1
			wantResult          = bmock.NewResult()
			codec               = tmock.NewClientCodec().RegisterDecode(
				func(r transport.Reader) (seq base.Seq, result base.Result, err error) {
					return wantSeq, wantResult, nil
				},
			)
			mocks     = []*mok.Mock{codec.Mock}
			transport = New[base.Cmd[any], base.Result](nil, nil, nil,
				codec)
			seq, result, err = transport.Receive()
		)
		if err != nil {
			t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
		}
		if seq != wantSeq {
			t.Errorf("unexpected seq, want '%v' actual '%v'", wantSeq, seq)
		}
		if result != wantResult {
			t.Errorf("unexpected result, want '%v' actual '%v'", wantResult, result)
		}
		if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
			t.Error(infomap)
		}
	})

	t.Run("If Codec.Decode fails with an error, Receive should return it",
		func(t *testing.T) {
			var (
				wantSeq    base.Seq    = 0
				wantResult base.Result = nil
				wantErr                = errors.New("Codec.Decode error")
				codec                  = tmock.NewClientCodec().RegisterDecode(
					func(r transport.Reader) (seq base.Seq, result base.Result, err error) {
						err = wantErr
						return
					},
				)
				mocks     = []*mok.Mock{codec.Mock}
				transport = New[base.Cmd[any], base.Result](nil, nil, nil,
					codec)
				seq, result, err = transport.Receive()
			)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
			if seq != wantSeq {
				t.Errorf("unexpected seq, want '%v' actual '%v'", wantSeq, seq)
			}
			if result != wantResult {
				t.Errorf("unexpected result, want '%v' actual '%v'", wantResult, result)
			}
			if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
				t.Error(infomap)
			}
		})

	t.Run("Close should close the conn", func(t *testing.T) {
		var (
			conn = bmock.NewConn().RegisterClose(
				func() (err error) { return nil },
			)
			mocks     = []*mok.Mock{conn.Mock}
			transport = New[base.Cmd[any], base.Result](conn, nil, nil, nil)
			err       = transport.Close()
		)
		if err != nil {
			t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
		}
		if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
			t.Error(infomap)
		}
	})

	t.Run("If Conn.Close fails with an error, Close should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Conn.Close error")
				conn    = bmock.NewConn().RegisterClose(
					func() (err error) { return wantErr },
				)
				mocks     = []*mok.Mock{conn.Mock}
				transport = New[base.Cmd[any], base.Result](conn, nil, nil, nil)
				err       = transport.Close()
			)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
			if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
				t.Error(infomap)
			}
		})

	t.Run("If Writer.Flus fails with an error, Flush should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Writer.Flush error")
				writer  = tmock.NewWriter().RegisterFlush(
					func() error { return wantErr },
				)
				mocks     = []*mok.Mock{writer.Mock}
				transport = Transport[any, any]{w: writer}
				err       = transport.Flush()
			)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
			if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
				t.Error(infomap)
			}
		})

}

func SameTime(t1, t2 time.Time) bool {
	return !(t1.Before(t2.Truncate(Delta)) || t1.After(t2.Add(Delta)))
}
