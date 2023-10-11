package cmd

import (
	"bufio"
	"fmt"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
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
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		// ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
		ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
		speaker.Play(ctrl)

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Press [ENTER] to pause/resume, 'q' to stop. ")
		for {
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "q" || input == "Q" {
				speaker.Close()
				f.Close()
				return
			}
			speaker.Lock()
			// pause/resume playback
			ctrl.Paused = !ctrl.Paused
			// output what second we are on now
			fmt.Print("\r                                                                 \r")
			fmt.Print(format.SampleRate.D(streamer.Position()).Round(time.Second))
			speaker.Unlock()

		}
	},
}

func init() {
	RootCmd.AddCommand(playCmd)
}
