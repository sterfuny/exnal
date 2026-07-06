package anal

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"exnal/tui/questions"
)

type Items struct{
	Trues  []string
	Falses []string
	Tip	   []string
}

type Section struct {
	Number int
	Title  string
	Type   questions.QuestionType
	Items
}

// 保护反引号 `...` 内的内容
// 移除常规注释 <!-- ... -->
// 保留其他所有文本
func cleanItemMix(raw string) string {
	// 反引号内容替换为占位符
	codeRegex := regexp.MustCompile("`[^`]*`")
	codeBlocks := codeRegex.FindAllString(raw, -1)
	placeholder := "__CODE_BLOCK_%d__"

	for i, block := range codeBlocks {
		raw = strings.Replace(raw, block, fmt.Sprintf(placeholder, i), 1)
	}

	// 移除注释
	commentRegex := regexp.MustCompile(`<!--.*?-->`)
	cleaned := commentRegex.ReplaceAllString(raw, "")

	// 还原反引号内容
	for i, block := range codeBlocks {
		cleaned = strings.Replace(cleaned, fmt.Sprintf(placeholder, i), block, -1)
	}

	return strings.TrimSpace(cleaned)
}

func ParseMarkdown(filePath string) ([]Section, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sections []Section
	var currentSection *Section

	titleRegex := regexp.MustCompile(`^(\d+)\.\s+(.*)`)
	itemRegex := regexp.MustCompile(`^([-+*])\s+(.*)`)
	typeRegex := regexp.MustCompile(`^#{1,6}\s+(.*)`)

	inGlobalComment := false
	scanner := bufio.NewScanner(file)
	var nowType questions.QuestionType

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// ---------- 1. 处理跨行全局注释 ----------
		if inGlobalComment {
			if strings.Contains(line, "-->") {
				inGlobalComment = false
			}
			continue
		}
		if strings.Contains(line, "<!--") && !strings.HasPrefix(trimmed, "- ") {
			inGlobalComment = true
			if strings.Contains(line, "-->") {
				inGlobalComment = false
			}
			continue
		}

		if trimmed == "" {
			continue
		}

		// ---------- 2. 截断逻辑 ----------
		isFlag := titleRegex.MatchString(trimmed) || itemRegex.MatchString(trimmed) || typeRegex.MatchString(trimmed)
		if currentSection != nil && !isFlag {
			currentSection = nil
			continue
		}

		if matches := typeRegex.FindStringSubmatch(trimmed); matches != nil {
			if t := questions.ParseType(matches[1]);t != -1 {
				nowType = t
			}
			continue
		}

		// ---------- 3. 匹配标题 ----------
		if matches := titleRegex.FindStringSubmatch(trimmed); matches != nil {
			num, _ := strconv.Atoi(matches[1])
			newSection := Section{
				Number: num,
				Title:  matches[2],
				Type: 	nowType,
			}
			sections = append(sections, newSection)
			currentSection = &sections[len(sections)-1]
			continue
		}

		// ---------- 4. 匹配列表项 ----------
		if matches := itemRegex.FindStringSubmatch(trimmed); matches != nil && currentSection != nil{
			i, v := matches[1], cleanItemMix(matches[2])
			item := &currentSection.Items

			if i == "-" {
				item.Falses = append(item.Falses, v)
			} else if i == "+" {
				item.Trues = append(item.Trues, v)
			}
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return sections, nil
}
