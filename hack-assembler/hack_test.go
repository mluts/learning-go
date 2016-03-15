package main

import "testing"

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
	case len(tokens) != 2:
		t.Fatalf("Wrong tokens size: %d", len(tokens))

	case tokens[0].t != T_AINST:
		t.Fatalf("Expected to have AInstruction, but have %d", tokens[0].t)

	case tokens[1].t != T_ID:
		t.Fatalf("Expected to have ID, but have %d", tokens[1].t)

	case tokens[1].val != "a":
		t.Fatalf("Second token should be \"a\"")
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

func TestParseCinstruction(t *testing.T) {
	tokens := parseLine("A=D")

	switch {
	case len(tokens) != 3:
		t.Fatalf("Expected 3 tokens, but have %d", len(tokens))
	case tokens[0].t != T_ID:
		t.Fatalf("Expected 1st token to be an id")
	case tokens[1].t != T_EQ:
		t.Fatalf("Expected 2nd token to be an equation sign")
	case tokens[2].t != T_ID:
		t.Fatalf("Expected 3rd token to be an id")
	case tokens[0].val != "A":
		t.Fatalf("Expected 1st token to eq %s, but have %s", "A", tokens[0].val)
	case tokens[2].val != "D":
		t.Fatalf("Expected 1st token to eq %s, but have %s", "D", tokens[1].val)
	}
}
