package transport

import (
	"testing"

	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestOptions(t *testing.T) {
	var (
		o                 = Options{}
		wantWriterBufSize = 1
		wantReaderBufSize = 1
	)
	Apply([]SetOption{
		WithWriterBufSize(wantWriterBufSize),
		WithReaderBufSize(wantReaderBufSize),
	}, &o)

	asserterror.Equal(o.WriterBufSize, wantWriterBufSize, t)
	asserterror.Equal(o.ReaderBufSize, wantReaderBufSize, t)
}
