package main

import (
	tea "charm.land/bubbletea/v2"
	"exnal/internal/core"
)

func main() {
	p := tea.NewProgram(core.InitialModel())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
