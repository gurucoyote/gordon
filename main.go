package main

import (
	"github.com/chzyer/readline"
	"gordon/cmd"
	"strings"
)

func main() {
	cmd.Execute()
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
