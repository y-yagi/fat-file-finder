package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/bytefmt"
)

func search(thresholdSize int64, w io.Writer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Size() > thresholdSize {
			fmt.Fprintf(w, "%s (%s)\n", path, bytefmt.ByteSize(uint64(info.Size())))
		}
		return nil
	}
}

func main() {
	const version = "0.1.0"

	var fileSizeStr string
	var path string
	var showVersion bool

	flag.StringVar(&fileSizeStr, "s", "100M", "Threshold size to display.")
	flag.StringVar(&path, "p", ".", "Search path.")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Println("version:", version)
		os.Exit(0)
		return
	}

	thresholdSize, err := bytefmt.ToBytes(fileSizeStr)
	if err != nil {
		fmt.Printf("Threshold size is invalid value. %v\n", err)
		os.Exit(1)
		return
	}

	err = filepath.Walk(path, search(int64(thresholdSize), os.Stdout))
	if err != nil {
		fmt.Printf("File read error. %v\n", err)
		os.Exit(1)
	}
}
