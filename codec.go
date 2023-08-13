package transport

import (
	"io"

	"github.com/cmd-stream/base-go"
)

// Writer is the interface that groups the WriteByte, Write, WriteString and
// Flush methods.
type Writer interface {
	io.ByteWriter
	io.Writer
	io.StringWriter
	Flush() error
}

// Reader is the interface that groups the basic ReadByte and Read methods.
type Reader interface {
	io.Reader
	io.ByteReader
}

// Codec encodes/decodes data to/from the connection.
//
// Client codec encodes commands and decodes results. If the server imposes a
// command size limit, the client delegate, using the Size() method, can
// determine whether the size of the command being sent is small enough.
//
// Server codec decodes commands and encodes results. In the Decode method it
// can check the length of the command.
type Codec[T, V any] interface {
	Encode(seq base.Seq, t T, w Writer) (err error)
	Decode(r Reader) (seq base.Seq, v V, err error)
	Size(t T) int
}
