package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	godaemon "github.com/sevlyar/go-daemon"

	"github.com/hoppxi/bpv/internal/daemon"
	"github.com/hoppxi/bpv/internal/logger"
	"github.com/hoppxi/bpv/internal/xdg"
	"github.com/spf13/cobra"
)

const version = "0.1.0"

var (
	verbose     bool
	noDaemonize bool
)

var rootCmd = &cobra.Command{
	Use:     "bpvd [flags]",
	Version: version,
	Short:   "BPV daemon — background music library service",
	Long: `bpvd is the BPV daemon that manages your music library.

  bpvd                   Start daemon (daemonize by default)
  bpvd --no-daemonize    Run in foreground (used for services)`,
	Run: func(cmd *cobra.Command, args []string) {
		if noDaemonize {
			logger.Init(verbose, true)
			runDaemon()
			return
		}

		logPath := xdg.LogPath()
		os.MkdirAll(xdg.StateDir(), 0755)

		pidPath := xdg.PidPath()
		os.MkdirAll(xdg.RuntimeDir(), 0755)

		ctx := &godaemon.Context{
			PidFileName: pidPath,
			PidFilePerm: 0644,
			LogFileName: logPath,
			LogFilePerm: 0640,
			WorkDir:     "/",
			Umask:       027,
		}

		child, err := ctx.Reborn()
		if err != nil {
			log.Fatalf("Failed to daemonize: %v", err)
		}

		if child != nil {
			fmt.Printf("bpvd started (pid %d)\n", child.Pid)
			fmt.Printf("  Log: %s\n", logPath)
			fmt.Printf("  PID: %s\n", pidPath)
			return
		}

		defer ctx.Release()

		logger.Init(verbose, false)
		runDaemon()
	},
}

func runDaemon() {
	d, err := daemon.NewDaemon()
	if err != nil {
		logger.Log.FatalErr(err, "Failed to create daemon")
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Log.Info("Shutting down daemon...")
		d.Stop()
		os.Exit(0)
	}()

	if err := d.Start(); err != nil {
		logger.Log.FatalErr(err, "Daemon stopped unexpectedly")
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose log")
	rootCmd.Flags().BoolVar(&noDaemonize, "no-daemonize", false, "run in foreground without forking")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
