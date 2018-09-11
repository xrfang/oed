//+build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	bindata "github.com/tmthrgd/go-bindata"
)

func main() {
	genOpts := &bindata.GenerateOptions{
		Package:        "main",
		MemCopy:        true,
		Compress:       true,
		Metadata:       true,
		Restore:        true,
		AssetDir:       true,
		DecompressOnce: true,
	}
	findOpts := new(bindata.FindFilesOptions)

	flag.Usage = func() {
		fmt.Printf("Usage: %s options...\n\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	output := flag.String("o", "", "Name of the output file to be generated.")
	flag.StringVar(&findOpts.Prefix, "r", "", "Root directory of resources to embed")
	flag.BoolVar(&genOpts.Debug, "debug", genOpts.Debug, "Do not embed the assets, but provide the embedding API. Contents will still be loaded from disk.")
	flag.Parse()

	if *output == "" {
		fmt.Fprintln(os.Stderr, "Missing -o\n")
		flag.Usage()
		os.Exit(1)
	}

	if findOpts.Prefix == "" {
		fmt.Fprintln(os.Stderr, "Missing -r\n")
		flag.Usage()
		os.Exit(2)
	}

	findOpts.Recursive = true
	files, err := bindata.FindFiles(findOpts.Prefix, findOpts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	f, err := os.OpenFile(*output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
	defer f.Close()
	err = files.Generate(f, genOpts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}
}
