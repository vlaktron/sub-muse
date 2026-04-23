package theme

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Colors struct {
	Accent              string `toml:"accent"`
	Foreground          string `toml:"foreground"`
	Background          string `toml:"background"`
	SelectionForeground string `toml:"selection_foreground"`
	SelectionBackground string `toml:"selection_background"`
	Color0              string `toml:"color0"`
	Color1              string `toml:"color1"`
	Color2              string `toml:"color2"`
	Color3              string `toml:"color3"`
	Color4              string `toml:"color4"`
	Color5              string `toml:"color5"`
	Color6              string `toml:"color6"`
	Color7              string `toml:"color7"`
	Color8              string `toml:"color8"`
	Color9              string `toml:"color9"`
	Color10             string `toml:"color10"`
	Color11             string `toml:"color11"`
	Color12             string `toml:"color12"`
	Color13             string `toml:"color13"`
	Color14             string `toml:"color14"`
	Color15             string `toml:"color15"`
}

var Default = Colors{
	Accent:              "#82FB9C",
	Foreground:          "#ddf7ff",
	Background:          "#0B0C16",
	SelectionForeground: "#0B0C16",
	SelectionBackground: "#ddf7ff",
	Color0:              "#0B0C16",
	Color1:              "#50f872",
	Color2:              "#4fe88f",
	Color3:              "#50f7d4",
	Color4:              "#829dd4",
	Color5:              "#86a7df",
	Color6:              "#7cf8f7",
	Color7:              "#85E1FB",
	Color8:              "#6a6e95",
	Color9:              "#85ff9d",
	Color10:             "#9cf7c2",
	Color11:             "#a4ffec",
	Color12:             "#c4d2ed",
	Color13:             "#cddbf4",
	Color14:             "#d1fffe",
	Color15:             "#ddf7ff",
}

func LoadOrDefault() Colors {
	home, err := os.UserHomeDir()
	if err != nil {
		return Default
	}

	themePath := filepath.Join(home, ".config", "omarchy", "current", "theme", "colors.toml")
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		return Default
	}

	var colors Colors
	if _, err := toml.DecodeFile(themePath, &colors); err != nil {
		fmt.Printf("Warning: Failed to parse theme file: %v\n", err)
		return Default
	}

	return colors
}
