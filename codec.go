package transport

import (
	"io"

	"github.com/cmd-stream/base-go"
)

// Writer is an interface that extends io.ByteWriter, io.Writer, and
// io.StringWriter. It also includes a Flush method to ensure buffered data is
// written out.
type Writer interface {
	io.ByteWriter
	io.Writer
	io.StringWriter
	Flush() error
}

// Reader is an interface that extends io.Reader and io.ByteReader.
type Reader interface {
	io.Reader
	io.ByteReader
}

// Codec is responsible for encoding and decoding data transmitted over a
// connection:
//   - On the client side: encodes Commands to send and decodes Results received
//     from the server.
//   - On the server side: decodes Commands from clients, validates them, and
//     encodes Results to send back.
type Codec[T, V any] interface {
	Encode(seq base.Seq, t T, w Writer) (n int, err error)
	Decode(r Reader) (seq base.Seq, v V, n int, err error)
}
