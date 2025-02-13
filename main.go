package main

import (
	"gordon/cmd"
)

func main() {
	// run any command implied by the command line (play a file)
	cmd.Execute()
	// toggle playback (workaround for now)
	cmd.RootCmd.SetArgs([]string{"pause"})
	cmd.RootCmd.Execute()
	// drop into keyboard mode
	cmd.RootCmd.SetArgs([]string{"keyboard"})
	cmd.RootCmd.Execute()
}
