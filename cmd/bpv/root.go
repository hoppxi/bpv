package bpv

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bpv",
	Version: "0.1.0",
	Short: "BPV is a browser-based music player",
	Long: `BPV (Browser Player for Video/Audio) is a music player that runs in your browser.
It serves music files from a local directory and provides a web interface for playback.

Usage:
  bpv serve /path/to/music/folder
  bpv serve --port 3000 /path/to/music/folder`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntP("port", "p", 8080, "port to run the server on")
}
