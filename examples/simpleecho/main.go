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
	if err = ret.AddStrArg([]string{"s", "text"}, "A string to echo"); err != nil {
		panic(err)
	}
	ret.Example("Echo a string", "simpleecho echo -s=example_string", "example_string")
	ret.Example("Echo a string", "simpleecho echo --s \"example string\"", "example string")
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
	cli := goldcmd.NewCli("latest", "A simple echo command line tool.")
	cli.HandleSubcommand(handler())
	cli.Run()
}
