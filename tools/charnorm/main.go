package main

import (
	"bufio"
	"flag"
	"fmt"
	"mojo/util"
	"os"
)

// Normalize characters in a file

func main() {
	flag.Parse()
	var err error
	var f *os.File
	switch flag.NArg() {
	case 0:
		f = os.Stdin
	case 1:
		f, err = os.Open(flag.Arg(0))
		if err != nil {
			fmt.Printf("error opening %s: %s\n", flag.Arg(0), err.Error())
			os.Exit(1)
		}
	default:
		fmt.Printf("input must be from stdin or file\n")
		os.Exit(1)
	}

	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		line = util.NonASCIIToASCII(line)
		fmt.Println(line)
	}
}
