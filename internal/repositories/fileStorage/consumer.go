package fileStorage

import (
	"bufio"
	"encoding/json"
	"os"
)

type consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewConsumer(filename string) (*consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *consumer) ReadURL() (*SavingURL, error) {
	data := c.scanner.Bytes()

	su := SavingURL{}
	if len(data) > 0 {
		err := json.Unmarshal(data, &su)
		if err != nil {
			return nil, err
		}
	}

	return &su, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}
