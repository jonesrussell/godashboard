package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/dashboard/internal/ui"
)

func openExternalWindow() error {
	var cmd *exec.Cmd
	exePath := os.Args[0]

	// If running with 'go run', build a proper executable
	if strings.Contains(exePath, "go-build") {
		fmt.Println("Please use 'task run-external' instead of 'go run'")
		os.Exit(1)
	}

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", exePath)
	case "darwin":
		cmd = exec.Command("open", "-a", "Terminal", exePath)
	default: // Linux and others
		terminals := []string{"gnome-terminal", "xterm", "konsole", "terminator"}
		for _, term := range terminals {
			if path, err := exec.LookPath(term); err == nil {
				cmd = exec.Command(path, "--", exePath)
				break
			}
		}
		if cmd == nil {
			return fmt.Errorf("no suitable terminal emulator found")
		}
	}

	cmd.Env = append(os.Environ(), "LAUNCH_IN_TERMINAL=1")
	return cmd.Start()
}

func main() {
	var external bool
	flag.BoolVar(&external, "external", false, "Launch in external window")
	flag.Parse()

	if external {
		if err := openExternalWindow(); err != nil {
			fmt.Printf("Error launching external window: %v\n", err)
			os.Exit(1)
		}
		return
	}

	opts := []tea.ProgramOption{
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	}

	p := tea.NewProgram(ui.NewDashboard(), opts...)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
