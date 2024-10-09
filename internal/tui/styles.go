package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func HasColors() bool {
	return lipgloss.ColorProfile() != termenv.Ascii
}
