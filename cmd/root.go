package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "app",
	Short: "A music player application",
	Long: `This is a command-line music player application. It supports playing music files in mp3, flac, or wav format.
You can use the 'play' command followed by the file path to play a music file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Put your code here
	},
}

var Interactive bool

func init() {
	RootCmd.PersistentFlags().BoolVarP(&Interactive, "interactive", "i", false, "Start the interactive mode for the app")
	RootCmd.AddCommand(exitCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

var exitCmd = &cobra.Command{
	Use:     "exit",
	Aliases: []string{"q", "Q", "bye"},
	Short:   "Exit the application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Goodbye!")
		os.Exit(0)
	},
}
