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

func (mock ClientCodec) RegisterDecode(
	fn func(r transport.Reader) (seq base.Seq, result base.Result, err error),
) ClientCodec {
	mock.Register("Decode", fn)
	return mock
}

func (mock ClientCodec) RegisterEncode(
	fn func(seq base.Seq, cmd base.Cmd[any], w transport.Writer) (err error),
) ClientCodec {
	mock.Register("Encode", fn)
	return mock
}

func (mock ClientCodec) RegisterSize(
	fn func(cmd base.Cmd[any]) (size int),
) ClientCodec {
	mock.Register("Size", fn)
	return mock
}

func (mock ClientCodec) Decode(r transport.Reader) (seq base.Seq, result base.Result, err error) {
	vals, err := mock.Call("Decode", mok.SafeVal[transport.Reader](r))
	if err != nil {
		panic(err)
	}
	seq = vals[0].(base.Seq)
	result, _ = vals[1].(base.Result)
	err, _ = vals[2].(error)
	return
}

func (mock ClientCodec) Encode(seq base.Seq, cmd base.Cmd[any], w transport.Writer) (
	err error) {
	vals, err := mock.Call("Encode", seq, mok.SafeVal[base.Cmd[any]](cmd),
		mok.SafeVal[transport.Writer](w))
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (mock ClientCodec) Size(cmd base.Cmd[any]) (size int) {
	vals, err := mock.Call("Size", mok.SafeVal[base.Cmd[any]](cmd))
	if err != nil {
		panic(err)
	}
	size = vals[0].(int)
	return
}
