package commands

import (
	"fmt"
	"github.com/mdapathy/architecture-4/engine"
	"strings"
)

type printCommand struct {
	arg string
}

type palindromCommand struct {
	pref string
}

func (p printCommand) Execute(loop engine.Handler) {
	fmt.Println(p.arg)
}


func (palindrome *palindromCommand) Execute(loop engine.Handler) {
	res := palindrome.pref
	for i := len(palindrome.pref) - 1; i >= 0; i-- {
		res += string(palindrome.pref[i])
	}

	loop.Post(&printCommand{arg: res})
}

func Parse(str string) engine.Command {
	parts := strings.Fields(str)

	switch {
	case len(parts) < 1:
		return &printCommand{arg: "SYNTAX ERROR: No command specified "}

	case parts[0] == "palindrom" && len(parts) == 2:
		return &palindromCommand{pref: parts[1]}

	case parts[0] == "print" && len(parts) == 2:
		return &printCommand{arg: parts[1]}

	case (parts[0] == "print" || parts[0] == "palindrom") && len(parts) < 2:
		return &printCommand{arg: "SYNTAX ERROR: Not enough arguments for `" + parts[0] + "` command"}

	case (parts[0] == "print" || parts[0] == "palindrom") && len(parts) > 2:
		return &printCommand{arg: "SYNTAX ERROR: Too many arguments in `" + str + "`"}

	default:
		return &printCommand{arg: "SYNTAX ERROR: No such command as `" + parts[0] + "`"}

	}
}
