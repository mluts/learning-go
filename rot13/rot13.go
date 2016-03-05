package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func usage(msg string) {
	if len(msg) > 0 {
		fmt.Printf(msg)
	}

	fmt.Printf(`
	USAGE:

	%s FILE...
	%s -

	Rotates ASCII chars and prints them to stdout
`, os.Args[0], os.Args[0])
	os.Exit(1)
}

func rotateCh(ch rune) (rch rune) {
	switch {
	case 'a' <= ch && ch <= 'z':
		return ((ch-'a')+13)%26 + 'a'
	case 'A' <= ch && ch <= 'Z':
		return ((ch-'A')+13)%26 + 'A'
	default:
		return ch
	}
}

type RuneReader interface {
	ReadRune() (rune, int, error)
}

func printRotated(reader RuneReader) {
	for {
		r, _, err := reader.ReadRune()
		switch err {
		case nil:
			fmt.Printf("%c", rotateCh(r))
		case io.EOF:
			return
		default:
			panic(err)
		}
	}
}

func printPathRotated(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Can't read file \"%s\": %v\n", path, err)
	} else {
		printRotated(bufio.NewReader(file))
	}
}

func main() {
	if len(os.Args) < 2 {
		usage("Wrong number of arguments!")
	} else if os.Args[1] == "-" {
		printRotated(bufio.NewReader(os.Stdin))
	} else {
		for _, path := range os.Args[1:] {
			printPathRotated(path)
		}
	}
}
