package anal

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)
type Items struct{
	Trues  []string
	Falses []string
	Tip	   []string
}

type Section struct {
	Number int
	Title  string
	Items
}

// 保护反引号 `...` 内的内容
// 移除常规注释 <!-- ... -->
// 保留其他所有文本
func cleanItemText(raw string) string {
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

	inGlobalComment := false
	scanner := bufio.NewScanner(file)

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

		isTitle := titleRegex.MatchString(trimmed)
		isItem := itemRegex.MatchString(trimmed)

		// ---------- 2. 截断逻辑 ----------
		if currentSection != nil && !isTitle && !isItem {
			currentSection = nil
			continue
		}

		// ---------- 3. 匹配标题 ----------
		if matches := titleRegex.FindStringSubmatch(trimmed); matches != nil {
			num, _ := strconv.Atoi(matches[1])
			newSection := Section{
				Number: num,
				Title:  matches[2],
			}
			sections = append(sections, newSection)
			currentSection = &sections[len(sections)-1]
			continue
		}

		// ---------- 4. 匹配列表项（使用清洗函数） ----------
		matches := itemRegex.FindStringSubmatch(trimmed)

		if matches != nil && currentSection != nil{
			i, v := matches[1], cleanItemText(matches[2])
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
