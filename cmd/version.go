package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gkwa/wholeoverride/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of wholeoverride",
	Long:  `All software has versions. This is wholeoverride's`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		logger.Info("Getting version information")

		buildInfo := version.GetBuildInfo()
		fmt.Println(buildInfo)

		logger.Info("Version command completed")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
