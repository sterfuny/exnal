package proctor

import(
	"fmt"

	. "exnal/tui/questions"
	// "exnal/internal/core"
	"exnal/internal/anal"
	"math/rand"
)

var questions = []Question{
	// NewChoiceSingle("1", []string{"1", "2", "3", "4"}),
	// NewInputSince("2"),
}
var rights [][]string

func init() {
	sections, err := anal.ParseMarkdown("test.md")
	var allitems [][]string = rights

	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}

	for _, sec := range sections {
		allitems = append(allitems, sec.Items)
	}

	if len(sections) == len(allitems) {
		for _ , sec := range sections {
			q := NewChoiceSingle(sec.Title, sec.Items)
			questions = append(questions, q)
		}
		rand.Shuffle(len(questions), func(i, j int) {
			questions[i], questions[j] = questions[j], questions[i]
			// sections[i], sections[j] = sections[j], sections[i]
			// allitems[i], allitems[j] = allitems[j], allitems[i]
		})
	} else {
		panic("漏答案")
	}

	rights = allitems

}

func GetRight(index int) []string {
	return rights[index]
}

func GetQuestion(index int) Question {
	return questions[index]
}

func GetQuestionAll() []Question {
	return questions
}
