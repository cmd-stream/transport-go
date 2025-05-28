package mock

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	"github.com/ymz-ncnk/mok"
)

type DecodeFn func(r transport.Reader) (seq base.Seq, result base.Result, n int, err error)
type EncodeFn func(seq base.Seq, cmd base.Cmd[any], w transport.Writer) (n int, err error)

func NewClientCodec() ClientCodec {
	return ClientCodec{
		Mock: mok.New("ClientCodec"),
	}
}

type ClientCodec struct {
	*mok.Mock
}

func (c ClientCodec) RegisterDecode(fn DecodeFn) ClientCodec {
	c.Register("Decode", fn)
	return c
}

func (c ClientCodec) RegisterEncode(fn EncodeFn) ClientCodec {
	c.Register("Encode", fn)
	return c
}

func (c ClientCodec) RegisterSize(
	fn func(cmd base.Cmd[any]) (size int),
) ClientCodec {
	c.Register("Size", fn)
	return c
}

func (c ClientCodec) Decode(r transport.Reader) (seq base.Seq, result base.Result, n int, err error) {
	vals, err := c.Call("Decode", mok.SafeVal[transport.Reader](r))
	if err != nil {
		panic(err)
	}
	seq = vals[0].(base.Seq)
	result, _ = vals[1].(base.Result)
	n = vals[2].(int)
	err, _ = vals[3].(error)
	return
}

func (c ClientCodec) Encode(seq base.Seq, cmd base.Cmd[any], w transport.Writer) (
	n int, err error) {
	vals, err := c.Call("Encode", seq, mok.SafeVal[base.Cmd[any]](cmd),
		mok.SafeVal[transport.Writer](w))
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
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
