package main

import (
	"github.com/chzyer/readline"
	"gordon/cmd"
	"strings"
)

func main() {
	cmd.Execute()
	// start in keyboard mode
	cmd.RootCmd.SetArgs([]string{"keyboard"})
	cmd.RootCmd.Execute()
	if cmd.Interactive {
		// enter repl loop
		rl, _ := readline.New("> ")
		defer rl.Close()
		for {
			input, _ := rl.Readline()
			args := strings.Fields(input)
			cmd.RootCmd.SetArgs(args)
			cmd.RootCmd.Execute()
		}
	}
}
