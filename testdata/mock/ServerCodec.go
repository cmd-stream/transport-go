package mock

import (
	"reflect"

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

func (mock ServerCodec) RegisterDecode(
	fn func(r transport.Reader) (seq base.Seq, cmd base.Cmd[any], err error),
) ServerCodec {
	mock.Register("Decode", fn)
	return mock
}

func (mock ServerCodec) RegisterEncode(
	fn func(seq base.Seq, result base.Result, w transport.Writer) (err error),
) ServerCodec {
	mock.Register("Encode", fn)
	return mock
}

func (mock ServerCodec) RegisterSize(
	fn func(result base.Result) (size int),
) ServerCodec {
	mock.Register("Size", fn)
	return mock
}

func (mock ServerCodec) Decode(r transport.Reader) (seq base.Seq, cmd base.Cmd[any],
	err error) {
	var rVal reflect.Value
	if r == nil {
		rVal = reflect.Zero(reflect.TypeOf((*transport.Reader)(nil)).Elem())
	} else {
		rVal = reflect.ValueOf(r)
	}
	vals, err := mock.Call("Decode", rVal)
	if err != nil {
		panic(err)
	}
	seq, _ = vals[0].(base.Seq)
	cmd, _ = vals[1].(base.Cmd[any])
	err, _ = vals[2].(error)
	return
}

func (mock ServerCodec) Encode(seq base.Seq, result base.Result, w transport.Writer) (
	err error) {
	var resultVal reflect.Value
	if result == nil {
		resultVal = reflect.Zero(reflect.TypeOf((*base.Result)(nil)).Elem())
	} else {
		resultVal = reflect.ValueOf(result)
	}
	var wVal reflect.Value
	if w == nil {
		wVal = reflect.Zero(reflect.TypeOf((*transport.Writer)(nil)).Elem())
	} else {
		wVal = reflect.ValueOf(w)
	}
	vals, err := mock.Call("Encode", seq, resultVal, wVal)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (mock ServerCodec) Size(result base.Result) (size int) {
	var resultVal reflect.Value
	if result == nil {
		resultVal = reflect.Zero(reflect.TypeOf((*base.Result)(nil)).Elem())
	} else {
		resultVal = reflect.ValueOf(result)
	}
	vals, err := mock.Call("Size", resultVal)
	if err != nil {
		panic(err)
	}
	size = vals[0].(int)
	return
}
