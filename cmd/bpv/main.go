package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/hoppxi/bpv/internal/daemon"
	"github.com/hoppxi/bpv/internal/logger"
	"github.com/hoppxi/bpv/internal/server"
	"github.com/hoppxi/bpv/internal/tui"
	"github.com/spf13/cobra"
)

const version = "0.1.0"

var (
	verbose bool
	port    int
	client  string
)

var rootCmd = &cobra.Command{
	Use:     "bpv [flags] [music-directory]",
	Version: version,
	Short:   "BPV — music player for your terminal and browser",
	Long: `BPV is a music player that serves your local music library with
a web interface or a TUI.

  bpv ~/Music              Start TUI player (default)
  bpv ~/Music --client web Start web server & open browser
  bpv                      Re-open last used directory`,
	Args: cobra.MaximumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Init(verbose, true)
	},
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		var musicDir string

		if len(args) >= 1 {
			musicDir = args[0]
		} else {
			musicDir = lastUsedDir()
			if musicDir == "" {
				logger.Log.Fatal("No music directory provided and no previous directory found.\n  Usage: bpv <music-directory>")
			}
			logger.Log.Info("Using previous directory: %s", musicDir)
		}

		if len(musicDir) >= 2 && musicDir[:2] == "~/" {
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

		switch client {
		case "web", "b", "browser":
			startServer(port, absPath, start)
			return
		case "tui", "t", "terminal":
			startTUI(absPath)
			return
		default:
			logger.Log.Fatal("client can only be either web or tui")
		}
	},
}

func lastUsedDir() string {
	c, err := daemon.Connect()
	if err != nil {
		return ""
	}
	defer c.Close()

	settings, err := c.GetSettings()
	if err != nil || settings == nil {
		return ""
	}
	if settings.LastDir == "" {
		return ""
	}

	if _, err := os.Stat(settings.LastDir); os.IsNotExist(err) {
		return ""
	}
	return settings.LastDir
}

func startServer(port int, absPath string, start time.Time) {
	addr := fmt.Sprintf("http://localhost:%d/", port)

	arrow := color.New(color.FgGreen, color.Bold).Sprint("➜")
	header := color.New(color.FgHiCyan, color.Bold).SprintFunc()
	label := color.New(color.FgWhite, color.Bold).SprintFunc()
	url := color.New(color.FgHiYellow).SprintFunc()

	fmt.Printf("\n%s %s\n\n", header("Using Port -"), color.HiYellowString("%d", port))
	fmt.Printf("  %s v%s  %s\n\n", header("BPV"), version, color.HiGreenString("started in %d ms", time.Since(start).Milliseconds()))
	fmt.Printf("  %s  %s: \t%s\n", arrow, label("Local"), url(addr))
	fmt.Printf("  %s  %s: \tServing music from %s\n", arrow, label("Files"), url(absPath))
	fmt.Printf("  %s  %s: \tPress Ctrl+C to stop the server\n\n", arrow, label("More"))

	if err := exec.Command("xdg-open", addr).Start(); err != nil {
		logger.Log.Error("Failed to open browser: %v", err)
	} else {
		logger.Log.Info("Opening BPV in the browser")
	}

	srv := server.NewServer(port, absPath)
	if err := srv.Start(); err != nil {
		logger.Log.FatalErr(err, "Failed to start server")
	}
	logger.Log.LogDuration(start, "Server startup")
}

func startTUI(absPath string) {
	model := tui.New(absPath)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose log")
	rootCmd.Flags().StringVarP(&client, "client", "c", "tui", "which client to use (web or tui based)")
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "port to run the server on")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
