package goldcmd

import (
	"strings"
	"testing"
)

func TestSubcommandExample(t *testing.T) {
	sc := subcommandExample{
		documentation: "Grep can be used to search for text in a file.",
		command:       "grep -r \"hello world\" .",
		output:        "",
	}
	m := sc.helpMessage()
	if !strings.Contains(m, sc.documentation) {
		t.Fail()
	}
	if !strings.Contains(m, sc.command) {
		t.Fail()
	}
	if !strings.Contains(m, "#") {
		t.Fail()
	}
}

func TestNoHash(t *testing.T) {
	sc := subcommandExample{
		documentation: "",
		command:       "grep -r \"hello world\" .",
		output:        "",
	}
	m := sc.helpMessage()
	if strings.Contains(m, "#") {
		t.Fail()
	}
}
