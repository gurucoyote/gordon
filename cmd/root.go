package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
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
	rootCmd.PersistentFlags().BoolVarP(&Interactive, "interactive", "i", false, "Start the interactive mode for the app")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
