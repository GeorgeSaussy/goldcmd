package goldcmd

import (
	"fmt"
)

// A subcommandExample stores an example of a subcommand in use.
type subcommandExample struct {
	// documentation for the example e.g. context or an explanation
	documentation string
	// the text of the command e.g. 'grep "hello world" -r ~'
	command string
	// the expected output of the command, which can be fabricated as an example
	output string
}

// Print the help message associated with the command example.
func (ex *subcommandExample) helpMessage() string {
	ret := ""
	if len(ex.documentation) > 0 {
		ret = ret + fmt.Sprintf("$ # %s\n", ex.documentation)
	}
	ret = ret + fmt.Sprintf("$ %s\n", ex.command)
	if len(ex.output) > 0 {
		ret = ret + fmt.Sprintf("%s\n\n", ex.output)
	}
	return ret
}
