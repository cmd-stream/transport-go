package mock

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	"github.com/ymz-ncnk/mok"
)

func NewClientCodec() ClientCodec {
	return ClientCodec{
		Mock: mok.New("ClientCodec"),
	}
}

type ClientCodec struct {
	*mok.Mock
}

func (c ClientCodec) RegisterDecode(
	fn func(r transport.Reader) (seq base.Seq, result base.Result, err error),
) ClientCodec {
	c.Register("Decode", fn)
	return c
}

func (c ClientCodec) RegisterEncode(
	fn func(seq base.Seq, cmd base.Cmd[any], w transport.Writer) (err error),
) ClientCodec {
	c.Register("Encode", fn)
	return c
}

func (c ClientCodec) RegisterSize(
	fn func(cmd base.Cmd[any]) (size int),
) ClientCodec {
	c.Register("Size", fn)
	return c
}

func (c ClientCodec) Decode(r transport.Reader) (seq base.Seq, result base.Result, err error) {
	vals, err := c.Call("Decode", mok.SafeVal[transport.Reader](r))
	if err != nil {
		panic(err)
	}
	seq = vals[0].(base.Seq)
	result, _ = vals[1].(base.Result)
	err, _ = vals[2].(error)
	return
}

func (c ClientCodec) Encode(seq base.Seq, cmd base.Cmd[any], w transport.Writer) (
	err error) {
	vals, err := c.Call("Encode", seq, mok.SafeVal[base.Cmd[any]](cmd),
		mok.SafeVal[transport.Writer](w))
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (c ClientCodec) Size(cmd base.Cmd[any]) (size int) {
	vals, err := c.Call("Size", mok.SafeVal[base.Cmd[any]](cmd))
	if err != nil {
		panic(err)
	}
	size = vals[0].(int)
	return
}
