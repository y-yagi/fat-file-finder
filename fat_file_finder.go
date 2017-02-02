package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/bytefmt"
)

func search(location string, thresholdSize int64) chan string {
	chann := make(chan string)

	go func() {
		_ = filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				chann <- fmt.Sprintf("%v", err)
				return nil
			}

			if info.IsDir() {
				return nil
			}

			if info.Size() > thresholdSize {
				chann <- fmt.Sprintf("%s (%s)", path, bytefmt.ByteSize(uint64(info.Size())))
			}
			return nil
		})
		defer close(chann)
	}()
	return chann
}

func main() {
	const version = "0.1.0"

	var fileSizeStr string
	var location string
	var showVersion bool

	flag.StringVar(&fileSizeStr, "s", "100M", "Threshold size to display.")
	flag.StringVar(&location, "l", ".", "Search location.")
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

	chann := search(location, int64(thresholdSize))
	for msg := range chann {
		fmt.Println(msg)
	}
}
