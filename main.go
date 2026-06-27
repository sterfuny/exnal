package main

import (
	tea "charm.land/bubbletea/v2"
	. "exnal/tui/questions"
)

type model struct {
    questions []Question
    current   int
    done      bool
    quit      bool
}

func initialModel() model {
    questions := []Question{
        NewChoiceSingal("你最喜欢的编程语言？", []string{"Go", "Python", "Rust", "Java"}),
        NewInputSince("你的名字是什么？"),
        NewChoiceSingal("你选择哪个难度？", []string{"简单", "普通", "困难"}),
        NewInputSince("你的目标是什么？"),
    }
    
    return model{
        questions: questions,
        current:   0,
        done:      false,
        quit:      false,
    }
}

func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        panic(err)
    }
}
