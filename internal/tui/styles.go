package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
)

func HasColors() bool {
	return lipgloss.ColorProfile() != termenv.Ascii
}

func IsATTY(fd uintptr) bool {
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}
