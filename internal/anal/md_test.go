package anal

import "testing"
import (
	"fmt"
	"math/rand"
)

func TestSomething(t *testing.T) {
	sections, err := ParseMarkdown("test.md")
	var allitems [][]string

	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}

	for _, sec := range sections {
		items := append(sec.Trues, sec.Falses...)
		allitems = append(allitems, items)
	}
	rand.Shuffle(len(sections), func(i, j int) {
		sections[i], sections[j] = sections[j], sections[i]
		allitems[i], allitems[j] = allitems[j], allitems[i]
	})
	for i := range len(sections) {
		fmt.Printf("%d %s\n", i, sections[i].Title)
		for _, v := range allitems[i] {
			fmt.Printf("- %s\n", v)
		}
	}
}
