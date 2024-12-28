// Package main is the entry point for the dashboard application
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Parse command line flags
	debug := flag.Bool("debug", false, "Run with debug output")
	external := flag.Bool("external", false, "Run in external window")
	flag.Parse()

	if *external {
		// Start a new process without the -external flag
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", "start", os.Args[0])
		} else {
			cmd = exec.Command(os.Args[0])
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		if err := cmd.Start(); err != nil {
			fmt.Printf("Error starting external process: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Initialize the dashboard
	dashboard, err := InitializeDashboard()
	if err != nil {
		fmt.Printf("Error initializing dashboard: %v\n", err)
		os.Exit(1)
	}

	if *debug {
		dashboard.EnableDebug()
	}

	// Create and run the program
	p := tea.NewProgram(
		dashboard,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
