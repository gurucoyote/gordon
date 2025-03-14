package cmd

import (
	"fmt"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var RootCmd = &cobra.Command{
	Use:   "app",
	Short: "A music player application",
	Long: `This is a command-line music player application. It supports playing music files in mp3, flac, or wav format.
You can use the 'play' command followed by the file path to play a music file.`,
	Args: cobra.ArbitraryArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		sr := beep.SampleRate(44100)
		speaker.Init(sr, sr.N(time.Second/10))
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if ap != nil {
			ap.play()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 {
			loadCmd.Run(cmd, args)
		} else {
			cmd.Help()
		}
	},
}

func init() {
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
