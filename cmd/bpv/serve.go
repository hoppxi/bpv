package bpv

import (
	"os"
	"path/filepath"

	"github.com/hoppxi/bpv/internal/server"
	"github.com/hoppxi/bpv/pkg/logger"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve [music directory]",
	Short: "Start the BPV server",
	Long: `Start the BPV server to serve music files from the specified directory.
The server will extract metadata and provide a web interface for playback.

Example:
  bpv serve /path/to/music/folder
  bpv serve --port 3000 /path/to/music/folder`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		musicDir := args[0]
		
		if _, err := os.Stat(musicDir); os.IsNotExist(err) {
			logger.Log.Error("Directory does not exist: %s", musicDir)
		}
		
		absPath, err := filepath.Abs(musicDir)
		if err != nil {
			logger.Log.Error("Error getting absolute path: %v", err)
		}
		
		port, _ := rootCmd.PersistentFlags().GetInt("port")
		
		logger.Log.Info("Starting BPV server on port %d", port)
		logger.Log.Info("Serving music from: %s", absPath)
		
		server := server.NewServer(port, absPath)
		if err := server.Start(); err != nil {
			logger.Log.Error("Failed to start server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}