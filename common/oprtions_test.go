package tcom

import (
	"testing"
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

	if o.WriterBufSize != wantWriterBufSize {
		t.Errorf("unexpected WriterBufSize, want %v actual %v", wantWriterBufSize,
			o.WriterBufSize)
	}

	if o.ReaderBufSize != wantReaderBufSize {
		t.Errorf("unexpected ReaderBufSize, want %v actual %v", wantReaderBufSize,
			o.ReaderBufSize)
	}

}
