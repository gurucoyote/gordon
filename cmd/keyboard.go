package cmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/eiannone/keyboard"
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
	JumpSec := "1"
	fmt.Println("Keyboard Control Mode (Normal Mode):")
	fmt.Println("  Space       : Toggle play/pause")
	fmt.Printf("  Left Arrow  : Rewind %s seconds\n", JumpSec)
	fmt.Printf("  Right Arrow : Forward %s seconds\n", JumpSec)
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
		if char == 'l' {
			c1, _, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			if c1 < '0' || c1 > '9' {
				continue // ignore if not a valid digit
			}
			c2, _, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			if c2 < '0' || c2 > '9' {
				continue // ignore if not a valid digit
			}
			RootCmd.SetArgs([]string{"loop", string(c1), string(c2)})
			if err := RootCmd.Execute(); err != nil {
				fmt.Println(err)
			}
			continue
		}
		if char == 'g' {
			digit, _, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			if digit < '0' || digit > '9' {
				continue // ignore if not a valid digit
			}
			RootCmd.SetArgs([]string{"goto", string(digit)})
			if err := RootCmd.Execute(); err != nil {
				fmt.Println(err)
			}
			continue
		}
		if char >= '1' && char <= '9' {
			RootCmd.SetArgs([]string{"setmarker", string(char)})
			if err := RootCmd.Execute(); err != nil {
				fmt.Println(err)
			}
			continue
		}
		if char == 'q' || char == 'Q' {
			RootCmd.SetArgs([]string{"exit"})
			if err := RootCmd.Execute(); err != nil {
				fmt.Println(err)
			}
			break
		}
		if char == ':' {
			keyboard.Close()
			commandMode()
			if err := keyboard.Open(); err != nil {
				panic(err)
			}
			// fmt.Println("Resuming Normal Mode...")
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
			// Rewind .5 seconds via existing subcommand.
			RootCmd.SetArgs([]string{"rewind", JumpSec})
			if err := RootCmd.Execute(); err != nil {
				fmt.Println(err)
			}
		case key == keyboard.KeyArrowRight: // Forward .5 seconds via existing subcommand.
			RootCmd.SetArgs([]string{"forward", JumpSec})
			if err := RootCmd.Execute(); err != nil {
				fmt.Println(err)
			}
		case key == keyboard.KeyArrowUp:
			// Increase volume via existing subcommand.
			currentPercent := int(math.Pow(ap.volume.Base, ap.volume.Volume) * 100)
			newPercent := currentPercent + 10
			if newPercent > 100 {
				newPercent = 100
			}
			RootCmd.SetArgs([]string{"volume", strconv.Itoa(newPercent)})
			if err := RootCmd.Execute(); err != nil {
				fmt.Println(err)
			}
		case key == keyboard.KeyArrowDown:
			// Decrease volume via existing subcommand.
			currentPercent := int(math.Pow(ap.volume.Base, ap.volume.Volume) * 100)
			newPercent := currentPercent - 10
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
		ControlLoop()
		return
	},
}

func init() {
	RootCmd.AddCommand(keyboardCmd)
}
