package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"code.cloudfoundry.org/bytefmt"
)

const (
	typeFile = "f"
	typeDir  = "d"
)

func search(location string, thresholdSize int64, outStream io.Writer) {
	entries, err := os.ReadDir(location)
	if err != nil {
		fmt.Fprintf(outStream, "%v\n", err)
		return
	}

	var dirSize int64
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Fprintf(outStream, "%v\n", err)
		}

		fullPath := filepath.Join(location, info.Name())
		if info.IsDir() {
			search(fullPath, int64(thresholdSize), outStream)
		} else if info.Size() > int64(thresholdSize) {
			fmt.Fprintf(outStream, "%s %s (%s)\n", typeFile, fullPath, bytefmt.ByteSize(uint64(info.Size())))
			dirSize += info.Size()
		}
	}
	if dirSize > thresholdSize {
		fmt.Fprintf(outStream, "%s %s/ (%s)\n", typeDir, location, bytefmt.ByteSize(uint64(dirSize)))
	}
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

	entries, err := os.ReadDir(location)
	if err != nil {
		fmt.Fprintf(outStream, "Location is invalid value. %v\n", err)
		exitCode = 1
		return
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Fprintf(outStream, "%v\n", err)
		}

		fullPath := filepath.Join(location, info.Name())
		if info.IsDir() {
			wg.Add(1)
			go func() {
				search(fullPath, int64(thresholdSize), outStream)
				wg.Done()
			}()
		} else if info.Size() > int64(thresholdSize) {
			fmt.Fprintf(outStream, "%s %s (%s)\n", typeFile, fullPath, bytefmt.ByteSize(uint64(info.Size())))
		}
	}
	wg.Wait()

	return
}

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}
