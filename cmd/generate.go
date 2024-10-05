package cmd

import (
	"github.com/spf13/cobra"

	"github.com/gkwa/wholeoverride/core"
)

var (
	baseDir string
	format  string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a markdown file from recipe files",
	Long:  `Generate a markdown file containing recipe and creator information from a directory of markdown files.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		logger.Info("Running generate command")

		if err := core.GenerateMarkdownWithFormat(logger, baseDir, format); err != nil {
			logger.Error(err, "Failed to generate markdown")
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().
		StringVar(&baseDir, "basedir", "", "Base directory containing markdown files")
	generateCmd.Flags().
		StringVar(&format, "format", "sections", "Output format (sections or table)")
	if err := generateCmd.MarkFlagRequired("basedir"); err != nil {
		panic(err)
	}
}
