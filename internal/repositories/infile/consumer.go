package infile

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewConsumer(filename string) (*consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *consumer) initializeStorage() map[string]string {
	defer c.file.Close()
	initializedStorage := make(map[string]string)
	for c.scanner.Scan() {
		reduceURL, readErr := readURL(c)
		if readErr != nil {
			log.Fatal(readErr)
		}
		initializedStorage[reduceURL.URLID] = reduceURL.URL
	}
	if err := c.scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return initializedStorage
}

func readURL(c *consumer) (*savingURL, error) {
	data := c.scanner.Bytes()

	su := savingURL{}
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
