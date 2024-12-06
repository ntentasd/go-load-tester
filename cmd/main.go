package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	models "github.com/ntentasd/go-load-tester/internal/models"
)

func main() {
	var file string
	flag.StringVar(&file, "file", "", "The yaml file to parse")
	flag.Parse()

	if file == "" {
		fmt.Println("Please select a file")
		os.Exit(1)
	}

	r := regexp.MustCompile(`(?i)\.ya?ml$`)
	if !r.MatchString(file) {
		fmt.Println("Invalid file type. Please provide a valid YAML file.")
		os.Exit(1)
	}

	fileHandle, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file. Does it exist?")
		os.Exit(1)
	}
	defer fileHandle.Close()

	f := io.Reader(fileHandle)

	yamlFile, err := models.UnmarshalYaml(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	yamlFile.Dbg()
}
