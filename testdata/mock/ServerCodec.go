package mock

import (
	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/transport-go"
	"github.com/ymz-ncnk/mok"
)

type (
	DecodeServerFn func(r transport.Reader) (seq core.Seq, cmd core.Cmd[any], n int, err error)
	EncodeServerFn func(seq core.Seq, result core.Result, w transport.Writer) (n int, err error)
)

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
	fn func(result core.Result) (size int),
) ServerCodec {
	c.Register("Size", fn)
	return c
}

func (c ServerCodec) Decode(r transport.Reader) (seq core.Seq, cmd core.Cmd[any],
	n int, err error,
) {
	vals, err := c.Call("Decode", mok.SafeVal[transport.Reader](r))
	if err != nil {
		panic(err)
	}
	seq = vals[0].(core.Seq)
	cmd, _ = vals[1].(core.Cmd[any])
	n = vals[2].(int)
	err, _ = vals[3].(error)
	return
}

func (c ServerCodec) Encode(seq core.Seq, result core.Result, w transport.Writer) (
	n int, err error,
) {
	vals, err := c.Call("Encode", seq, mok.SafeVal[core.Result](result),
		mok.SafeVal[transport.Writer](w))
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (c ServerCodec) Size(result core.Result) (size int) {
	vals, err := c.Call("Size", mok.SafeVal[core.Result](result))
	if err != nil {
		panic(err)
	}
	size = vals[0].(int)
	return
}
