package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"code.cloudfoundry.org/bytefmt"
)

func search(location string, thresholdSize int64, outStream io.Writer, wg *sync.WaitGroup) {
	err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Size() > thresholdSize {
			fmt.Fprintf(outStream, "%s (%s)\n", path, bytefmt.ByteSize(uint64(info.Size())))
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(outStream, "%v\n", err)
	}
	wg.Done()
}

func run(args []string, outStream, errStream io.Writer) (exitCode int) {
	const version = "0.1.0"

	var fileSizeStr string
	var location string
	var showVersion bool
	var wg sync.WaitGroup

	flags := flag.NewFlagSet("fat-file-finder", flag.ExitOnError)
	flags.SetOutput(errStream)
	flags.StringVar(&fileSizeStr, "s", "100M", "Threshold size to display.")
	flags.StringVar(&location, "l", ".", "Search location.")
	flags.BoolVar(&showVersion, "v", false, "show version")
	flags.Parse(args[1:])

	exitCode = 0

	if showVersion {
		fmt.Fprintln(outStream, "version:", version)
		return
	}

	thresholdSize, err := bytefmt.ToBytes(fileSizeStr)
	if err != nil {
		fmt.Fprintf(outStream, "Threshold size is invalid value. %v\n", err)
		exitCode = 1
		return
	}

	fileInfos, err := ioutil.ReadDir(location)
	if err != nil {
		fmt.Fprintf(outStream, "Location is invalid value. %v\n", err)
		exitCode = 1
		return
	}

	for _, fileInfo := range fileInfos {
		fullPath := filepath.Join(location, fileInfo.Name())
		if fileInfo.IsDir() {
			wg.Add(1)
			go search(fullPath, int64(thresholdSize), outStream, &wg)
		} else if fileInfo.Size() > int64(thresholdSize) {
			fmt.Fprintf(outStream, "%s (%s)\n", fullPath, bytefmt.ByteSize(uint64(fileInfo.Size())))
		}
	}
	wg.Wait()

	return
}

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}
