package main

import (
	tea "charm.land/bubbletea/v2"
	"exnal/internal/core"
	"exnal/internal/proctor"
	"flag"
	"os"
)

func main() {
	filePath := flag.String("f", "", "Markdown file path")

	flag.Parse()

	// 检查是否提供了 -f
	if *filePath == "" {
		flag.Usage()
		os.Exit(64)
	}

	if err := proctor.PreTopic(*filePath); err != nil {
		os.Exit(66)
	}

	p := tea.NewProgram(core.InitialModel())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
