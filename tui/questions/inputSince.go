package questions

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type InputSinceQ struct {
    question string
    input    string
    cursor   int
    done     bool
}

func NewInputSince(q string) *InputSinceQ {
    return &InputSinceQ{
        question: q,
        input:    "",
        cursor:   0,
        done:     false,
    }
}

func (q *InputSinceQ) Render() string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("❓ %s (填空题)\n\n", q.question))
    
    // 显示输入内容（带光标）
    display := q.input
    if !q.done {
        // 在光标位置显示下划线
        cursorPos := q.cursor
        if cursorPos < len(display) {
            display = display[:cursorPos] + "|" + display[cursorPos:]
        } else {
            display = display + "|"
        }
    }
    sb.WriteString(fmt.Sprintf("答案: %s\n", display))
    
    if q.done {
        sb.WriteString(fmt.Sprintf("\n✅ 已填写: %s", q.input))
    } else {
        sb.WriteString("\n输入字符 • Backspace 删除 • Enter 确认")
    }
    return sb.String()
}

func (q *InputSinceQ) HandleKey(msg tea.KeyMsg) (done bool, quit bool) {
    if q.done {
        return true, false
    }
    
    switch msg.String() {
    case "ctrl+c", "q":
        return false, true
    case "enter":
        // 允许空答案，也可以加校验
        q.done = true
        return true, false
    case "backspace":
        if q.cursor > 0 && len(q.input) > 0 {
            q.input = q.input[:q.cursor-1] + q.input[q.cursor:]
            q.cursor--
        }
    case "left":
        if q.cursor > 0 {
            q.cursor--
        }
    case "right":
        if q.cursor < len(q.input) {
            q.cursor++
        }
    default:
        // 普通字符（只接受可打印字符）
        if len(msg.String()) == 1 {
            // 简单过滤控制字符
            if msg.String() >= " " && msg.String() <= "~" {
                q.input = q.input[:q.cursor] + msg.String() + q.input[q.cursor:]
                q.cursor++
            }
        }
    }
    return false, false
}

func (q *InputSinceQ) GetQuestionText() string {
    return q.question
}

func (q *InputSinceQ) GetAnswer() any {
    return q.input
}

func (q *InputSinceQ) IsDone() bool {
    return q.done
}
