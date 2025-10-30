package cmd

import (
	"fmt"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/spf13/cobra"
	"os"
	"sync"
	"time"
)

const defaultSampleRate = beep.SampleRate(44100)

var (
	speakerOnce    sync.Once
	speakerInitErr error
)

func ensureSpeaker() error {
	speakerOnce.Do(func() {
		sr := defaultSampleRate
		speakerInitErr = speaker.Init(sr, sr.N(time.Second/10))
	})
	if speakerInitErr != nil {
		return fmt.Errorf("failed to init speaker: %w", speakerInitErr)
	}
	return nil
}

var RootCmd = &cobra.Command{
	Use:   "app",
	Short: "A music player application",
	Long: `This is a command-line music player application. It supports playing music files in mp3, flac, or wav format.
You can use the 'play' command followed by the file path to play a music file.`,
	Args: cobra.ArbitraryArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return ensureSpeaker()
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
