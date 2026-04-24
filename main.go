package main

import (
	"fmt"
	"log"
	"os"

	"sub-muse/internal/config"
	"sub-muse/internal/setup"
	"sub-muse/internal/theme"
	"sub-muse/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
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
	cfgPath, _ := config.GetConfigPath()
	fmt.Printf("DEBUG: config path from GetConfigPath=%s\n", cfgPath)
	configured, err := config.IsConfigured()
	if err != nil {
		log.Fatalf("Failed to check if configured: %v", err)
	}

	cfg, err2 := config.LoadConfig()
	fmt.Printf("DEBUG: configured=%v, reconfigure=%v, err=%v\n", configured, reconfigure, err2)
	if err2 == nil && cfg != nil {
		fmt.Printf("DEBUG: server_url=%s, username=%s\n", cfg.ServerURL, cfg.Username)
	}

	shouldRunWizard := !configured || reconfigure
	fmt.Printf("DEBUG: shouldRunWizard=%v (!configured=%v || reconfigure=%v)\n", shouldRunWizard, !configured, reconfigure)

	if shouldRunWizard {
		fmt.Println("DEBUG: Running setup wizard")
		runSetupWizard()
		return
	}

	fmt.Println("DEBUG: After check - launching TUI")

	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	colors := theme.LoadOrDefault()
	model := ui.NewModel(cfg, colors)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}

func runSetupWizard() {
	wizard := tea.NewProgram(setup.NewModel())
	if _, err := wizard.Run(); err != nil {
		log.Fatalf("Error running setup wizard: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config after setup: %v", err)
	}

	fmt.Printf("Sub-muse - TUI Music Streaming Client\n")
	fmt.Printf("Configured for: %s\n", cfg.ServerURL)
	fmt.Printf("Password: %s\n", cfg.Password)
	fmt.Println("Setup complete! Launching application...")

	colors := theme.LoadOrDefault()
	model := ui.NewModel(cfg, colors)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}
