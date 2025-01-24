package common

// Conf configures the Transport.
//
// WriterBufSize defines a buffer size for the Writer, if == 0, the default
// bufio.Writer size is used.
// ReaderBufSize defines a buffer size for the Reader, if == 0, the default
// bufio.Reader size is used.
type Conf struct {
	WriterBufSize int
	ReaderBufSize int
}
