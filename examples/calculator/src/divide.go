package src

import (
	"fmt"

	"github.com/GeorgeSaussy/goldcmd"
)

/// Get a subcommand handler for divide.
/// # Returns
/// A SubCmdHandler instance that can divide command line arguments
func Divider() *goldcmd.SubcommandHandler {
	ret, err := goldcmd.NewSubcommandHandler("divide", "Divide two numbers.")
	if err != nil {
		panic(err)
	}
	if err := ret.AddIntArg([]string{"first", "f"}, "first integer argument"); err != nil {
		panic(err)
	}
	if err := ret.AddIntArg([]string{"second", "s"}, "second integer argument"); err != nil {
		panic(err)
	}
	ret.Example("with mixed arguments", "calculator divide -f 1 -second 34", "0")
	ret.Example("again with flags", "calculator divide -f=4 --second=2", "2")
	ret.Handle(func(handler *goldcmd.SubcommandHandler) {
		a, _ := handler.GetInt("f")
		b, _ := handler.GetInt("s")
		fmt.Printf("%d\n", a/b)
	})
	return ret
}
