package bpv

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/hoppxi/bpv/internal/logger"
	"github.com/hoppxi/bpv/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve [music directory]",
	Short: "Start the BPV server",
	Long: `Start the BPV server to serve music files from the specified directory.
The server will extract metadata and provide a web interface for playback.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		musicDir := args[0]
		addr := fmt.Sprintf("http://localhost:%d/", port)
		
		if musicDir[:2] == "~/" {
			if home, err := os.UserHomeDir(); err == nil {
				musicDir = filepath.Join(home, musicDir[2:])
			}
		}
		
		if _, err := os.Stat(musicDir); os.IsNotExist(err) {
			logger.Log.FatalP("Directory check", "Directory does not exist: %s", musicDir)
		}
		
		absPath, err := filepath.Abs(musicDir)
		if err != nil {
			logger.Log.FatalP("Path resolution", "Error getting absolute path: %v", err)
		}
		
		arrow := color.New(color.FgGreen, color.Bold).Sprint("➜")
		header := color.New(color.FgHiCyan, color.Bold).SprintFunc()
		label := color.New(color.FgWhite, color.Bold).SprintFunc()
		url := color.New(color.FgHiYellow).SprintFunc()

		fmt.Printf("\n%s %s\n\n", header("Using Port -"), color.HiYellowString("%d", port))
		fmt.Printf("  %s v%s  %s\n\n", header("BPV"), rootCmd.Version, color.HiGreenString("started in %d ms", time.Since(start).Milliseconds()))
		fmt.Printf("  %s  %s: \t%s\n", arrow, label("Local"), url(addr))
		fmt.Printf("  %s  %s: \tServing music from %s\n", arrow, label("Files"), url(absPath))
		fmt.Printf("  %s  %s: \tPress Ctrl+C to stop the server\n\n", arrow, label("More"))
		
		if open {
			if err := exec.Command("xdg-open", addr).Start(); err != nil {
				logger.Log.Error("Failed to open browser: %v", err)
			} else {
				logger.Log.Info("Opening BPV in the browser")
			}
		}
		
		server := server.NewServer(port, absPath)
		if err := server.Start(); err != nil {
			logger.Log.FatalErr(err, "Failed to start server")
		}
		logger.Log.LogDuration(start, "Server startup")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}