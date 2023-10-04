package cmd

import (
	"fmt"
	"os"
	"time"
	// "github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play [file]",
	Short: "Play a music file",
	Long:  `Play a music file. The file must be in either mp3, flac, or wav format.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("File %s does not exist\n", file)
			return
		}

		// Open the file
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Failed to open file: %s\n", err)
			return
		}

		// Decode the file
		streamer, format, err := mp3.Decode(f)
		if err != nil {
			fmt.Printf("Failed to decode file: %s\n", err)
			return
		}
		defer streamer.Close()

		// Initialize the speaker
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))

		// Play the music
		speaker.Play(streamer)
	},
}

func init() {
	RootCmd.AddCommand(playCmd)
}
