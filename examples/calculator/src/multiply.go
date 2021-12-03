package src

import (
	"fmt"

	"github.com/GeorgeSaussy/goldcmd"
)

// Get a subcommand handler for multiplication.
func Multiplier() *goldcmd.SubcommandHandler {
	ret, err := goldcmd.NewSubcommandHandler("multiply", "Multiply two numbers together.")
	if err != nil {
		panic(err)
	}
	if err := ret.AddIntArg([]string{"first", "f"}, "first integer argument"); err != nil {
		panic(err)
	}
	if err := ret.AddIntArg([]string{"second", "s"}, "second integer argument"); err != nil {
		panic(err)
	}
	ret.Example("with mixed arguments", "calculator multiply -f 1 -second 34", "34")
	ret.Example("again with flags", "calculator multiply -f=1 --second=34", "34")
	ret.Handle(func(handler *goldcmd.SubcommandHandler) {
		a, _ := handler.GetInt("f")
		b, _ := handler.GetInt("s")
		fmt.Printf("%d\n", a*b)
	})
	return ret
}
