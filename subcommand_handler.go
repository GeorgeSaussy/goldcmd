package goldcmd

import (
	"errors"
	"fmt"
	"unicode"
)

// Handle a CLI subcommand.
//
// For each row in ints, each element in the row can be used internally to
// refer to the value given at the command line. For example, if ints is
// [["a", "apple"]], then GetInt("a") should equal GetInt("apple") no matter
// what the user wrote at the command line.
type SubcommandHandler struct {
	// subcommand name, functions as the label of the subcomamnd
	name string
	// documentation for the subcommand
	documentation string
	// the actual function to handle the execution of the function
	handle func(h *SubcommandHandler)
	// examples of the subcommand in use
	examples []subcommandExample
	// the argument parsing handler
	argparser *commandParser
	// the param parsing handler
	paramparser *commandParser
}

// Check that a flag name is valid.
// The name is valid if it is composed of letters, numbers, hyphens, and
// underscores and it begins with a letter. Otherwise, the name is invalid.
func aliasIsValid(name string) bool {
	if len(name) == 0 {
		return false
	}
	for i, c := range name {
		if i == 0 {
			if !unicode.IsLetter(c) {
				return false
			} else {
				if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' && c != '-' {
					return false
				}
			}
		}
	}
	return true
}

// Create a new SubcommandHandler instance.
// An error may occur if the subcommand labels are not valid e.g. if the one
// of the labels had a leading '-'.
func NewSubcommandHandler(name string, doc string) (*SubcommandHandler, error) {
	if !aliasIsValid(name) {
		return nil, errors.New("flag name \"" + name + "\" is not valid")
	}
	return &SubcommandHandler{
		name:          name,
		documentation: doc,
		handle:        func(h *SubcommandHandler) {},
		examples:      make([]subcommandExample, 0),
		argparser:     newCommandParser(),
		paramparser:   newCommandParser(),
	}, nil
}

// Check if a set of aliases are valid and not already in use.
// The function returns true if any of the aliases can be used.
func (h *SubcommandHandler) checkAliasesAllowed(aliases []string) bool {
	for _, alias := range aliases {
		if !aliasIsValid(alias) {
			return false
		}
	}
	// check that the label is not in use
	for _, alias := range aliases {
		if h.argparser.hasAlias(alias) || h.paramparser.hasAlias(alias) {
			return false
		}
	}
	return true
}

// Add an integer argument to the subcommand.
// Warning, the subcommand will always fail if the argument is not set.
// If one of the aliases is used by another argument or
// parameter, the function will return an error and handler will not be mutated.
func (h *SubcommandHandler) AddIntArg(aliases []string, doc string) error {
	// TODO(gs): This can get slow for a large number of arguments. It may be
	// better to re-implement this with a hash map to support more complex
	// subcommands. The same goes for the float, string, and param versions of
	// this command

	if !h.checkAliasesAllowed(aliases) {
		return errors.New("invalid label value")
	}
	h.argparser.addIntArg(aliases, doc)
	return nil
}

// Add a float argument to the subcommand.
// Warning, the subcommand will always fail if the argument is not set.
// If one of the aliases is used by another argument or
// parameter, the function will return an error and handler will not be mutated.
func (h *SubcommandHandler) AddFloatArg(aliases []string, doc string) error {
	if !h.checkAliasesAllowed(aliases) {
		return errors.New("invalid label value")
	}
	h.argparser.addFloatArg(aliases, doc)
	return nil
}

// Add a string argument to the subcommand.
// Warning, the subcommand will always fail if the argument is not set.
// If one of the aliases is used by another argument or
// parameter, the function will return an error and handler will not be mutated.
func (h *SubcommandHandler) AddStrArg(aliases []string, doc string) error {
	if !h.checkAliasesAllowed(aliases) {
		return errors.New("invalid label value")
	}
	h.argparser.addStrArg(aliases, doc)
	return nil
}

// Add a boolean argument to the subcommand.
// Warning, the subcommand will always fail if the argument is not set.
// If one of the aliases is used by another argument or
// parameter, the function will return an error and handler will not be mutated.
func (h *SubcommandHandler) AddBoolArg(aliases []string, doc string) error {
	if !h.checkAliasesAllowed(aliases) {
		return errors.New("invalid label value")
	}
	h.argparser.addBoolArg(aliases, doc)
	return nil
}

// Add an integer parameter to the subcommand with a default value.
func (h *SubcommandHandler) AddIntParamWithDefault(aliases []string, doc string, deflt int) error {
	if !h.checkAliasesAllowed(aliases) {
		return errors.New("invalid label value")
	}
	h.paramparser.addIntArg(aliases, doc)
	h.paramparser.setIntArg(aliases, deflt)
	return nil
}

// Add a float parameter to the subcommand with a default value.
func (h *SubcommandHandler) AddFloatParamWithDefault(aliases []string, doc string, deflt float64) error {
	if !h.checkAliasesAllowed(aliases) {
		return errors.New("invalid label value")
	}
	h.paramparser.addFloatArg(aliases, doc)
	h.paramparser.setFloatArg(aliases, deflt)
	return nil
}

// Add string parameter to the subcommand with a default value.
func (h *SubcommandHandler) AddStrParamWithDefault(aliases []string, doc string, deflt string) error {
	if !h.checkAliasesAllowed(aliases) {
		return errors.New("invalid label value")
	}
	h.paramparser.addStrArg(aliases, doc)
	h.paramparser.setStrArg(aliases, deflt)
	return nil
}

// Add a boolean parameter to the subcommand with a default value.
func (h *SubcommandHandler) AddBoolParamWithDefault(aliases []string, doc string, deflt bool) error {
	if !h.checkAliasesAllowed(aliases) {
		return errors.New("invalid label value")
	}
	h.paramparser.addBoolArg(aliases, doc)
	h.paramparser.setBoolArg(aliases, deflt)
	return nil
}

// Add an example to a SubcommandHandler instance.
// The argument `doc` is documentation for the example.
// The argument `cmd` is the text of the command with command line arguments.
// The argument `out` is the plausible output for the command
func (h *SubcommandHandler) Example(doc string, cmd string, out string) {
	h.examples = append(h.examples, subcommandExample{
		documentation: doc, command: cmd, output: out})
}

// Set the handler function for a SubcommandHandler instance.
// # Arguments
// - sub: a SubcommandHandler instance to be mutated
// - f: the handler function
func (h *SubcommandHandler) Handle(f func(insub *SubcommandHandler)) {
	h.handle = f
}

// Get an integer argument or parameter from the command line.
// Warning: This function will fail if the command line arguments have
// not already been parsed.
func (h *SubcommandHandler) GetInt(key string) (int, error) {
	if val, b := h.argparser.intValues[key]; b {
		return val, nil
	}
	if val, b := h.paramparser.intValues[key]; b {
		return val, nil
	}
	return 0, errors.New("key not available")
}

// Get a string argument or parameter from the command line.
// Warning: This function will fail if the command line arguments have
// not already been parsed.
func (h *SubcommandHandler) GetStr(key string) (string, error) {
	if val, b := h.argparser.strValues[key]; b {
		return val, nil
	}
	if val, b := h.paramparser.strValues[key]; b {
		return val, nil
	}
	return "", errors.New("key not available")
}

// Get a string argument or parameter from the command line.
// Warning: This function will fail if the command line arguments have
// not already been parsed.
func (h *SubcommandHandler) GetFloat(key string) (float64, error) {
	if val, b := h.argparser.floatValues[key]; b {
		return val, nil
	}
	if val, b := h.paramparser.floatValues[key]; b {
		return val, nil
	}
	return 0.0, errors.New("key not available")
}

// Get a string argument or parameter from the command line.
// Warning: This function will fail if the command line arguments have
// not already been parsed.
func (h *SubcommandHandler) GetBool(key string) (val bool, err error) {
	if val, b := h.argparser.boolValues[key]; b {
		return val, nil
	}
	if val, b := h.paramparser.boolValues[key]; b {
		return val, nil
	}
	return false, errors.New("key not available")
}

// Parse the command line flags.
func (h *SubcommandHandler) parseFlags(args []string) {
	h.argparser.parseFlags(args)
	h.paramparser.parseFlags(args)
}

// Print argument documentation.
func (h *SubcommandHandler) printArgumentHelp() {
	s := h.argparser.helpString()
	if len(s) > 0 {
		fmt.Printf("ARGUMENTS\n%s\n\n", s)
	}
}

// Print option documentation.
func (h *SubcommandHandler) printOptionHelp() {
	s := h.paramparser.helpString()
	if len(s) > 0 {
		fmt.Printf("OPTIONS\n%s\n\n", s)
	}
}

// Print example documentation.
func (h *SubcommandHandler) printExampleHelp() {
	if len(h.examples) > 0 {
		fmt.Printf("EXAMPLES\n")
		for _, ex := range h.examples {
			fmt.Print(ex.helpMessage())
		}
	}
}
