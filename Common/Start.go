package Common

import (
	"fmt"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

func Start() {
	// 1. 生成艺术字
	myFigure := figure.NewFigure("codeql_n1ght", "", true)
	asciiArt := myFigure.String()

	// 2. 拆分行
	lines := strings.Split(asciiArt, "\n")

	// 3. 计算最长行长度
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// 4. 定义颜色（标准蓝色）
	blue := color.New(color.FgBlue).SprintFunc()

	// 5. 打印上边框
	fmt.Println(strings.Repeat("-", maxLen+4))

	// 6. 打印正文带边框和颜色
	for _, line := range lines {
		fmt.Printf("| %s%s |\n", blue(line), strings.Repeat(" ", maxLen-len(line)))
	}

	// 7. 打印下边框
	fmt.Println(strings.Repeat("-", maxLen+4))
}
