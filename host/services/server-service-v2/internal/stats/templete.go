package stats

import (
	"io"
	"os"
)

func CreateTemplete() ([]byte, error) {
	file, err := os.Open("templete.json")
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	templete, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}

	return templete, nil
}