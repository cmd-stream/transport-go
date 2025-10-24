package mock

import (
	"github.com/ymz-ncnk/mok"
)

func NewWriter() Writer {
	return Writer{
		Mock: mok.New("Writer"),
	}
}

type Writer struct {
	*mok.Mock
}

func (w Writer) RegisterFlush(fn func() error) Writer {
	w.Register("Flush", fn)
	return w
}

func (w Writer) RegisterWriteByte(fn func(b byte) error) Writer {
	w.Register("WriteByte", fn)
	return w
}

func (w Writer) RegisterWrite(fn func(p []byte) (int, error)) Writer {
	w.Register("Write", fn)
	return w
}

func (w Writer) WriteByte(b byte) (err error) {
	vals, err := w.Call("WriteByte", b)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (w Writer) Write(p []byte) (n int, err error) {
	vals, err := w.Call("Write", p)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (w Writer) WriteString(s string) (n int, err error) {
	panic("not implemented")
}

func (w Writer) Flush() (err error) {
	vals, err := w.Call("Flush")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}
