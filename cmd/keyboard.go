package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/chzyer/readline"
	// "github.com/gopxl/beep/speaker"
	"github.com/spf13/cobra"
)

// commandMode prompts the user for a full command input similar to vim's command mode.
func commandMode() {
	// enter repl loop
	rl, _ := readline.New("> ")
	defer rl.Close()
	input, _ := rl.Readline()
	args := strings.Fields(input)
	RootCmd.SetArgs(args)
	RootCmd.Execute()
}

// ControlLoop starts in normal mode where keys control media playback.
// In normal mode:
//   - Space toggles play/pause
//   - Left Arrow rewinds 5 seconds
//   - Right Arrow forwards 5 seconds
//   - Up Arrow increases volume
//   - Down Arrow decreases volume
//   - ':' enters command mode
//   - Q exits keyboard control mode
func ControlLoop() {
	fmt.Println("Keyboard Control Mode (Normal Mode):")
	fmt.Println("  Space       : Toggle play/pause")
	fmt.Println("  Left Arrow  : Rewind 5 seconds")
	fmt.Println("  Right Arrow : Forward 5 seconds")
	fmt.Println("  Up Arrow    : Increase volume")
	fmt.Println("  Down Arrow  : Decrease volume")
	fmt.Println("  :           : Enter command mode")
	fmt.Println("  Q           : Quit control mode")

	// Using GetKey()
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() { _ = keyboard.Close() }()

	for {
		char, key, err := keyboard.GetKey()
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
			RootCmd.SetArgs([]string{"pause"})
			if err := RootCmd.Execute(); err != nil {
				fmt.Println(err)
			}
		case key == keyboard.KeyArrowLeft:
			// Rewind 5 seconds via existing subcommand.
			RootCmd.SetArgs([]string{"rewind", "5"})
			if err := RootCmd.Execute(); err != nil {
				fmt.Println(err)
			}
		case key == keyboard.KeyArrowRight:
			// Forward 5 seconds via existing subcommand.
			RootCmd.SetArgs([]string{"forward", "5"})
			if err := RootCmd.Execute(); err != nil {
				fmt.Println(err)
			}
		case key == keyboard.KeyArrowUp:
			// Increase volume via existing subcommand.
			newPercent := int((ap.volume.Volume + 0.1) * 100)
			if newPercent > 100 {
				newPercent = 100
			}
			RootCmd.SetArgs([]string{"volume", strconv.Itoa(newPercent)})
			if err := RootCmd.Execute(); err != nil {
				fmt.Println(err)
			}
		case key == keyboard.KeyArrowDown:
			// Decrease volume via existing subcommand.
			newPercent := int((ap.volume.Volume - 0.1) * 100)
			if newPercent < 0 {
				newPercent = 0
			}
			RootCmd.SetArgs([]string{"volume", strconv.Itoa(newPercent)})
			if err := RootCmd.Execute(); err != nil {
				fmt.Println(err)
			}
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
		// ensure to not interrupt audio playback
		done := make(chan bool)
		go func() {
			ControlLoop()
			done <- true
		}()
		<-done
		return
	},
}

func init() {
	RootCmd.AddCommand(keyboardCmd)
}
