package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play [file]",
	Short: "Play a music file",
	Long: `Play a music file. The file must be in either mp3, flac, or wav format.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("File %s does not exist\n", file)
			return
		}
		// TODO: Add code to play the file
	},
}

func init() {
	RootCmd.AddCommand(playCmd)
}
