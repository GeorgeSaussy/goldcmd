package src

import (
	"fmt"

	"github.com/GeorgeSaussy/goldcmd"
)

// Get a subcommand handler for subtraction.
func Subtracter() *goldcmd.SubcommandHandler {
	ret, err := goldcmd.NewSubcommandHandler("subtract", "Subtract two numbers.")
	if err != nil {
		panic(err)
	}
	if err := ret.AddIntArg([]string{"first", "f"}, "first integer argument"); err != nil {
		panic(err)
	}
	if err := ret.AddIntArg([]string{"second", "s"}, "second integer argument"); err != nil {
		panic(err)
	}
	ret.Example("with mixed arguments", "calculator subtract -f 1 -second 34", "-33")
	ret.Example("again with flags", "calculator subtract -f=1 --second=34", "-33")
	ret.Handle(func(handler *goldcmd.SubcommandHandler) {
		a, _ := handler.GetInt("f")
		b, _ := handler.GetInt("s")
		fmt.Printf("%d\n", a-b)
	})
	return ret
}
