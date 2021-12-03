package src

import (
	"fmt"

	"github.com/GeorgeSaussy/goldcmd"
)

// Get the subcommand handler for addition.
func Adder() *goldcmd.SubcommandHandler {
	ret, err := goldcmd.NewSubcommandHandler("add", "Add two numbers together.")
	if err != nil {
		panic(err)
	}
	if err := ret.AddIntArg([]string{"first", "f"}, "first integer argument"); err != nil {
		panic(err)
	}
	if err := ret.AddIntArg([]string{"second", "s"}, "second integer argument"); err != nil {
		panic(err)
	}
	ret.Example("with mixed arguments", "calculator add -f 1 -second 34", "35")
	ret.Example("again with flags", "calculator add -f=1 --second=34", "35")

	ret.Handle(func(handler *goldcmd.SubcommandHandler) {
		a, _ := handler.GetInt("f")
		b, _ := handler.GetInt("s")
		fmt.Printf("%d\n", a+b)
	})
	return ret
}
