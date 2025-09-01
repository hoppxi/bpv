package bpv

import (
	"os"

	"github.com/hoppxi/bpv/internal/logger"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	port    int
	open    bool
)

var rootCmd = &cobra.Command{
	Use:   "bpv",
	Version: "0.1.0",
	Short: "BPV is a browser-based music player",
	Long: `BPV (Browser Player for Video/Audio) is a music player that runs in your browser.
It serves music files from a local directory and provides a web interface for playback.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Init(verbose)
		logger.Log.Info("BPV Music Player starting up...")
		if verbose {
			logger.Log.Debug("Verbose logging enabled")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&open, "open", "o", false, "open BPV with the default browser")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "port to run the server on")
}