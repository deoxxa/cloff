package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	file = flag.String("file", "", "File to read.")
	line = flag.Int("line", 0, "Line to compute offset of.")
	col  = flag.Int("col", 0, "Column in line to compute offset of.")
)

func main() {
	flag.Parse()

	if *file == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var l, c, t int

	b := make([]byte, 1024)
	for {
		n, err := f.Read(b)

		for _, r := range string(b[0:n]) {
			if l == *line-1 && c == *col-1 {
				goto done
			}

			if r == '\n' {
				if l == *line-1 {
					fmt.Fprintf(os.Stderr, "ERROR: file %q line %d was only %d columns wide; you requested column %d\n", *file, *line, c, *col)
					os.Exit(1)
				}

				l++
				c = 0
			} else {
				c++
			}

			t++
		}

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	fmt.Fprintf(os.Stderr, "ERROR: file %q could only be read to %d,%d; you requested %d,%d\n", *file, l, c, *line, *col)
	os.Exit(1)

done:
	fmt.Printf("%d", t)
}
