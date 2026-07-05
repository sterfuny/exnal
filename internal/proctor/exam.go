package proctor

import (
	"fmt"

	"exnal/internal/anal"
	. "exnal/tui/questions"
	"math/rand"
)

var questions = []Question{
	// NewChoiceSingle("1", []string{"1", "2", "3", "4"}),
	// NewInputSince("2"),
}

type itemWithFlag struct {
	opt  string
	flag bool
}

// 全部选项带真假
var dataOpts [][]itemWithFlag

func PreTopic(path string) error {
	sections, err := anal.ParseMarkdown(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	dataOpts = make([][]itemWithFlag, len(sections))
	for i, sec := range sections {
		for _, v := range sec.Items.Trues {
			dataOpts[i] = append(dataOpts[i],
				itemWithFlag{v, true},
			)
		}
		for _, v := range sec.Items.Falses {
			dataOpts[i] = append(dataOpts[i],
				itemWithFlag{v, false},
			)
		}
	}

	if len(sections) == len(dataOpts) {
		for k := range len(dataOpts) {
			rand.Shuffle(len(dataOpts[k]), func(i, j int) {
				dataOpts[k][i], dataOpts[k][j] = dataOpts[k][j], dataOpts[k][i]
			})
		}

		for i, sec := range sections {
			var opts []string
			for _, item := range dataOpts[i] {
				opts = append(opts, item.opt)
			}

			if len(opts) == 0 {
				return fmt.Errorf("Some questions do not have an item yet.")
			}

			var q Question
			switch sec.Type {
			case ChoiceSingle:
				q = NewChoiceSingle(sec.Title, opts)
			case ChoiceMulti:
				continue
			case InputSince:
				continue
			case InputChunk:
				continue
			default:
				return fmt.Errorf("unknown questions type")
			}
			questions = append(questions, q)
		}

		rand.Shuffle(len(dataOpts), func(i, j int) {
			dataOpts[i], dataOpts[j] = dataOpts[j], dataOpts[i]
			questions[i], questions[j] = questions[j], questions[i]
		})
	}
	return nil
}

func GetRight(index int) []string {
	var rights []string // (*dataOpts)[index]
	for _, k := range dataOpts[index] {
		if k.flag == true {
			rights = append(rights, k.opt)
		}
	}
	return rights
}

func GetQuestion(index int) Question {
	return questions[index]
}

func GetQuestionAll() []Question {
	return questions
}
