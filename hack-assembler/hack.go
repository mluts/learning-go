package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	COMMENT = "//"
	T_AINST = iota
	T_EQ
	T_SEMIC
	T_ID
	T_LABEL
)

type Token struct {
	t   int
	val string
}

func stripComment(line string) string {
	return strings.Split(line, COMMENT)[0]
}

func stripWhitespace(line string) (result string) {
	result = strings.Replace(line, " ", "", -1)
	result = strings.Replace(result, "\t", "", -1)
	result = strings.Replace(result, "\r", "", -1)
	return
}

func isAinstruction(line string) bool {
	return strings.HasPrefix(line, "@")
}

func isLabel(line string) bool {
	return line[0] == '(' && line[len(line)-1] == ')'
}

func parseLine(line string) []Token {
	switch {
	case isAinstruction(line):
		return []Token{Token{T_AINST, ""}, Token{T_ID, line[1:]}}
	case isLabel(line):
		return []Token{Token{T_LABEL, line[1 : len(line)-1]}}
	default:
		tokens := make([]Token, 0)

		for _, ch := range line {
			switch {
			case 'A' <= ch && ch <= 'Z':
				tokens = append(tokens, Token{T_ID, string(ch)})
			case ch == '=':
				tokens = append(tokens, Token{T_EQ, string(ch)})
			case ch == ';':
				tokens = append(tokens, Token{T_SEMIC, string(ch)})
			}
		}
		return tokens
	}
}

func showUsage() {
	fmt.Printf(`
	USAGE:

	%s ASSEMBLY-FILE OUTPUT-FILE

	Compiles ASSEMBLY-FILE to HACK machine binary
`, os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 {
		showUsage()
	}

	input, output := os.Args[1], os.Args[2]

	inputFile, err := os.OpenFile(input, os.O_RDONLY, 0666)

	if err != nil {
		fmt.Printf("Can't open file for reading %s: %v", input, err)
		showUsage()
	}

	outputFile, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)

	if err != nil {
		fmt.Printf("Can't open file for writing %s: %v", output, err)
		showUsage()
	}

	_, err = io.Copy(outputFile, inputFile)

	if err != nil {
		fmt.Printf("Can't copy from %s to %s: %v", output, input, err)
		showUsage()
	}
}
