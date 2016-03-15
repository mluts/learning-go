package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	COMMENT   = "//"
	VAR       = "_var"
	A         = "@"
	LEFT_PAR  = '('
	RIGHT_PAR = ')'
	A_REG     = "A"
	D_REG     = "D"
	M_REG     = "M"
	T_AINST   = iota
	T_EQ
	T_SEMIC
	T_ID
	T_LABEL
)

type SymbolTable map[string]uint16

var defaultSymbolTable = SymbolTable{
	"R0":   0,
	"R1":   1,
	"R2":   2,
	"R3":   3,
	"R4":   4,
	"R5":   5,
	"R6":   6,
	"R7":   7,
	"R8":   8,
	"R9":   9,
	"R10":  10,
	"R11":  11,
	"R12":  12,
	"R13":  13,
	"R14":  14,
	"R15":  15,
	"SP":   0,
	"LCL":  1,
	"ARG":  2,
	"THIS": 3,
	"THAT": 4,

	"SCREEN": 0x4000,
	"KBD":    0x6000,

	"_var": 16,
}

type Token struct {
	t   uint16
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
	return strings.HasPrefix(line, A)
}

func isLabel(line string) bool {
	return line[0] == LEFT_PAR && line[len(line)-1] == RIGHT_PAR
}

func parseLine(line string) []Token {
	if len(line) == 0 {
		return nil
	}

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
			default:
				return nil
			}
		}
		return tokens
	}
}

func parseLines(r io.Reader) (lines [][]Token, symbols map[string]uint16) {
	scanner := bufio.NewScanner(r)

	symbols = make(map[string]uint16)

	for k, v := range defaultSymbolTable {
		symbols[k] = v
	}

	labels := make([][]Token, 0)
	lineIndex := uint16(0)

	for scanner.Scan() {
		line := parseLine(scanner.Text())

		if line[0].t == T_LABEL {
			labels = append(labels, line)
		} else {
			for _, label := range labels {
				symbols[label[0].val] = lineIndex
			}
			labels = make([][]Token, 0)
			lines = append(lines, line)
			lineIndex++
		}
	}

	if scanner.Err() != nil {
		panic("Can't parse source!")
	}

	return
}

func symbolToAddr(symbol string, symbols SymbolTable) uint16 {
	_, ok := symbols[symbol]

	if !ok {
		symbols[symbol] = symbols[VAR]
		symbols[VAR]++
	}

	return symbols[symbol]
}

func isAddr(str string) bool {
	for _, ch := range str {
		if ch > '9' || ch < '0' {
			return false
		}
	}

	return true
}

func compileAinstruction(line []Token, symbols SymbolTable) (i uint16) {
	var addr uint16
	if isAddr(line[0].val) {
		res, err := strconv.ParseUint(line[0].val, 10, 16)
		addr = uint16(res)
		if err != nil {
			panic("Can't parse number!")
		}
	} else {
		addr = symbolToAddr(line[0].val, symbols)
	}
	return addr &^ uint16(1<<15)
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func getDestRegisters(line []Token) ([]string, int) {
	dest := -1
	for i, t := range line {
		if t.t == T_EQ {
			dest = i
			break
		}
	}

	if dest > 0 {
		destRegisters = make([]string)
		for _, t := range line[0:dest] {
			switch t.val {
			case A_REG, D_REG, M_REG:
				if contains(destRegisters, t.val) {
					panic("Duplicated A register")
				} else {
					destRegisters = append(destRegist)
				}
			default:
				panic(fmt.Sprintf("Unknown register: %s", t.val))
			}
		}
		return destRegisters, dest
	} else {
		return nil, 0
	}
}

func getCompRegisters(line []Token) []string {
	compRegisters := make([]string, 0)
	for i, t := range line {
		if t.val == T_SEMIC {
			break
		}

		switch t.val {
		case A_REG:
			if contains(compRegisters, A_REG) {
				panic("Double A comparison register")
			}
		}
	}
	return compRegisters
}

func compileCinstruction(line []Token) (i uint16) {
	destRegisters, i := getDestRegisters(line)
	if i > 0 {
		line = line[i+1:]
	}
	compRegisters := getCompRegisters(line)
	return 0
}

func compileLine(line []Token, symbols SymbolTable) (i uint16) {
	if line[0].t == T_AINST {
		return compileAinstruction(line, symbols)
	} else {
		return compileCinstruction(line)
	}
}

// func compile(r io.Reader) string {
// 	lines, symbols := parseLines(r)
// }

func showUsage() {
	fmt.Printf(`
	USAGE:

	%s ASSEMBLY-FILE OUTPUT-FILE

	Compiles HACK-ASSEMBLY to HACK machine code
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
