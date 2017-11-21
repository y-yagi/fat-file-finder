package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"code.cloudfoundry.org/bytefmt"
)

func search(location string, thresholdSize int64, wg *sync.WaitGroup) {
	err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Size() > thresholdSize {
			fmt.Printf("%s (%s)\n", path, bytefmt.ByteSize(uint64(info.Size())))
		}
		return nil
	})

	if err != nil {
		fmt.Printf("%v\n", err)
	}
	wg.Done()
}

func main() {
	const version = "0.1.0"

	var fileSizeStr string
	var location string
	var showVersion bool
	var wg sync.WaitGroup

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

	fileInfos, err := ioutil.ReadDir(location)
	if err != nil {
		fmt.Printf("Location is invalid value. %v\n", err)
		os.Exit(1)
		return
	}

	for _, fileInfo := range fileInfos {
		fullPath := filepath.Join(location, fileInfo.Name())
		if fileInfo.IsDir() {
			wg.Add(1)
			go search(fullPath, int64(thresholdSize), &wg)
		} else if fileInfo.Size() > int64(thresholdSize) {
			fmt.Printf("%s (%s)\n", fullPath, bytefmt.ByteSize(uint64(fileInfo.Size())))
		}
	}
	wg.Wait()
}
