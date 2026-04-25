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

	playIcon := "□"
	if p.Playing {
		playIcon = "■"
	}
	metadata := truncateString(fmt.Sprintf("%s %s - %s", playIcon, p.SongTitle, p.Artist), width)

	barWidth := max(0, width-8)
	progressFilled := int(p.Progress * float64(barWidth))
	progressEmpty := barWidth - progressFilled
	progressBar := "[" + strings.Repeat("█", progressFilled) + strings.Repeat("░", progressEmpty) + "]"

	timeStr := fmt.Sprintf("%s / %s", p.FormatTime(p.CurrentTime), p.FormatTime(p.Duration))

	content := strings.Join([]string{metadata, progressBar, timeStr}, "\n")
	return p.Styles.PlayerBar.Render(content)
}

func (p *PlayerBar) renderMinimal(width int) string {
	playIcon := "□"
	if p.Playing {
		playIcon = "■"
	}
	titleWidth := max(0, width-4)
	metadata := fmt.Sprintf("%s %s", playIcon, truncateString(p.SongTitle, titleWidth))
	timeStr := fmt.Sprintf("%s / %s", p.FormatTime(p.CurrentTime), p.FormatTime(p.Duration))

	content := strings.Join([]string{metadata, timeStr}, "\n")
	return p.Styles.PlayerBar.Render(content)
}

func truncateString(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if len(s) <= maxWidth {
		return s
	}
	if maxWidth < 3 {
		return s[:maxWidth]
	}
	return s[:maxWidth-3] + "..."
}
