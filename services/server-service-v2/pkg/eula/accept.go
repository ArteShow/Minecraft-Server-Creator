package eula

import (
	"os"
	"path/filepath"
)

func Accept(dir string) error {
	path := filepath.Join(dir, "eula.txt")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("eula=true")
	if err != nil {
		return err
	}

	return nil
}
