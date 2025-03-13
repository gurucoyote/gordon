package cmd

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
	"github.com/spf13/cobra"
)

var pinkNoiseActive bool

var pinknoiseCmd = &cobra.Command{
	Use:   "pinknoise [volume]",
	Short: "Toggle pink noise playback; optional volume between 0 and 100%",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if pinkNoiseActive {
			speaker.Lock()
			speaker.Clear()
			speaker.Unlock()
			pinkNoiseActive = false
			fmt.Println("Pink noise stopped")
		} else {
			volumePercent := 100.0
			if len(args) > 0 {
				v, err := strconv.ParseFloat(args[0], 64)
				if err != nil {
					fmt.Printf("Failed to parse volume: %v\n", err)
					return
				}
				if v < 0 || v > 100 {
					fmt.Println("Volume must be between 0 and 100")
					return
				}
				volumePercent = v
			}
			pn := NewPinkNoise()
			vol := &effects.Volume{
				Streamer: pn,
				Base:     2,
				Volume:   math.Log(volumePercent/100) / math.Log(2),
			}
			speaker.Lock()
			speaker.Play(vol)
			speaker.Unlock()
			pinkNoiseActive = true
			fmt.Printf("Pink noise started at volume: %.0f%%\n", volumePercent)
		}
	},
}

func init() {
	RootCmd.AddCommand(pinknoiseCmd)
}
