package main

import (
	"bufio"
	"flag"
	"fmt"
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
		// line = util.NonASCIIToASCII(line)
		b := []byte(line)
		lenb := len(b)
		lenbm1 := lenb - 1
		c := make([]byte, 0)
		for i := 0; i < len(b); i++ {
			if b[i] == 0xc3 && i < lenbm1 {
				c = append(c, 'a')
				i++
			} else {
				c = append(c, b[i])
			}
		}
		fmt.Println(string(c))
	}
}
