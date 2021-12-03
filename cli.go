package goldcmd

import (
	"fmt"
	"os"
)

// The CLI "server" object
type Cli struct {
	// Brief documentation of the CLI app
	documentation string
	// TODO(gs): The version is not publically exposed. Add 'version' subcommand
	// preset. Also check that 'help' is a reserved subcommand.
	// Version of the CLI app
	version string
	// The subcommands available with this CLI
	subcommands []*SubcommandHandler
}

// Create a new Cli instance.
//
// This Cli will do nothing when run is called unless it is mutated.
// The argument `version` should be a string describing the version of
// the CLI e.g. "2.1.3", "beta", "latest", &c
// The argument `doc` should contain brief documentation of the command.
func NewCli(version string, doc string) Cli {
	return Cli{version: version,
		documentation: doc, subcommands: make([]*SubcommandHandler, 0)}
}

// Add a subcommand to the CLI instance.
func (cli *Cli) HandleSubcommand(subcmd *SubcommandHandler) {
	cli.subcommands = append(cli.subcommands, subcmd)
}

// Print a help message.
// The variable sub is either one of the subcommands or an empty string.
// If it matches a a subcommand, then the help message for that subcommand
// is printed. If it is empty (or does not match a subcommand), then the
// help message for the CLI app is printed.
func (cli *Cli) printHelp(sub string) {
	if sub != "" {
		for _, subcmd := range cli.subcommands {
			if subcmd.name == sub {
				fmt.Printf("%s\n\n", subcmd.documentation)
				subcmd.printArgumentHelp()
				subcmd.printOptionHelp()
				subcmd.printExampleHelp()
				return
			}
		}
	}
	fmt.Printf("%s\n\n", cli.documentation)
	fmt.Printf("SUBCOMMANDS\n")
	for _, subcmd := range cli.subcommands {
		fmt.Printf("  %s\t%s\n", subcmd.name, subcmd.documentation)
	}
	fmt.Printf("  help\tthis help message\n\n")
	fmt.Printf("Get help with a subcommand with by passing it as an argument to the 'help' subcommand.\n")
}

// Either print help, or run a subcommand.
func (cli *Cli) Run() {
	if len(os.Args) == 1 {
		cli.printHelp("")
		return
	} else if len(os.Args) == 2 {
		cmd := os.Args[1]
		if cmd == "--help" || cmd == "help" || cmd == "-h" {
			cli.printHelp("")
			return
		} else {
			for _, subcmd := range cli.subcommands {
				if subcmd.name == cmd {
					// parse command line arguments and pass them to the subcommand
					subcmd.parseFlags(os.Args)
					// if the values are valid, then run the subcommand's handle function
					subcmd.handle(subcmd)
					return
				}
			}
			cli.printHelp("")
			return
		}
	} else {
		cmd := os.Args[1]
		if cmd == "--help" || cmd == "help" || cmd == "-h" {
			arg := os.Args[2]
			cli.printHelp(arg)
			return
		} else {
			for _, subcmd := range cli.subcommands {
				if subcmd.name == cmd {
					// parse command line arguments and pass them to the subcommand
					subcmd.parseFlags(os.Args)
					// if the values are valid, then run the subcommand's handle function
					subcmd.handle(subcmd)
					return
				}
			}
			cli.printHelp("")
			return
		}
	}
}
