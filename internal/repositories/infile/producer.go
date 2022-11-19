package infile

import (
	"encoding/json"
	"io"
	"os"
)

type producer struct {
	file    io.ReadWriteCloser
	encoder *json.Encoder
}

func NewProducer(filePath string) (*producer, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) WriteURL(su *savingURL) error {
	return p.encoder.Encode(su)
}

func (p *producer) Close() error {
	return p.file.Close()
}
