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

// Codec is responsible for encoding and decoding data transmitted over the
// connection:
//
//   - On the client side: Encodes Commands and decodes Results received from the
//     server. If the server imposes a Command size limit, the client can use the
//     Codec.Size() method to ensure that Commands being sent comply with the size
//     restriction.
//
//   - On the server side: Decodes Commands received from the client and encodes
//     Results to be sent back. During the Decode process, the server can validate
//     the Command length to enforce size constraints.
type Codec[T, V any] interface {
	Encode(seq base.Seq, t T, w Writer) (err error)
	Decode(r Reader) (seq base.Seq, v V, err error)
	Size(t T) int
}
