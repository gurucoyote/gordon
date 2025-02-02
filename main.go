package main

import (
	"gordon/cmd"
)

func main() {
	cmd.Execute()
	// start in keyboard mode
	cmd.RootCmd.SetArgs([]string{"keyboard"})
	cmd.RootCmd.Execute()
}
