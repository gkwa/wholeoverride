package cmd

import (
	"github.com/spf13/cobra"

	"github.com/gkwa/wholeoverride/core"
)

var baseDir string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a markdown table from recipe files",
	Long:  `Generate a markdown table containing recipe and creator information from a directory of markdown files.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		logger.Info("Running generate command")
		if err := core.GenerateTable(logger, baseDir); err != nil {
			logger.Error(err, "Failed to generate table")
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVar(&baseDir, "basedir", "", "Base directory containing markdown files")
	if err := generateCmd.MarkFlagRequired("basedir"); err != nil {
		panic(err)
	}
}
