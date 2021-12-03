# README

This is a small library for implementing command line utilities.
This was originally implemented as part of the larger [geode](https://github.com/GeorgeSaussy/geode) project.

## Quick Start 

This script can be found in `goldcmd/examples/simpleecho/main.go`.
It implements a simple version of the UNIX `echo` command.

```golang
package main

import (
	"fmt"

	"github.com/GeorgeSaussy/goldcmd"
)

/// Get a subcommand handler to echo a string.
func handler() *goldcmd.SubcommandHandler {
	ret, err := goldcmd.NewSubcommandHandler("echo", "Echo a string.")
	if err != nil {
		panic(err)
	}
	if err = ret.AddStrArg([]string{"s"}, "A string to echo"); err != nil {
		panic(err)
	}
	ret.Example("Echo a string", "dumbecho echo -s=example_string", "example_string")
	ret.Example("Echo a string", "dumbecho echo --s \"example string\"", "example string")
	ret.Handle(func(handler *goldcmd.SubcommandHandler) {
		if a, err := handler.GetStr("s"); err != nil {
			fmt.Printf("No string found!\n")
		} else {
			fmt.Printf("%s\n", a)
		}
	})
	return ret
}

func main() {
	cli := goldcmd.NewCli("latest", "A dumb echo app.")
	cli.HandleSubcommand(*handler())
	cli.Run()
}

```

If you compile this script to `./a.out`, it could then be called like this:


- Example: No arguments

```
% ./a.out 
A simple echo command line tool.

SUBCOMMANDS
  echo  Echo a string.
  help  this help message

Get help with a subcommand with by passing it as an argument to the 'help' subcommand.
```

The same output would be printed for `./a.out -h` and `./a.out --help`.

- Example: Help with the echo subcommand

```
% ./a.out help echo
Echo a string.

ARGUMENTS
 --s, --text    A string to echo


EXAMPLES
$ # Echo a string
$ simpleecho echo -s=example_string
example_string

$ # Echo a string
$ simpleecho echo --s "example string"
example string
```

Note, everything after the first line is the output of the command.

- Example: Echoing things

```
% ./a.out echo -s "Hello, World!"
Hello, World!
```

```
% ./a.out echo --text="Bonsoir, Elliot."
Bonsoir, Elliot.
```