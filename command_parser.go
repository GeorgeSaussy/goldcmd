package goldcmd

import (
	"errors"
	"fmt"
	"strconv"
)

// A struct to parse arguments
type commandParser struct {
	// all argument labels this subcommand uses
	allLabels []string

	// Here, the command line options are called "labels".
	// Labels are associated with one of the supported types.
	// Each label can have a number of aliases.
	//
	// Example:
	// $ some_command --a=some_string_argument -b 1
	// In this example, the labels are "a" (of string type) and "b" (of int type).
	// Note, as far as the library is concerned, this command is equivalent to
	// $ some_command -a "some_string_argument" --b=1
	//
	// All of the aliases for a label are stored together in an array.
	// The set of sets of aliases are then stored in an array, one for each supported type.
	// This way, strong typing is enforced without need for reflection.

	intLabels   [][]string
	strLabels   [][]string
	floatLabels [][]string
	boolLabels  [][]string

	// After parsing the values, all of the arguments are stored in a map from label aliases to
	// the value set by the user.
	// If a label has multiple aliases, all of them will be mapped to the value set by the user.

	intValues   map[string]int
	strValues   map[string]string
	floatValues map[string]float64
	boolValues  map[string]bool

	// Documentation for the arguments. Keys are the documentation value, values
	// are lists of labels associated with the argument.
	menu map[string][]string
}

// Create a new, but empty commandParser instance. Application code will set the
// fill its state.
func newCommandParser() *commandParser {
	return &commandParser{
		allLabels:   make([]string, 0),
		intLabels:   make([][]string, 0),
		strLabels:   make([][]string, 0),
		floatLabels: make([][]string, 0),
		boolLabels:  make([][]string, 0),
		intValues:   make(map[string]int),
		strValues:   make(map[string]string),
		floatValues: make(map[string]float64),
		boolValues:  make(map[string]bool),
		menu:        make(map[string][]string),
	}
}

// Return true if a label alias is already in use.
func (cp *commandParser) hasAlias(newAlias string) bool {
	for _, label := range cp.allLabels {
		if newAlias == label {
			return true
		}
	}
	return false
}

// Add a set of label aliases to a commandParser instance.
// This should only be called by the add*Arg functions.
func (cp *commandParser) addLabel(aliases []string, doc string, typeLabels *[][]string) {
	cp.allLabels = append(cp.allLabels, aliases...)
	*typeLabels = append(*typeLabels, aliases)
	cp.menu[doc] = aliases
}

// Add a label set for a new integer argument.
// Warning, the caller should check the labels are not in use before calling this function.
func (cp *commandParser) addIntArg(aliases []string, doc string) {
	cp.addLabel(aliases, doc, &cp.intLabels)
}

// Add a label set for a new string argument.
// Warning, the caller should check the labels are not in use before calling this function.
func (cp *commandParser) addStrArg(aliases []string, doc string) {
	cp.addLabel(aliases, doc, &cp.strLabels)
}

// Add a label set for a new float argument.
// Warning, the caller should check the labels are not in use before calling this function.
func (cp *commandParser) addFloatArg(aliases []string, doc string) {
	cp.addLabel(aliases, doc, &cp.floatLabels)
}

// Add a label set for a new Boolean argument.
// Warning, the caller should check the labels are not in use before calling this function.
func (cp *commandParser) addBoolArg(aliases []string, doc string) {
	cp.addLabel(aliases, doc, &cp.boolLabels)
}

// Set the value for an integer argument.
func (cp *commandParser) setIntArg(aliases []string, value int) {
	for _, alias := range aliases {
		cp.intValues[alias] = value
	}
}

// Set the value for a string argument.
func (cp *commandParser) setStrArg(aliases []string, value string) {
	for _, alias := range aliases {
		cp.strValues[alias] = value
	}
}

// Set the value for a float argument.
func (cp *commandParser) setFloatArg(aliases []string, value float64) {
	for _, alias := range aliases {
		cp.floatValues[alias] = value
	}
}

// Set the value for a Boolean argument.
func (cp *commandParser) setBoolArg(aliases []string, value bool) {
	for _, alias := range aliases {
		cp.boolValues[alias] = value
	}
}

// Parse the command line flags. The variable "args" is the command line arguments.
func (cp *commandParser) parseFlags(args []string) {
	// The first and second argument are the command invocation and the
	// subcommand, so we can skip them. If the binary for the CLI is 'cli' and
	// the subcommand is 'sub' with only one argument set by label 'label' and
	// value 'val', then the following invocations are valid:
	// - $ cli sub --label val
	// - $ cli sub -label val
	// - $ cli sub --label=val
	// - $ cli sub -label=val
	// If the argument type is Boolean, then argument can be implicitly set:
	// - $ cli sub --label # implicit true
	// If multiple labels set a value the last one is used. Fight me.
	k := 2
	for k < len(args) {
		arg := args[k]
		start := 0
		for arg[start] == '-' && start < len(arg) {
			start += 1
		}
		if start == 1 || start == 2 {
			end := 0
			for end < len(arg) && arg[end] != '=' {
				end += 1
			}
			label := arg[start:end]
			if end < len(arg)-1 {
				cp.tryToUseFlag(label, arg[end+1:])
				k++
			} else if k < len(args)-1 {
				if cp.tryToUseFlag(label, args[k+1]) == nil {
					k += 2
				} else {
					k++
				}
			} else {
				cp.tryToUseFlag(label, "")
				k++
			}
		} else {
			k++
		}
	}
}

// Check if a string is in a (2D) list of strings. The function returns
// the row in which the string can be found or -1 if it is not in the string.
func strInStrList(s string, l [][]string) int {
	for n, row := range l {
		for _, elm := range row {
			if s == elm {
				return n
			}
		}
	}
	return -1
}

// Try to use a possible flag and value. The function will return an error if the
// label / value combination cannot be used by the command specification.
func (cp *commandParser) tryToUseFlag(alias string, possibleValue string) error {
	if row := strInStrList(alias, cp.intLabels); row >= 0 {
		if val, err := strconv.Atoi(possibleValue); err == nil {
			cp.setIntArg(cp.intLabels[row], val)
			return nil
		} else {
			return err
		}
	}

	if row := strInStrList(alias, cp.strLabels); row >= 0 {
		cp.setStrArg(cp.strLabels[row], possibleValue)
		return nil
	}
	if row := strInStrList(alias, cp.floatLabels); row >= 0 {
		if val, err := strconv.ParseFloat(possibleValue, 64); err == nil {
			cp.setFloatArg(cp.floatLabels[row], val)
			return nil
		} else {
			return err
		}
	}

	if row := strInStrList(alias, cp.boolLabels); row >= 0 {
		if possibleValue == "true" {
			cp.setBoolArg(cp.boolLabels[row], true)
			return nil
		} else if possibleValue == "false" {
			cp.setBoolArg(cp.boolLabels[row], false)
			return nil
		} else {
			cp.setBoolArg(cp.boolLabels[row], true)
			return nil
		}
	}

	return errors.New(fmt.Sprintf("The label \"%s\" was not found to be a supported type.", alias))
}

// Get the help string for a commandParser instance.
func (cp *commandParser) helpString() string {
	var s string = ""
	for doc, labels := range cp.menu {
		var ls string = ""
		for k, label := range labels {
			if k > 0 {
				ls = ls + ","
			}
			ls = ls + " --" + label
		}
		s = s + ls + "\t" + doc + "\n"
	}
	return s
}
