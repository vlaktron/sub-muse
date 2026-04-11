package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"sub-muse/internal/config"
	"sub-muse/internal/setup"
)

func main() {
	// Check if we're in development mode
	if len(os.Args) > 1 && os.Args[1] == "dev" {
		fmt.Println("Development mode - UI not fully implemented yet")
		os.Exit(0)
	}

	// Check for reconfigure flag
	reconfigure := false
	for _, arg := range os.Args {
		if arg == "--reconfigure" {
			reconfigure = true
			break
		}
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		if err == config.ErrNotConfigured || reconfigure {
			runSetupWizard()
			return
		}
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Sub-muse - TUI Music Streaming Client\n")
	fmt.Printf("Configured for: %s\n", cfg.ServerURL)

	// For now, just show a simple message
	fmt.Println("This application is a work in progress.")
	fmt.Println("UI components are being implemented.")
}

func runSetupWizard() {
	p := tea.NewProgram(setup.NewModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running setup wizard: %v", err)
	}

	// Reload config after setup
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config after setup: %v", err)
	}

	fmt.Printf("Sub-muse - TUI Music Streaming Client\n")
	fmt.Printf("Configured for: %s\n", cfg.ServerURL)
	fmt.Println("Setup complete! Starting application...")
}
