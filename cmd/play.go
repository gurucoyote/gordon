package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
	"github.com/gopxl/beep/wav"
	"strings"

	"github.com/spf13/cobra"
)

var (
	Markers []PlaybackPosition
	format  beep.Format
)

type audioPanel struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	loop       *loopBetween
	volume     *effects.Volume
}

func newAudioPanel(sampleRate beep.SampleRate, streamer beep.StreamSeeker) *audioPanel {
	// ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	loop := LoopBetween(-1, 0, streamer.Len(), streamer)
	ctrl := &beep.Ctrl{Streamer: loop}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: ctrl, Base: 2}
	return &audioPanel{sampleRate,
		streamer,
		ctrl,
		resampler,
		loop,
		volume}
}

func (ap *audioPanel) play() {
	speaker.Play(ap.volume)
}

var ap *audioPanel

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

		// Determine the file type and decode accordingly
		var streamer beep.StreamSeeker
		var decodedFormat beep.Format

		switch {
		case strings.HasSuffix(file, ".mp3"):
			streamer, decodedFormat, err = mp3.Decode(f)
		case strings.HasSuffix(file, ".wav"):
			streamer, decodedFormat, err = wav.Decode(f)
		case strings.HasSuffix(file, ".flac"):
			streamer, decodedFormat, err = flac.Decode(f)
		case strings.HasSuffix(file, ".ogg"):
			streamer, decodedFormat, err = vorbis.Decode(f)
		default:
			fmt.Printf("Unsupported file format: %s\n", file)
			return
		}
		format = decodedFormat
		if err != nil {
			fmt.Printf("Failed to decode file: %s\n", err)
			return
		}
		// what info do we get here?
		fmt.Println(format)
		// defer streamer.Close()

		ap = newAudioPanel(format.SampleRate, streamer)
		ap.play()
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

var setMarkerCmd = &cobra.Command{
	Use:     "setmarker [marker]",
	Aliases: []string{"m"},
	Short:   "Set a marker",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Failed to parse argument: %s\n", err)
			return
		}

		speaker.Lock()
		samplePosition := ap.streamer.Position()
		playPosition := ap.sampleRate.D(samplePosition).Seconds()
		speaker.Unlock()

		newMarker := PlaybackPosition{
			SamplePosition: samplePosition,
			PlayPosition:   playPosition,
		}

		for len(Markers) <= index {
			Markers = append(Markers, PlaybackPosition{})
		}
		Markers[index] = newMarker

		fmt.Printf("Marker %d set to sample position %d (play position %.2f seconds)\n", index, samplePosition, playPosition)
	},
}

var gotoCmd = &cobra.Command{
	Use:   "goto [marker]",
	Short: "Go to a marker",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		markerIndex, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Failed to parse argument: %s\n", err)
			return
		}
		if markerIndex < 0 || markerIndex >= len(Markers) {
			fmt.Printf("Marker %d does not exist\n", markerIndex)
			return
		}
		marker := Markers[markerIndex]
		if err := ap.streamer.Seek(marker.SamplePosition); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Go to marker %d\n", markerIndex)
	},
}

var loopCmd = &cobra.Command{
	Use:   "loop [start_marker] [end_marker]",
	Short: "Loop between two markers",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		startMarker, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Failed to parse start_marker argument: %s\n", err)
			return
		}
		endMarker, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Failed to parse end_marker argument: %s\n", err)
			return
		}
		fmt.Printf("Loop between markers %d and %d\n", startMarker, endMarker)
		startPos := Markers[startMarker].SamplePosition
		endPos := Markers[endMarker].SamplePosition
		speaker.Lock()
		ap.loop.start = startPos
		ap.loop.end = endPos
		speaker.Unlock()
		ap.play()

	},
}

var saveCmd = &cobra.Command{
	Use:   "save [start_marker] [end_marker] [output_file]",
	Short: "Save the loop between two markers to a file",
	Long:  `Save the audio loop between two specified markers to a .wav file.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		startMarker, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Failed to parse start_marker argument: %s\n", err)
			return
		}
		endMarker, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Failed to parse end_marker argument: %s\n", err)
			return
		}
		outputFile := args[2]

		if startMarker < 0 || startMarker >= len(Markers) || endMarker < 0 || endMarker >= len(Markers) {
			fmt.Println("Invalid marker indices")
			return
		}

		startPos := Markers[startMarker].SamplePosition
		endPos := Markers[endMarker].SamplePosition

		if startPos >= endPos {
			fmt.Println("Start marker must be before end marker")
			return
		}

		// Create the output file
		f, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("Failed to create output file: %s\n", err)
			return
		}
		defer f.Close()

		// Seek to the start position
		if err := ap.streamer.Seek(startPos); err != nil {
			fmt.Printf("Failed to seek to start position: %s\n", err)
			return
		}

		// Create a buffer for the segment
		buffer := beep.NewBuffer(format)
		segment := beep.Take(endPos-startPos, ap.streamer)
		buffer.Append(segment)

		// Create a streamer from the buffer
		streamer := buffer.Streamer(0, buffer.Len())

		// Encode the streamer to a wav file
		if err := wav.Encode(f, streamer, format); err != nil {
			fmt.Printf("Failed to encode wav file: %s\n", err)
			return
		}

		fmt.Printf("Saved segment between markers %d and %d to %s\n", startMarker, endMarker, outputFile)
	},
}

func init() {
	RootCmd.AddCommand(playCmd, pauseCmd, rewindCmd, forwardCmd, volumeCmd, setMarkerCmd, gotoCmd, loopCmd, saveCmd)
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

type PlaybackPosition struct {
	SamplePosition int
	PlayPosition   float64
}

// LoopBetween takes a StreamSeeker and plays it between start and end positions. If count is negative, s is looped infinitely.
//
// The returned Streamer propagates s's errors.
func LoopBetween(count int, start int, end int, s beep.StreamSeeker) *loopBetween {
	return &loopBetween{
		s:       s,
		remains: count,
		start:   start,
		end:     end,
	}
}

type loopBetween struct {
	s       beep.StreamSeeker
	remains int
	start   int
	end     int
}

func (l *loopBetween) Stream(samples [][2]float64) (n int, ok bool) {
	if l.remains == 0 || l.s.Err() != nil {
		return 0, false
	}
	for len(samples) > 0 {
		sn, sok := l.s.Stream(samples)
		if !sok || l.s.Position() >= l.end {
			if l.remains > 0 {
				l.remains--
			}
			if l.remains == 0 {
				break
			}
			err := l.s.Seek(l.start)
			if err != nil {
				return n, true
			}
			continue
		}
		samples = samples[sn:]
		n += sn
	}
	return n, true
}

func (l *loopBetween) Err() error {
	return l.s.Err()
}
