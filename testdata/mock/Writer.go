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

func (mock Writer) RegisterFlush(fn func() error) Writer {
	mock.Register("Flush", fn)
	return mock
}

func (mock Writer) RegisterWriteByte(fn func(b byte) error) Writer {
	mock.Register("WriteByte", fn)
	return mock
}

func (mock Writer) WriteByte(b byte) (err error) {
	vals, err := mock.Call("WriteByte", b)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (mock Writer) Write(p []byte) (n int, err error) {
	panic("not implemented")
}

func (mock Writer) WriteString(s string) (n int, err error) {
	panic("not implemented")
}

func (mock Writer) Flush() (err error) {
	vals, err := mock.Call("Flush")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}
