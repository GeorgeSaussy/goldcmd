package goldcmd

import (
	"math"
	"testing"
)

func TestNewCommandParser(t *testing.T) {
	cp := newCommandParser()
	if len(cp.allLabels) != 0 {
		t.Fatalf("new parser should not have labels")
	}
	if len(cp.intLabels) != 0 {
		t.Fatalf("new parser should not have int labels declared")
	}
	if len(cp.intValues) != 0 {
		t.Fatalf("new parser should not have int values set")
	}
	if len(cp.menu) != 0 {
		t.Fatalf("new parser should not have menu elements set")
	}
}

func TestHasAlias(t *testing.T) {
	cp := newCommandParser()
	cp.addIntArg([]string{"a", "b"}, "some documentation")
	if !cp.hasAlias("a") {
		t.Fatalf("alias \"a\" should be present")
	}
	if !cp.hasAlias("b") {
		t.Fatalf("alias \"b\" should be present")
	}
	if cp.hasAlias("c") {
		t.Fatalf("alias \"c\" should not be present")
	}
}

func sampleCommandParser() *commandParser {
	cp := newCommandParser()
	cp.addIntArg([]string{"a", "b"}, "some int arg documentation")
	cp.addStrArg([]string{"c"}, "some string arg documentation")
	cp.addFloatArg([]string{"d", "e", "f"}, "some float arg documentation")
	cp.addBoolArg([]string{"g", "h"}, "some bool arg documentation")
	return cp
}

func TestAddArgs(t *testing.T) {
	cp := sampleCommandParser()
	for _, c := range "abcdefgh" {
		if !cp.hasAlias(string(c)) {
			t.Fatalf("alias '%c' should be present", c)
		}
	}
}

func TestSetArgs(t *testing.T) {
	cp := sampleCommandParser()
	cp.setIntArg([]string{"a", "b"}, 1337)
	i1, ok1 := cp.intValues["a"]
	i2, ok2 := cp.intValues["b"]
	if !ok1 || !ok2 || i1 != 1337 || i2 != 1337 {
		t.Fatalf("int value should be 1337 for aliases \"a\" and \"b\"")
	}
	cp.setStrArg([]string{"c"}, "hello world")
	s, ok := cp.strValues["c"]
	if !ok || s != "hello world" {
		t.Fatalf("string value should be \"hello world\" for alias \"c\"")
	}
	cp.setFloatArg([]string{"d", "e", "f"}, 0.42)
	f, ok := cp.floatValues["d"]
	if !ok || math.Abs(0.42-f) > 0.001 {
		t.Fatalf("float value should be 0.42 for alias \"d\", f=%b, ok=%t", f, ok)
	}
	f, ok = cp.floatValues["e"]
	if !ok || math.Abs(0.42-f) > 0.001 {
		t.Fatalf("float value should be 0.42 for alias \"e\"")
	}
	f, ok = cp.floatValues["f"]
	if !ok || math.Abs(0.42-f) > 0.001 {
		t.Fatalf("float value should be 0.42 for alias \"f\"")
	}
	cp.setBoolArg([]string{"g", "h"}, true)
	b1, ok1 := cp.boolValues["g"]
	b2, ok2 := cp.boolValues["h"]
	if !ok1 || !ok2 || !b1 || !b2 {
		t.Fatalf("bool value should be true for aliases \"g\" and \"f\"")
	}
}

func TestBasicParseFlags(t *testing.T) {
	argsSet := [][]string{
		{"cli", "sub", "-a", "1337", "-x", "42"},
		{"cli", "sub", "-a", "1337", "--x", "42"},
		{"cli", "sub", "-a", "1337", "-x=42"},
		{"cli", "sub", "-a", "1337", "--x=42"},

		{"cli", "sub", "-x", "42", "-a", "1337"},
		{"cli", "sub", "--x", "42", "-a", "1337"},
		{"cli", "sub", "-x=42", "-a", "1337"},
		{"cli", "sub", "--x=42", "-a", "1337"},

		{"cli", "sub", "--a", "1337", "-x", "42"},
		{"cli", "sub", "--a", "1337", "--x", "42"},
		{"cli", "sub", "--a", "1337", "-x=42"},
		{"cli", "sub", "--a", "1337", "--x=42"},

		{"cli", "sub", "-x", "42", "--a", "1337"},
		{"cli", "sub", "--x", "42", "--a", "1337"},
		{"cli", "sub", "-x=42", "--a", "1337"},
		{"cli", "sub", "--x=42", "--a", "1337"},

		{"cli", "sub", "-a=1337", "-x", "42"},
		{"cli", "sub", "-a=1337", "--x", "42"},
		{"cli", "sub", "-a=1337", "-x=42"},
		{"cli", "sub", "-a=1337", "--x=42"},

		{"cli", "sub", "-x", "42", "-a=1337"},
		{"cli", "sub", "--x", "42", "-a=1337"},
		{"cli", "sub", "-x=42", "-a=1337"},
		{"cli", "sub", "--x=42", "-a=1337"},

		{"cli", "sub", "--a=1337", "-x", "42"},
		{"cli", "sub", "--a=1337", "--x", "42"},
		{"cli", "sub", "--a=1337", "-x=42"},
		{"cli", "sub", "--a=1337", "--x=42"},

		{"cli", "sub", "-x", "42", "--a=1337"},
		{"cli", "sub", "--x", "42", "--a=1337"},
		{"cli", "sub", "-x=42", "--a=1337"},
		{"cli", "sub", "--x=42", "--a=1337"},
	}
	for _, args := range argsSet {
		cp := sampleCommandParser()
		cp.addIntArg([]string{"x"}, "another int arg")
		cp.parseFlags(args)
		i1, ok1 := cp.intValues["a"]
		i2, ok2 := cp.intValues["b"]
		if !ok1 || !ok2 || i1 != 1337 || i2 != 1337 {
			t.Fatalf("int value should be 1337 for aliases \"a\" and \"b\"")
		}
		i, ok := cp.intValues["x"]
		if !ok || i != 42 {
			t.Fatalf("int value should be 42 for alias \"x\"")
		}
	}
}
