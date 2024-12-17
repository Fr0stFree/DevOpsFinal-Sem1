package archivers

import (
	"io"
)

type Format string

const (
	ZipFmt Format = "zip"
	TarFmt Format = "tar"
)

type Archiver interface {
	Archive(w io.Writer, fileName string) (io.WriteCloser, error)
	Extract(r io.Reader) (io.ReadCloser, error)
}

func New(fmt Format) Archiver {
	switch fmt {
	case ZipFmt:
		return &zipArchiver{}
	case TarFmt:
		return &tarArchiver{}
	default:
		return &zipArchiver{}
	}
}
