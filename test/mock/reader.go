package mock

import (
	"github.com/ymz-ncnk/mok"
)

func NewReader() Reader {
	return Reader{
		Mock: mok.New("Reader"),
	}
}

type Reader struct {
	*mok.Mock
}

func (r Reader) RegisterRead(fn func(p []byte) (n int, err error)) Reader {
	r.Register("Read", fn)
	return r
}

func (r Reader) RegisterReadByte(fn func() (b byte, err error)) Reader {
	r.Register("ReadByte", fn)
	return r
}

func (r Reader) Read(p []byte) (n int, err error) {
	vals, err := r.Call("Read", p)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (r Reader) ReadByte() (b byte, err error) {
	vals, err := r.Call("ReadByte")
	if err != nil {
		panic(err)
	}
	b = vals[0].(byte)
	err, _ = vals[1].(error)
	return
}
