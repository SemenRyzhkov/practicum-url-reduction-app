package fileRepository

import (
	"encoding/json"
	"os"
)

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
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
