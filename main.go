package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	projectDir := "." // adjust if needed
	outputFile := "all_code.txt"

	out, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	err = filepath.WalkDir(projectDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		ext := filepath.Ext(path)

		// Skip unwanted files
		if base == "config.go" || base == "connect.go" || base == "Dockerfile" {
			return nil
		}
		if ext == ".pb.go" {
			return nil
		}

		// Only read .go and .proto files
		if ext != ".go" && ext != ".proto" {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		_, err = out.WriteString(fmt.Sprintf("==== FILE: %s ====\n%s\n\n", path, string(data)))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("All relevant code has been written to", outputFile)
}