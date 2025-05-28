package mock

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	"github.com/ymz-ncnk/mok"
)

type DecodeServerFn func(r transport.Reader) (seq base.Seq, cmd base.Cmd[any], n int, err error)
type EncodeServerFn func(seq base.Seq, result base.Result, w transport.Writer) (n int, err error)

func NewServerCodec() ServerCodec {
	return ServerCodec{
		Mock: mok.New("ServerCodec"),
	}
}

type ServerCodec struct {
	*mok.Mock
}

func (c ServerCodec) RegisterDecode(fn DecodeServerFn) ServerCodec {
	c.Register("Decode", fn)
	return c
}

func (c ServerCodec) RegisterEncode(fn EncodeServerFn) ServerCodec {
	c.Register("Encode", fn)
	return c
}

func (c ServerCodec) RegisterSize(
	fn func(result base.Result) (size int),
) ServerCodec {
	c.Register("Size", fn)
	return c
}

func (c ServerCodec) Decode(r transport.Reader) (seq base.Seq, cmd base.Cmd[any],
	n int, err error) {
	vals, err := c.Call("Decode", mok.SafeVal[transport.Reader](r))
	if err != nil {
		panic(err)
	}
	seq = vals[0].(base.Seq)
	cmd, _ = vals[1].(base.Cmd[any])
	n = vals[2].(int)
	err, _ = vals[3].(error)
	return
}

func (c ServerCodec) Encode(seq base.Seq, result base.Result, w transport.Writer) (
	n int, err error) {
	vals, err := c.Call("Encode", seq, mok.SafeVal[base.Result](result),
		mok.SafeVal[transport.Writer](w))
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (c ServerCodec) Size(result base.Result) (size int) {
	vals, err := c.Call("Size", mok.SafeVal[base.Result](result))
	if err != nil {
		panic(err)
	}
	size = vals[0].(int)
	return
}
