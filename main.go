package main

import (
	tea "charm.land/bubbletea/v2"
	"exnal/internal/core"
	"exnal/internal/proctor"
	"flag"
	"fmt"
	"os"
)

func main() {
	filePath := flag.String("f", "", "Markdown 文件路径")

	// 解析命令行参数
	flag.Parse()

	// 检查是否提供了 -f
	if *filePath == "" {
		fmt.Println("错误: 请使用 -f 指定文件路径")
		flag.Usage()
		os.Exit(1)
	}

	proctor.PreTopic(*filePath)
	p := tea.NewProgram(core.InitialModel())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
