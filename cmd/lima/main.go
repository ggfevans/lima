package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ggfevans/lima/internal/app"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("lima %s (%s, %s)\n", version, commit, date)
			return
		}
	}

	m := app.New()

	p := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Send program reference for real-time event bridging
	go func() {
		p.Send(app.ProgramRefMsg{Program: p})
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
