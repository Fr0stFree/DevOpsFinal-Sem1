package utils

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"
	"path/filepath"
)

func UntarFile(r io.Reader) (io.ReadCloser, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	tarReader := tar.NewReader(bytes.NewReader(data))
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if filepath.Ext(header.Name) != ".csv" || filepath.Base(header.Name)[0] == '.' {
			continue
		}
		return io.NopCloser(tarReader), nil
	}
	return nil, errors.New("file not found in the tar archive")
}