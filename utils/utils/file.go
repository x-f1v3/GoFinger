package utils

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func IsFile(path string) bool {
	_, err := os.Stat(path)

	return err == nil

}

func ReadFile(path string) ([]string, error) {
	list := make([]string, 0)

	if !IsFile(path) {
		return nil, os.ErrNotExist
	}
	file, _ := os.Open(path)
	reader := bufio.NewReader(file)
	for {

		buf, _, err := reader.ReadLine()
		line := string(buf)
		line = strings.TrimSpace(line)

		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}

		}
		list = append(list, line)
	}
	return list, nil

}
