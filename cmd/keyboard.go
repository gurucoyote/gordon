package cmd

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gopxl/beep/speaker"
	"github.com/spf13/cobra"
)

// commandMode prompts the user for a full command input similar to vim's command mode.
func commandMode() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Command Mode - Enter command: ")
	command, _ := reader.ReadString('\n')
	fmt.Printf("Command entered: %s\r\n", command)
}

// ControlLoop starts in normal mode where keys control media playback.
// In normal mode:
//   • Space toggles play/pause
//   • Left Arrow rewinds 5 seconds
//   • Right Arrow forwards 5 seconds
//   • Up Arrow increases volume
//   • Down Arrow decreases volume
//   • ':' enters command mode
//   • Q exits keyboard control mode
func ControlLoop() {
	fmt.Println("Keyboard Control Mode (Normal Mode):")
	fmt.Println("  Space       : Toggle play/pause")
	fmt.Println("  Left Arrow  : Rewind 5 seconds")
	fmt.Println("  Right Arrow : Forward 5 seconds")
	fmt.Println("  Up Arrow    : Increase volume")
	fmt.Println("  Down Arrow  : Decrease volume")
	fmt.Println("  :           : Enter command mode")
	fmt.Println("  Q           : Quit control mode")

	for {
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}
		if char == 'q' || char == 'Q' {
			fmt.Println("Exiting keyboard control mode.")
			break
		}
		if char == ':' {
			commandMode()
			fmt.Println("Resuming Normal Mode...")
			continue
		}
		// Ensure an audio file is loaded (global variable ap from play.go)
		if ap == nil {
			fmt.Println("No audio loaded!")
			continue
		}
		switch {
		case key == keyboard.KeySpace:
			speaker.Lock()
			ap.ctrl.Paused = !ap.ctrl.Paused
			position := ap.sampleRate.D(ap.streamer.Position())
			length := ap.sampleRate.D(ap.streamer.Len())
			volume := ap.volume.Volume
			speaker.Unlock()
			ap.play()
			fmt.Printf("Toggled play/pause: %v / %v (Volume: %.1f)\n", position.Round(time.Second), length.Round(time.Second), volume)
		case key == keyboard.KeyArrowLeft:
			// Rewind 5 seconds.
			seekPos(-5.0)
			ap.play()
			fmt.Println("Rewinded 5 seconds.")
		case key == keyboard.KeyArrowRight:
			// Forward 5 seconds.
			seekPos(5.0)
			ap.play()
			fmt.Println("Forwarded 5 seconds.")
		case key == keyboard.KeyArrowUp:
			speaker.Lock()
			newVol := ap.volume.Volume + 0.1
			if newVol > 1.0 {
				newVol = 1.0
			}
			ap.volume.Volume = newVol
			speaker.Unlock()
			ap.play()
			fmt.Println("Increased volume.")
		case key == keyboard.KeyArrowDown:
			speaker.Lock()
			newVol := ap.volume.Volume - 0.1
			if newVol < 0.0 {
				newVol = 0.0
			}
			ap.volume.Volume = newVol
			speaker.Unlock()
			ap.play()
			fmt.Println("Decreased volume.")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

var keyboardCmd = &cobra.Command{
	Use:   "keyboard",
	Short: "Enter keyboard control mode (vim-like normal mode for media controls)",
	Run: func(cmd *cobra.Command, args []string) {
		if err := keyboard.Open(); err != nil {
			panic(err)
		}
		defer keyboard.Close()
		ControlLoop()
	},
}

func init() {
	RootCmd.AddCommand(keyboardCmd)
}
