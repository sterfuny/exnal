package proctor

import(
	// "fmt"

	. "exnal/tui/questions"
	// "exnal/internal/core"
)

var questions = []Question{
	NewChoiceSingal("1", []string{"1", "2", "3", "4"}),
	NewInputSince("2"),
}

func GetQuestion(index int) Question {
	return questions[index]
}

func GetQuestionAll() []Question {
	return questions
}
