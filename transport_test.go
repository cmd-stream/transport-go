package transport_test

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/cmd-stream/core-go"
	cmock "github.com/cmd-stream/core-go/testdata/mock"
	"github.com/cmd-stream/transport-go"
	tmock "github.com/cmd-stream/transport-go/testdata/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

func TestTransport(t *testing.T) {
	t.Run("LocalAddr should return local address of the conn",
		func(t *testing.T) {
			var (
				wantAddr = &net.TCPAddr{IP: net.ParseIP("127.0.0.1:9000")}
				conn     = cmock.NewConn().RegisterLocalAddr(
					func() (addr net.Addr) {
						return wantAddr
					},
				)
				mocks     = []*mok.Mock{conn.Mock}
				transport = transport.New[core.Cmd[any], core.Result](conn, nil, nil, nil)
			)
			addr := transport.LocalAddr()
			if addr != wantAddr {
				t.Errorf("unexpected addr, want '%v' actual '%v'", wantAddr, addr)
			}
			asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
		})

	t.Run("RemoteAddr should return remote address of the conn",
		func(t *testing.T) {
			var (
				wantAddr = &net.TCPAddr{IP: net.ParseIP("127.0.0.1:9000")}
				conn     = cmock.NewConn().RegisterRemoteAddr(
					func() (addr net.Addr) {
						return wantAddr
					},
				)
				mocks     = []*mok.Mock{conn.Mock}
				transport = transport.New[core.Cmd[any], core.Result](conn, nil, nil, nil)
			)
			addr := transport.RemoteAddr()
			if addr != wantAddr {
				t.Errorf("unexpected addr, want '%v' actual '%v'", wantAddr, addr)
			}
			asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
		})

	t.Run("Conn.SetWriteDeadline should receive same deadline as SetSendDeadline",
		func(t *testing.T) {
			var (
				wantDeadline = time.Now()
				conn         = cmock.NewConn().RegisterSetWriteDeadline(
					func(deadline time.Time) (err error) {
						asserterror.Equal(deadline, wantDeadline, t)
						return
					},
				)
				mocks     = []*mok.Mock{conn.Mock}
				transport = transport.New[core.Cmd[any], core.Result](conn, nil, nil, nil)
			)
			transport.SetSendDeadline(wantDeadline)
			asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
		})

	t.Run("If Conn.SetWriteDeadline fails with an error, SetSendDeadline should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Conn.SetWriteDeadline error")
				conn    = cmock.NewConn().RegisterSetWriteDeadline(
					func(deadline time.Time) (err error) { return wantErr },
				)
				mocks     = []*mok.Mock{conn.Mock}
				transport = transport.New[core.Cmd[any], core.Result](conn, nil, nil, nil)
				err       = transport.SetSendDeadline(time.Time{})
			)
			asserterror.Equal(err, wantErr, t)
			asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
		})

	t.Run("Send should encode data with help of the Codec", func(t *testing.T) {
		var (
			wantSeq core.Seq = 1
			wantCmd          = cmock.NewCmd()
			wantN   int      = 3
			wantErr error    = nil
			writer           = tmock.NewWriter()
			codec            = tmock.NewClientCodec().RegisterEncode(
				func(seq core.Seq, cmd core.Cmd[any], w transport.Writer) (n int, err error) {
					asserterror.Equal(seq, wantSeq, t)
					asserterror.Equal[any](cmd, wantCmd, t)
					return wantN, wantErr
				},
			)
			mocks     = []*mok.Mock{writer.Mock, codec.Mock}
			transport = transport.New(nil, writer, nil, codec)
			n, err    = transport.Send(wantSeq, wantCmd)
		)
		asserterror.EqualError(err, wantErr, t)
		asserterror.Equal(n, wantN, t)
		asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
	})

	t.Run("If Codec.Encode fails with an error, Send should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Codec.Encode error")
				codec   = tmock.NewClientCodec().RegisterEncode(
					func(seq core.Seq, cmd core.Cmd[any], w transport.Writer) (n int, err error) {
						return 0, wantErr
					},
				)
				mocks     = []*mok.Mock{codec.Mock}
				transport = transport.New(nil, nil, nil, codec)
				_, err    = transport.Send(1, nil)
			)
			asserterror.EqualError(err, wantErr, t)
			asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
		})

	t.Run("Conn.SetReadDeadline should receive same deadline as SetReceiveDeadline",
		func(t *testing.T) {
			var (
				wantDeadline = time.Now()
				conn         = cmock.NewConn().RegisterSetReadDeadline(
					func(deadline time.Time) (err error) {
						asserterror.Equal(deadline, wantDeadline, t)
						return
					},
				)
				mocks     = []*mok.Mock{conn.Mock}
				transport = transport.New[core.Cmd[any], core.Result](conn, nil, nil, nil)
			)
			transport.SetReceiveDeadline(wantDeadline)
			asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
		})

	t.Run("If Conn.SetReadDeadline fails with an error, SetReceiveDeadline should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Conn.SetReadDeadline error")
				conn    = cmock.NewConn().RegisterSetReadDeadline(
					func(deadline time.Time) (err error) { return wantErr },
				)
				mocks     = []*mok.Mock{conn.Mock}
				transport = transport.New[core.Cmd[any], core.Result](conn, nil, nil, nil)
				err       = transport.SetReceiveDeadline(time.Time{})
			)
			asserterror.EqualError(err, wantErr, t)
			asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
		})

	t.Run("Receive should decode data with help of the Codec", func(t *testing.T) {
		var (
			wantSeq    core.Seq = 1
			wantResult          = cmock.NewResult()
			wantN      int      = 3
			wantErr    error    = nil
			codec               = tmock.NewClientCodec().RegisterDecode(
				func(r transport.Reader) (seq core.Seq, result core.Result, n int, err error) {
					return wantSeq, wantResult, wantN, wantErr
				},
			)
			mocks               = []*mok.Mock{codec.Mock}
			transport           = transport.New(nil, nil, nil, codec)
			seq, result, n, err = transport.Receive()
		)
		asserterror.EqualError(err, wantErr, t)
		asserterror.Equal(seq, wantSeq, t)
		asserterror.EqualDeep(result, wantResult, t)
		asserterror.Equal(n, wantN, t)
		asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
	})

	t.Run("If Codec.Decode fails with an error, Receive should return it",
		func(t *testing.T) {
			var (
				wantSeq    core.Seq    = 0
				wantResult core.Result = nil
				wantErr                = errors.New("Codec.Decode error")
				codec                  = tmock.NewClientCodec().RegisterDecode(
					func(r transport.Reader) (seq core.Seq, result core.Result, n int, err error) {
						err = wantErr
						return
					},
				)
				mocks               = []*mok.Mock{codec.Mock}
				transport           = transport.New(nil, nil, nil, codec)
				seq, result, _, err = transport.Receive()
			)
			asserterror.EqualError(err, wantErr, t)
			asserterror.Equal(seq, wantSeq, t)
			asserterror.EqualDeep(result, wantResult, t)
			asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
		})

	t.Run("Close should close the conn", func(t *testing.T) {
		var (
			wantErr error = nil
			conn          = cmock.NewConn().RegisterClose(
				func() (err error) { return nil },
			)
			mocks     = []*mok.Mock{conn.Mock}
			transport = transport.New[core.Cmd[any], core.Result](conn, nil, nil, nil)
			err       = transport.Close()
		)
		asserterror.EqualError(err, wantErr, t)
		asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
	})

	t.Run("If Conn.Close fails with an error, Close should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Conn.Close error")
				conn    = cmock.NewConn().RegisterClose(
					func() (err error) { return wantErr },
				)
				mocks     = []*mok.Mock{conn.Mock}
				transport = transport.New[core.Cmd[any], core.Result](conn, nil, nil, nil)
				err       = transport.Close()
			)
			asserterror.EqualError(err, wantErr, t)
			asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
		})

	t.Run("If Writer.Flus fails with an error, Flush should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Writer.Flush error")
				writer  = tmock.NewWriter().RegisterFlush(
					func() error { return wantErr },
				)
				mocks     = []*mok.Mock{writer.Mock}
				transport = transport.Transport[any, any]{W: writer}
				err       = transport.Flush()
			)
			asserterror.EqualError(err, wantErr, t)
			asserterror.EqualDeep(mok.CheckCalls(mocks), mok.EmptyInfomap, t)
		})
}
