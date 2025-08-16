package Scanner

import (
	"codeql_n1ght/Common"
	"fmt"

	"github.com/fatih/color"
)

// showPackInstallHint 显示pack install提示
func showPackInstallHint() {
	yellow := color.New(color.FgYellow).SprintFunc()
	Common.LogWarn("如果遇到package相关错误，请尝试以下解决方案：")
	fmt.Printf("%s\n", yellow("1. 进入QL库目录: cd "+Common.QLLibsPath))
	fmt.Printf("%s\n", yellow("2. 运行命令: ../tools/codeql/codeql pack install"))
	fmt.Printf("%s\n", yellow("3. 或者运行: codeql pack install"))
	fmt.Println()
}
