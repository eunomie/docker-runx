package tui

import (
	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
)

func GetTermRenderer() *glamour.TermRenderer {
	w, _, err := term.GetSize(0)
	if err != nil {
		w = 80
	}

	var r *glamour.TermRenderer
	if HasColors() {
		r, _ = glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(w-10))
	} else {
		r, _ = glamour.NewTermRenderer(
			glamour.WithStandardStyle("notty"),
			glamour.WithWordWrap(w-10))
	}
	return r
}

func Markdown(str string) string {
	r := GetTermRenderer()
	out, err := r.Render(str)
	if err != nil {
		return str
	}
	return out
}

func BackQuoteItems(items []string) []string {
	var out []string
	for _, item := range items {
		out = append(out, "`"+item+"`")
	}
	return out
}
