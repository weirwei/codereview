package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of CodeReview",
	Long:  `All software has versions. This is CodeReview's`,
	Run: func(cmd *cobra.Command, args []string) {
		versionFile := "VERSION"
		content, err := os.ReadFile(versionFile)
		if err != nil {
			fmt.Printf("Error reading VERSION file: %v\n", err)
			return
		}
		fmt.Printf("CodeReview %s", string(content))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
