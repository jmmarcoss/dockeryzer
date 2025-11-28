package cmd

import (
	"fmt"
	"os"

	"github.com/jorgevvs2/dockeryzer/src/functions"
	"github.com/spf13/cobra"
)

var analyzeDockerfile bool

var analyzeCmd = &cobra.Command{
	Use:   "analyze [image|Dockerfile]",
	Short: "Analyze Docker image or Dockerfile based on CIS Docker Benchmark",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide an image name or Dockerfile path")
			os.Exit(0)
		}

		target := args[0]

		if analyzeDockerfile {
			functions.AnalyzeDockerfile(target)
		} else {
			functions.AnalyzeImage(target)
		}
	},
}

func init() {
	analyzeCmd.Flags().BoolVarP(&analyzeDockerfile, "dockerfile", "d", false, "Analyze a Dockerfile instead of an image")
	rootCmd.AddCommand(analyzeCmd)
}
