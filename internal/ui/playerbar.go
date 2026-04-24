package ui

import (
	"fmt"
	"strings"
)

type PlayerBar struct {
	Progress    float64
	CurrentTime int
	Duration    int
	Playing     bool
	SongTitle   string
	Artist      string
	Album       string
	Styles      Styles
}

func (p *PlayerBar) FormatTime(seconds int) string {
	minutes := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}

func (p *PlayerBar) Render(width int) string {
	if width < 20 {
		return p.renderMinimal(width)
	}

	lines := []string{}

	playIcon := "□"
	if p.Playing {
		playIcon = "■"
	}
	metadata := fmt.Sprintf("%s %s - %s", playIcon, p.SongTitle, p.Artist)
	metadata = truncateString(metadata, width)
	lines = append(lines, p.Styles.PlayerBar.Render(metadata))

	progressFilled := int(p.Progress * float64(width-8))
	progressEmpty := width - 8 - progressFilled

	progressBar := "["
	for i := 0; i < progressFilled; i++ {
		progressBar += "█"
	}
	for i := 0; i < progressEmpty; i++ {
		progressBar += "░"
	}
	progressBar += "]"
	lines = append(lines, p.Styles.PlayerBar.Render(progressBar))

	timeStr := fmt.Sprintf("%s / %s", p.FormatTime(p.CurrentTime), p.FormatTime(p.Duration))
	lines = append(lines, p.Styles.PlayerBar.Render(timeStr))

	return strings.Join(lines, "\n")
}

func (p *PlayerBar) renderMinimal(width int) string {
	lines := []string{}

	playIcon := "□"
	if p.Playing {
		playIcon = "■"
	}
	metadata := fmt.Sprintf("%s %s", playIcon, truncateString(p.SongTitle, width-4))
	lines = append(lines, p.Styles.PlayerBar.Render(metadata))

	timeStr := fmt.Sprintf("%s / %s", p.FormatTime(p.CurrentTime), p.FormatTime(p.Duration))
	lines = append(lines, p.Styles.PlayerBar.Render(timeStr))

	return strings.Join(lines, "\n")
}

func truncateString(s string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}
	if maxWidth < 3 {
		return s[:maxWidth]
	}
	return s[:maxWidth-3] + "..."
}
