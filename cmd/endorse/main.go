package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ggfevans/endorse/internal/app"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

func main() {
	demoMode := false
	themeName := ""
	for _, arg := range os.Args[1:] {
		switch {
		case arg == "-v" || arg == "--version":
			fmt.Printf("endorse %s (%s, %s)\n", version, commit, date)
			return
		case arg == "--demo":
			demoMode = true
		case strings.HasPrefix(arg, "--theme="):
			themeName = strings.TrimPrefix(arg, "--theme=")
		}
	}

	m := app.New(app.Options{DemoMode: demoMode, ThemeName: themeName})

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
