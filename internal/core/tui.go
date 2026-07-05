package core

import (
	"fmt"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	"exnal/internal/proctor"
	. "exnal/tui/questions"
)

func InitialModel() model {

	var questions []Question = proctor.GetQuestionAll()
	// for i := 0; {
	// 	questions = append(questions, proctor.GetQuestion(i))
	// }

	return model{
		questions: questions,
		current:   0,
		done:      false,
		quit:      false,
	}
}

type model struct {
	questions []Question
	current   int
	done      bool
	quit      bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.done {
			if msg.String() == "ctrl+c" || msg.String() == "q" {
				m.quit = true
				return m, tea.Quit
			}
			return m, nil
		}

		q := m.questions[m.current]
		done, quit := q.HandleKey(msg)

		if quit {
			m.quit = true
			return m, tea.Quit
		}

		if done {
			m.current++
			if m.current >= len(m.questions) {
				m.done = true
			}
		}
	}
	return m, nil
}

func compStrList(a []string, b []string) bool {
	for _, A := range a {
		if slices.Contains(b, A) {
			return true
		}
	}
	return false
}

func (m model) View() tea.View {
	var sb strings.Builder

	// 显示已完成的题目
	for i := 0; i < m.current; i++ {
		q := m.questions[i]
		fmt.Fprintf(&sb, "%s:\n", q.GetQuestionText())
		if answer, ok := q.GetAnswer().([]string); ok {
			if compStrList(answer, proctor.GetRight(i)) {
				fmt.Fprintf(&sb, "√■ ")
			} else {
				fmt.Fprintf(&sb, "x■ ")
			}
			fmt.Fprintf(&sb, "%v\n", answer)
		} else if answer, ok := q.GetAnswer().(string); ok {
			if slices.Contains(proctor.GetRight(i), answer) {
				fmt.Fprintf(&sb, "√■ ")
			} else {
				fmt.Fprintf(&sb, "x■ ")
			}
			fmt.Fprintf(&sb, "%v\n", answer)
		}
	}

	if m.done {
		sb.WriteString("\n🎉 全部完成！\n\n")
		// sb.WriteString("📊 汇总：\n")
		// for _, q := range m.questions {
		//     sb.WriteString(fmt.Sprintf("  %s: %v\n", q.GetQuestionText(), q.GetAnswer()))
		// }
		sb.WriteString("\n按 q 退出")
		return tea.View{Content: sb.String()}
	}

	if m.current > 0 {
		sb.WriteString("\n" + strings.Repeat("─", 40) + "\n\n")
	}

	// 当前题目
	fmt.Fprintf(&sb, "[%d/%d] ", m.current+1, len(m.questions))
	sb.WriteString(m.questions[m.current].Render())

	if m.quit {
		fmt.Fprintf(&sb, "\n已退出\n")
	}

	return tea.View{Content: sb.String()}
}
