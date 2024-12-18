package archivers

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"path/filepath"
)

type zipArchiver struct{}

func (z *zipArchiver) Extract(r io.Reader) (io.ReadCloser, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	for _, file := range zipReader.File {
		if filepath.Ext(file.Name) != ".csv" {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		return rc, nil
	}
	return nil, errors.New("file not found in the zip archive")
}

func (z *zipArchiver) Archive(w io.Writer, fileName string) (io.WriteCloser, error) {
	zipWriter := zip.NewWriter(w)
	file, err := zipWriter.Create(fileName)
	if err != nil {
		return nil, err
	}
	return &zipWriteCloser{zipWriter, file}, nil
}

type zipWriteCloser struct {
	zipWriter *zip.Writer
	file      io.Writer
}

func (z *zipWriteCloser) Write(p []byte) (n int, err error) {
	return z.file.Write(p)
}

func (z *zipWriteCloser) Close() error {
	return z.zipWriter.Close()
}
