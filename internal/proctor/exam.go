package proctor

import(
	"fmt"

	. "exnal/tui/questions"
	"exnal/internal/anal"
	"math/rand"
)

var questions = []Question{
	// NewChoiceSingle("1", []string{"1", "2", "3", "4"}),
	// NewInputSince("2"),
}
type itemWithFlag struct {
	opt string
	flag bool
}

// 全部选项带真假
var dataOpts [][]itemWithFlag

func init() {
	sections, err := anal.ParseMarkdown("test.md")
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}

	dataOpts = make([][]itemWithFlag, len(sections))
	for i, sec := range sections {
		for _, v := range sec.Items.Trues {
			dataOpts[i] = append(dataOpts[i],
				itemWithFlag{v, true,},
			)
		}
		for _, v := range sec.Items.Falses {
			dataOpts[i] = append(dataOpts[i],
				itemWithFlag{v, false,},
			)
		}
	}

	if len(sections) == len(dataOpts) {
		for k := range len(dataOpts) {
			rand.Shuffle(len(dataOpts[k]), func(i, j int) {
				dataOpts[k][i], dataOpts[k][j] = dataOpts[k][j], dataOpts[k][i]
			})
		}

		for i , sec := range sections {
			var opts []string
			for _, item := range dataOpts[i] {
				opts = append(opts, item.opt)
			}
			q := NewChoiceSingle(sec.Title, opts)
			questions = append(questions, q)
		}

		rand.Shuffle(len(dataOpts), func(i, j int) {
			dataOpts[i], dataOpts[j] = dataOpts[j], dataOpts[i]
			questions[i], questions[j] = questions[j], questions[i]
		})
	} else {
		panic("漏答案")
	}
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
