package cmd

import (
	"fmt"
	"os"

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
		// output what second we are on now
		// fmt.Print(format.SampleRate.D(streamer.Position()).Round(time.Second))
		speaker.Unlock()
		ap.play()
		// speaker.Play(ctrl)
		return
	},
}

var rewindCmd = &cobra.Command{
	Use:     "rewind [seconds]",
	Aliases: []string{"rw"},
	Short:   "Rewind playback position by n seconds",
	Long:    `Rewind playback position by n seconds.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rewind command stub")
		ap.play()
		/*
			if event.Rune() == 'w' {
				newPos += ap.sampleRate.N(time.Second)
			}
						if newPos < 0 {
				newPos = 0
			}
			if newPos >= ap.streamer.Len() {
				newPos = ap.streamer.Len() - 1
			}
			if err := ap.streamer.Seek(newPos); err != nil {
				report(err)
			}

		*/
	},
}

var forwardCmd = &cobra.Command{
	Use:     "forward [seconds]",
	Aliases: []string{"fw"},
	Short:   "Forward playback position by n seconds",
	Long:    `Forward playback position by n seconds.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Forward command stub")
		ap.play()
	},
}

var stopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{"s"},
	Short:   "Stop playback",
	Long:    `Stop playback.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stop command stub")
		return
	},
}

func init() {
	RootCmd.AddCommand(playCmd)
	RootCmd.AddCommand(pauseCmd, rewindCmd, forwardCmd, stopCmd)
}
