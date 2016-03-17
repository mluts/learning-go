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
	A_REG     = 'A'
	D_REG     = 'D'
	M_REG     = 'M'
	PLUS      = '+'
	MINUS     = '-'
	AND       = '&'
	OR        = '|'
	ONE       = '1'
	ZERO      = '0'
	NEG       = '!'

	C_INST_MASK = 0xE000

	ZX     = 0x800
	NX     = 0x400
	ZY     = 0x200
	NY     = 0x100
	F      = 0x80
	NO     = 0x40
	A_COMP = 0x1000

	T_AINST = iota
	T_DEST
	T_COMP
	T_JMP
	T_LABEL

	STATE_START
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

func isAddr(str string) bool {
	for _, ch := range str {
		if ch > '9' || ch < '0' {
			return false
		}
	}

	return true
}

func parseCInstruction(line string) []Token {
	tokens := []Token{}
	for _, str := range []string{"=", ";"} {
		if strings.Count(line, str) > 1 {
			panic(fmt.Sprintf("Only one \"%s\" allowed", str))
		}
	}

	var (
		dest, comp, jmp string
	)

	destComp := strings.Split(line, "=")

	if len(destComp) > 1 {
		dest = destComp[0]
		comp = destComp[1]
	} else {
		comp = destComp[0]
	}

	compJmp := strings.Split(comp, ";")

	if len(compJmp) > 1 {
		comp = compJmp[0]
		jmp = compJmp[1]
	}

	if dest != "" {
		tokens = append(tokens, Token{T_DEST, dest})
	}

	if comp == "" {
		panic("comp can't be nil!")
	} else {
		tokens = append(tokens, Token{T_COMP, comp})
	}

	if jmp != "" {
		tokens = append(tokens, Token{T_JMP, jmp})
	}

	return tokens
}

func parseLine(line string) []Token {
	line = stripComment(line)
	line = stripWhitespace(line)

	if len(line) == 0 {
		return nil
	}

	switch {
	case isAinstruction(line):
		return []Token{Token{T_AINST, line[1:]}}
	case isLabel(line):
		return []Token{Token{T_LABEL, line[1 : len(line)-1]}}
	default:
		return parseCInstruction(line)
	}
}

func parseLines(r io.Reader) (lines [][]Token, symbols map[string]uint16) {
	scanner := bufio.NewScanner(r)

	symbols = SymbolTable{}

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

func compileDest(t Token) (mask uint16) {
	for i, ch := range t.val {
		switch ch {
		case A_REG, D_REG, M_REG:
			if i+1 < len(t.val) && strings.ContainsRune(t.val[i+1:], ch) {
				panic(fmt.Sprintf("Duplicated dest register: %c", ch))
			}
			fallthrough
		case A_REG:
			mask |= A_DEST_MASK
		case M_REG:
			mask |= M_DEST_MASK
		case D_REG:
			mask |= D_DEST_MASK
		default:
			panic(fmt.Sprintf("Unknown register: %c", ch))
		}
	}

	return
}

func compileComp1(ch rune) (mask uint16) {
	switch ch {
	case ZERO:
		return ZX | ZY | F
	case ONE:
		return ZX | NX | ZY | NY | F | NO
	case D_REG:
		return ZY | NY
	case A_REG:
		return ZX | NY
	case M_REG:
		return ZX | NY | A_COMP
	default:
		panic(fmt.Sprintf("Unexpected comp string: \"%c\"", ch))
	}
}

func compileComp2(operator rune, operand rune) (mask uint16) {
	switch operand {
	case A_REG:
		mask |= ZX | NX
	case M_REG:
		mask |= ZX | NX | A_COMP
	case D_REG:
		mask |= ZY | NY
	case ONE:
		mask |= ZX | NX | ZY | F
	default:
		panic(fmt.Sprintf("Unknown operand: %c", operand))
	}

	switch operator {
	case MINUS:
		mask |= F | NO
	case NEG:
		mask |= NO
	default:
		panic(fmt.Sprintf("Unexpected operator: %c", operator))
	}

	return
}

func compileComp3(operand1 rune, operator rune, operand2 rune) (mask uint16) {
	if operand1 == operand2 {
		panic("Equal operands not allowed")
	}

	if (operand1 == A_REG || operand1 == M_REG) &&
		(operand2 == A_REG || operand2 == M_REG) {
		panic("Cant operate on A and M simultaneously")
	}

	if !(operand1 == A_REG || operand1 == M_REG || operand1 == D_REG) {
		panic(fmt.Sprintf("Unknown register: %c", operand1))
	}

	if !(operand2 == A_REG || operand2 == M_REG || operand2 == D_REG) {
		panic(fmt.Sprintf("Unknown register: %c", operand2))
	}

	if operand1 == M_REG || operand2 == M_REG {
		mask |= A_COMP
	}

	if operator != MINUS {
		if operand2 == 'D' {
			operand1, operand2 = operand2, operand1
		}
	}

	switch operator {
	case OR:
		mask |= NX | NY | NO
	case AND:
		switch operand1 {
		case D_REG:

		}
	case MINUS:
	case PLUS:
	}

	return
}

func compileComp(t Token) uint16 {
	switch len(t.val) {
	case 1:
		return compileComp1(t)
	case 2:
		return compileComp2(t.val[0], t.val[1])
	case 3:
		return compileComp3(t.val[0], t.val[1], t.val[2])
	default:
		panic(fmt.Sprintf("Don't know how to handle comp \"%s\"", t.val))
	}
}

func compileJmp(t Token) (mask uint16) {
}

func compileCinstruction(line []Token) (i uint16) {
	i |= C_INST_MASK

	for _, t := range line {
		switch t.t {
		case T_DEST:
			i |= compileDest(t)
		case T_COMP:
			i |= compileComp(t)
		case T_JMP:
			i |= compileJmp(t)
		default:
			panic("Unknown token type!")
		}
	}

	return i
}

func compileLine(line []Token, symbols SymbolTable) uint16 {
	if line[0].t == T_AINST {
		return compileAinstruction(line, symbols)
	} else {
		return compileCinstruction(line)
	}
}

func compile(r io.Reader) (code []uint16) {
	lines, symbols := parseLines(r)

	for _, l := range lines {
		code = append(code, compileLine(l, symbols))
	}

	return
}

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
