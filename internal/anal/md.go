package anal

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Section struct {
	Number int
	Title  string
	Items  []string
	// Right  int
}

// cleanItemText 清洗列表项文本：
//  1. 保护反引号 `...` 内的内容（不解析其中的注释）
//  2. 移除常规注释 <!-- ... -->
//  3. 保留其他所有文本
func cleanItemText(raw string) string {
	// 第一步：提取所有反引号内容，替换为占位符
	codeRegex := regexp.MustCompile("`[^`]*`")
	codeBlocks := codeRegex.FindAllString(raw, -1)
	placeholder := "__CODE_BLOCK_%d__"

	for i, block := range codeBlocks {
		raw = strings.Replace(raw, block, fmt.Sprintf(placeholder, i), 1)
	}

	// 第二步：移除注释（此时反引号内容已被保护）
	commentRegex := regexp.MustCompile(`<!--.*?-->`)
	cleaned := commentRegex.ReplaceAllString(raw, "")

	// 第三步：还原反引号内容
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
	itemRegex := regexp.MustCompile(`^-\s+(.*)`)

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
				Items:  []string{},
			}
			sections = append(sections, newSection)
			currentSection = &sections[len(sections)-1]
			continue
		}

		// ---------- 4. 匹配列表项（使用清洗函数） ----------
		if matches := itemRegex.FindStringSubmatch(trimmed); matches != nil {
			if currentSection != nil {
				cleaned := cleanItemText(matches[1])
				currentSection.Items = append(currentSection.Items, cleaned)
			}
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return sections, nil
}
