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
		allitems = append(allitems, sec.Items)
	}
	rand.Shuffle(len(sections), func(i, j int) {
		sections[i], sections[j] = sections[j], sections[i]
		allitems[i], allitems[j] = allitems[j], allitems[i]
	})
	for i := range len(sections) {
		i = i
		// fmt.Printf("%d %s\n%s\n", i, sections[i].Title, allitems[i])
	}
}
