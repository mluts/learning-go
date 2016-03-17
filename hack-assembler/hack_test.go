package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestStripComment(t *testing.T) {
	examples := map[string]string{
		"abcd//abcd":  "abcd",
		"abcd //abcd": "abcd ",
		"//abcd":      "",
	}

	for example, expectedResult := range examples {
		result := stripComment(example)
		if result != expectedResult {
			t.Errorf("Have \"%s\" from \"%s\", but should have \"%s\"", result, example, expectedResult)
		}
	}
}

func TestStripWhitespace(t *testing.T) {
	example := "  abc  \t   \rqwerty  "
	expected := "abcqwerty"
	result := stripWhitespace(example)

	if result != expected {
		t.Errorf("Expected \"%s\", but have \"%s\"", expected, result)
	}
}

func TestIsAinstruction(t *testing.T) {
	res := isAinstruction("@a")

	if res != true {
		t.Error("@a should be A-instruction, but it is not")
	}

	res = isAinstruction("D=A")

	if res == true {
		t.Error("D=A should not be A-instruction")
	}
}

func TestIsLabel(t *testing.T) {
	res := isLabel("(ABC)")

	if res != true {
		t.Error("(ABC) should be recognised as label")
	}

	res = isLabel("A=D")

	if res == true {
		t.Error("A=D should not be recognised as label")
	}
}

func TestParseAinstruction(t *testing.T) {
	tokens := parseLine("@a")

	switch {
	case len(tokens) != 1:
		t.Fatalf("Wrong tokens size: %d", len(tokens))

	case tokens[0].t != T_AINST:
		t.Fatalf("Expected to have AInstruction, but have %d", tokens[0].t)

	case tokens[0].val != "a":
		t.Fatal("Expected first token value to be \"a\"")
	}

}

func TestParseLabel(t *testing.T) {
	tokens := parseLine("(ABC)")

	switch {
	case len(tokens) != 1:
		t.Fatalf("Wrong tokens size: %d", len(tokens))

	case tokens[0].val != "ABC":
		t.Fatalf("Token should eq ABC, but have: %s", tokens[0].val)
	}

}

func TestParseDestComp(t *testing.T) {
	tokens := parseLine("A=D")

	switch {

	case len(tokens) != 2:
		t.Fatalf("Expected 3 tokens, but have: %d", len(tokens))

	case tokens[0].t != T_DEST:
		t.Fatal("First token type should be a T_DEST")

	case tokens[0].val != "A":
		t.Fatal("First token value should be \"A\"")

	case tokens[1].t != T_COMP:
		t.Fatal("Third token type should be a T_COMP")

	case tokens[1].val != "D":
		t.Fatal("Third token value should be \"A\"")
	}
}

func TestParseComp(t *testing.T) {
	tokens := parseLine("M")

	switch {
	case len(tokens) != 1:
		t.Fatalf("Expected 1 token, but have: %d", len(tokens))

	case tokens[0].t != T_COMP:
		t.Fatal("Expected first token to be T_COMP")

	case tokens[0].val != "M":
		t.Fatal("Expected first token val to be \"M\"")
	}
}

func TestParseMinusComp(t *testing.T) {
	tokens := parseLine("-M")

	switch {
	case len(tokens) != 1:
		t.Fatalf("Expected 2 tokens, but have: %d", len(tokens))

	case tokens[0].t != T_COMP:
		t.Fatal("Expected first token to be T_COMP")

	case tokens[0].val != "-M":
		t.Fatal("Expected first token value to be \"-M\"")
	}
}

func TestParseCompPlusComp(t *testing.T) {
	tokens := parseLine("M+A")

	switch {
	case len(tokens) != 1:
		t.Fatalf("Expected 3 tokens, but have: %d", len(tokens))

	case tokens[0].t != T_COMP:
		t.Fatal("Expected first token to be T_COMP")

	case tokens[0].val != "M+A":
		t.Fatal("Expected first token val to eq \"M+A\"")
	}
}

func TestParseOneComp(t *testing.T) {
	tokens := parseLine("1")

	switch {
	case len(tokens) != 1:
		t.Fatal("Expected to have 1 token")

	case tokens[0].t != T_COMP:
		t.Fatal("Expected first token to be T_COMP")

	case tokens[0].val != "1":
		t.Fatal("Expected first token value to eq \"1\"")
	}
}

func TestParseZeroComp(t *testing.T) {
	tokens := parseLine("0")

	switch {
	case len(tokens) != 1:
		t.Fatal("Expected to have one token")

	case tokens[0].t != T_COMP:
		t.Fatal("Expected first token to be T_COMP")

	case tokens[0].val != "0":
		t.Fatal("Expected first token value to be \"0\"")
	}
}

func TestParseJMP(t *testing.T) {
	tokens := parseLine("D=0;JMP")

	switch {
	case len(tokens) != 3:
		t.Fatal("Expected to have 3 tokens")
	case tokens[0].t != T_DEST:
		t.Fatal("Expected first token to be T_DEST")
	case tokens[0].val != "D":
		t.Fatal("Expected first token value to eq \"D\"")
	case tokens[1].t != T_COMP:
		t.Fatal("Expected second token to be T_COMP")
	case tokens[1].val != "0":
		t.Fatal("Expected second token value to eq \"0\"")
	case tokens[2].t != T_JMP:
		t.Fatal("Expected third token to be T_JMP")
	case tokens[2].val != "JMP":
		t.Fatal("Expected third token value to eq \"JMP\"")
	}
}

func TestParseEmptyLine(t *testing.T) {
	if parseLine("") != nil {
		t.Fatalf("Empty line should be parsed as nil")
	}
}

func TestParseLines(t *testing.T) {
	lines, symbols := parseLines(strings.NewReader("(A)\n@A\nD;JMP\nAM=D+1;JLE"))

	switch {
	case len(lines) != 3:
		t.Fatalf("Size of lines should be 2, but have: %d", len(lines))
	case symbols["A"] != 0:
		t.Fatal("Symbol A should eq 0")
	}

	for k, v := range defaultSymbolTable {
		if symbols[k] != v {
			t.Errorf("Symbol \"%s\" should be defined as %d, have: %v", k, v, symbols[k])
		}
	}
}

func TestSymbolToAddr(t *testing.T) {
	table := make(SymbolTable)
	for k, v := range defaultSymbolTable {
		table[k] = v
	}

	res := symbolToAddr("i", table)

	if res != 16 {
		t.Fatal("\"i\" symbol should be %d, but have %d", 16, res)
	}

	if table["i"] != 16 {
		t.Fatal("table should contain i symbol")
	}

	res = symbolToAddr("j", table)

	if res != 17 {
		t.Fatal("\"j\" symbol should be %d, but have %d", 17, res)
	}

	if table["j"] != 17 {
		t.Fatal("table should contain i symbol")
	}
}

func TestIsAddr(t *testing.T) {
	switch {
	case !isAddr("1234"):
		t.Error("1234 should be index")
	case isAddr("123a4"):
		t.Error("123a4 should not be index")
	}
}

func TestCompileAInstruction(t *testing.T) {
	res := compileLine([]Token{Token{T_AINST, fmt.Sprintf("%d", 0x7fff)}}, defaultSymbolTable)
	switch {
	case res != 0x7fff:
		t.Errorf("@32767 should be 0111111111111111, but have: %b", res)
	}
}
