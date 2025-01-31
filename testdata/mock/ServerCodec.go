package mock

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	"github.com/ymz-ncnk/mok"
)

func NewServerCodec() ServerCodec {
	return ServerCodec{
		Mock: mok.New("ServerCodec"),
	}
}

type ServerCodec struct {
	*mok.Mock
}

func (c ServerCodec) RegisterDecode(
	fn func(r transport.Reader) (seq base.Seq, cmd base.Cmd[any], err error),
) ServerCodec {
	c.Register("Decode", fn)
	return c
}

func (c ServerCodec) RegisterEncode(
	fn func(seq base.Seq, result base.Result, w transport.Writer) (err error),
) ServerCodec {
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
	err error) {
	vals, err := c.Call("Decode", mok.SafeVal[transport.Reader](r))
	if err != nil {
		panic(err)
	}
	seq, _ = vals[0].(base.Seq)
	cmd, _ = vals[1].(base.Cmd[any])
	err, _ = vals[2].(error)
	return
}

func (c ServerCodec) Encode(seq base.Seq, result base.Result, w transport.Writer) (
	err error) {
	vals, err := c.Call("Encode", seq, mok.SafeVal[base.Result](result),
		mok.SafeVal[transport.Writer](w))
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
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
