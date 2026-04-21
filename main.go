package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"sub-muse/internal/config"
	"sub-muse/internal/setup"
	"sub-muse/internal/subsonic"

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

	fmt.Println("DEBUG: After check - starting player")

	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Sub-muse - TUI Music Streaming Client\n")
	fmt.Printf("Configured for: %s\n", cfg.ServerURL)
	fmt.Printf("Password: %s\n", cfg.Password)
	fmt.Println("Starting music playback...")

	// Start playing music directly without TUI
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		client := subsonic.NewClient(cfg.ServerURL, cfg.Username, cfg.Password, cfg.ClientName)
		songs, err := client.GetSongs()
		if err != nil {
			log.Fatalf("Failed to get songs: %v", err)
		}

		if len(songs) == 0 {
			log.Fatalf("No songs found")
		}

		song := songs[0]
		audioData, err := client.Stream(subsonic.WithID(song.ID))
		if err != nil {
			log.Fatalf("Failed to stream: %v", err)
		}

		// Play audio
		ext := ".wav"
		if len(audioData) > 4 && string(audioData[:4]) != "RIFF" {
			ext = ".ogg"
		}

		tmpFile, err := os.CreateTemp("", "sub-muse-*"+ext)
		if err != nil {
			log.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write(audioData)
		if err != nil {
			log.Fatalf("Failed to write temp file: %v", err)
		}
		tmpFile.Close()

		cmd := exec.CommandContext(ctx, "ffplay", "-nodisp", "-autoexit", tmpFile.Name())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Failed to play audio: %v", err)
		}
	}()

	// Wait for Ctrl+C
	select {}
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
	fmt.Println("Setup complete! Starting application...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		client := subsonic.NewClient(cfg.ServerURL, cfg.Username, cfg.Password, cfg.ClientName)
		songs, err := client.GetSongs()
		if err != nil {
			log.Fatalf("Failed to get songs: %v", err)
		}

		if len(songs) == 0 {
			log.Fatalf("No songs found")
		}

		song := songs[0]
		audioData, err := client.Stream(subsonic.WithID(song.ID))
		if err != nil {
			log.Fatalf("Failed to stream: %v", err)
		}

		ext := ".wav"
		if len(audioData) > 4 && string(audioData[:4]) != "RIFF" {
			ext = ".ogg"
		}

		tmpFile, err := os.CreateTemp("", "sub-muse-*"+ext)
		if err != nil {
			log.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write(audioData)
		if err != nil {
			log.Fatalf("Failed to write temp file: %v", err)
		}
		tmpFile.Close()

		cmd := exec.CommandContext(ctx, "ffplay", "-nodisp", "-autoexit", tmpFile.Name())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Failed to play audio: %v", err)
		}
	}()

	select {}
}
