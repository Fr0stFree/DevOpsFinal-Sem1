package utils

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"path/filepath"
)

func UnzipFile(r io.Reader) (io.ReadCloser, error) {
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
	return nil, errors.New("file not found in the zi[] archive")
}

func ZipFile(buffer *bytes.Buffer, w io.Writer, fileName string) error {
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	file, err := zipWriter.Create(fileName)
	if err != nil {
		return err
	}
	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}
