package main

import (
	"github.com/GeorgeSaussy/goldcmd"
	"github.com/GeorgeSaussy/goldcmd/examples/calculator/src"
)

func main() {
	cli := goldcmd.NewCli("latest", "A simple calculator CLI app.")
	cli.HandleSubcommand(src.Adder())
	cli.HandleSubcommand(src.Subtracter())
	cli.HandleSubcommand(src.Multiplier())
	cli.HandleSubcommand(src.Divider())
	cli.Run()
}
