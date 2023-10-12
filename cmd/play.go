package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"

	"github.com/spf13/cobra"
)

// var ctrl *beep.Ctrl
var ap *audioPanel

type audioPanel struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
}

func newAudioPanel(sampleRate beep.SampleRate, streamer beep.StreamSeeker) *audioPanel {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}
	return &audioPanel{sampleRate, streamer, ctrl, resampler, volume}
}

func (ap *audioPanel) play() {
	speaker.Play(ap.volume)
}

var playCmd = &cobra.Command{
	Use:   "play [file]",
	Short: "Play a music file",
	Long:  `Play a music file. The file must be in either mp3, flac, or wav format.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// make sure we get to enter commands after playback starts
		// also, this will 'block' so that the sound can play before the program ends
		Interactive = true
		// load the file
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
		// what info do we get here?
		fmt.Println(format)
		// defer streamer.Close()

		// ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
		ap = newAudioPanel(format.SampleRate, streamer)
		ap.play()
		// speaker.Play(ctrl)
		// this should drop us into interactive mode and continue playing
		return
	},
}

var pauseCmd = &cobra.Command{
	Use:     "pause",
	Aliases: []string{"p"},
	Short:   "Toggle play/pause of current sound",
	Long:    `Toggle play/pause of current sound.`,
	Run: func(cmd *cobra.Command, args []string) {
		// pause/resume playback
		speaker.Lock()
		ap.ctrl.Paused = !ap.ctrl.Paused
		position := ap.sampleRate.D(ap.streamer.Position())
		length := ap.sampleRate.D(ap.streamer.Len())
		volume := ap.volume.Volume
		speaker.Unlock()
		ap.play()
		positionStatus := fmt.Sprintf("%v / %v", position.Round(time.Second), length.Round(time.Second))
		volumeStatus := fmt.Sprintf("%.1f", volume)
		fmt.Println(positionStatus, volumeStatus)
		// speaker.Play(ctrl)
		return
	},
}

var rewindCmd = &cobra.Command{
	Use:     "rewind [seconds]",
	Aliases: []string{"rw"},
	Short:   "Rewind playback position by n seconds",
	Run: func(cmd *cobra.Command, args []string) {
		var relpos float64 = 1.0
		if len(args) > 0 {
			var err error
			relpos, err = strconv.ParseFloat(args[0], 64)
			if err != nil {
				fmt.Printf("Failed to parse argument: %s\n", err)
				return
			}
		}
		// negate it so we go backward
		relpos = relpos * -1
		fmt.Printf("rewind command with relative position: %f\n", relpos)
		seekPos(relpos)
		ap.play()
	},
}

var forwardCmd = &cobra.Command{
	Use:     "forward [seconds]",
	Aliases: []string{"fw"},
	Short:   "Forward playback position by n seconds",
	Run: func(cmd *cobra.Command, args []string) {
		var relpos float64 = 1.0
		if len(args) > 0 {
			var err error
			relpos, err = strconv.ParseFloat(args[0], 64)
			if err != nil {
				fmt.Printf("Failed to parse argument: %s\n", err)
				return
			}
		}
		fmt.Printf("Forward command with relative position: %f\n", relpos)
		seekPos(relpos)
		ap.play()
	},
}

var volumeCmd = &cobra.Command{
	Use:     "volume",
	Aliases: []string{"vol"},
	Short:   "set volume in 0-100%",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vol, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Failed to parse argument: %s\n", err)
			return
		}
		if vol < 0 || vol > 100 {
			fmt.Println("Volume must be between 0 and 100")
			return
		}
		speaker.Lock()
		ap.volume.Volume = float64(vol) / 100
		speaker.Unlock()
		ap.play()
		fmt.Printf("Volume set to %d%%\n", vol)
	},
}

func init() {
	RootCmd.AddCommand(playCmd)
	RootCmd.AddCommand(pauseCmd, rewindCmd, forwardCmd, volumeCmd)
}
func seekPos(pos float64) {
	newPos := ap.streamer.Position()
	// move this by the passed float seconds
	newPos += ap.sampleRate.N(time.Duration(pos) * time.Second)
	if newPos < 0 {
		newPos = 0
	}
	if newPos >= ap.streamer.Len() {
		newPos = ap.streamer.Len() - 1
	}
	if err := ap.streamer.Seek(newPos); err != nil {
		fmt.Println(err)
	}

}
