package cmd

import (
	"fmt"
	// "time"

	"github.com/gopxl/beep/speaker"
	"github.com/spf13/cobra"
)

var posCmd = &cobra.Command{
	Use:   "pos",
	Short: "Show playback position",
	Long:  "Show current playback position and total length in seconds with at least 3 digits of precision",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if ap == nil {
			fmt.Println("No audio loaded!")
			return
		}
		speaker.Lock()
		position := ap.sampleRate.D(ap.streamer.Position()).Seconds()
		length := ap.sampleRate.D(ap.streamer.Len()).Seconds()
		volume := ap.volume.Volume
		speaker.Unlock()
		fmt.Printf("%.3f / %.3f (Volume: %.1f)\n", position, length, volume)
	},
}

func init() {
	RootCmd.AddCommand(posCmd)
}
