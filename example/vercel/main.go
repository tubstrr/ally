package ally_theme

import (
	"embed"
)

var (
	//go:embed theme
	theme embed.FS
)

func GetTheme() embed.FS {
	return theme
}