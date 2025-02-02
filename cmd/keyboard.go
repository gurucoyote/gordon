package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	// "time"

	"github.com/eiannone/keyboard"
	// "github.com/gopxl/beep/speaker"
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
			// Delegate play/pause toggling to the existing pause subcommand.
			pauseCmd.Run(pauseCmd, []string{})
		case key == keyboard.KeyArrowLeft:
			// Rewind 5 seconds via existing subcommand.
			rewindCmd.Run(rewindCmd, []string{"5"})
		case key == keyboard.KeyArrowRight:
			// Forward 5 seconds via existing subcommand.
			forwardCmd.Run(forwardCmd, []string{"5"})
		case key == keyboard.KeyArrowUp:
			// Increase volume via existing subcommand.
			newPercent := int((ap.volume.Volume + 0.1) * 100)
			if newPercent > 100 {
				newPercent = 100
			}
			volumeCmd.Run(volumeCmd, []string{strconv.Itoa(newPercent)})
		case key == keyboard.KeyArrowDown:
			// Decrease volume via existing subcommand.
			newPercent := int((ap.volume.Volume - 0.1) * 100)
			if newPercent < 0 {
				newPercent = 0
			}
			volumeCmd.Run(volumeCmd, []string{strconv.Itoa(newPercent)})
		default:
			// Ignore unknown keys.
		}
		// time.Sleep(100 * time.Millisecond)
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
