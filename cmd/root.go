package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "app",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Put your code here
	},
}

var Interactive bool

func init() {
	RootCmd.PersistentFlags().BoolVarP(&Interactive, "interactive", "i", false, "Start the interactive mode for the app")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

var exitCmd = &cobra.Command{
	Use:     "exit",
	Aliases: []string{"q", "Q", "bye"},
	Short:   "Exit the application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Goodbye!")
		os.Exit(0)
	},
}
