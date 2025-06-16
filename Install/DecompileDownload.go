package Install

import (
	"codeql_n1ght/Common"
	"fmt"
	"os"
	"path/filepath"
)

// CheckDecompileInstalled 检查反编译器是否已安装在tools目录下
func CheckDecompileInstalled() bool {
	toolsDir := "./tools"
	procyonPath := filepath.Join(toolsDir, "procyon-decompiler-0.6.0.jar")
	fernflowerPath := filepath.Join(toolsDir, "java-decompiler.jar")
	jsp2classPath := filepath.Join(toolsDir, "jsp2class.jar")

	procyonExists := false
	fernflowerExists := false
	jsp2classExists := false

	if _, err := os.Stat(procyonPath); err == nil {
		fmt.Println("procyon-decompiler-0.6.0.jar 已经安装在 ./tools 目录下")
		procyonExists = true
	}

	if _, err := os.Stat(fernflowerPath); err == nil {
		fmt.Println("java-decompiler.jar 已经安装在 ./tools 目录下")
		fernflowerExists = true
	}

	if _, err := os.Stat(jsp2classPath); err == nil {
		fmt.Println("jsp2class.jar 已经安装在 ./tools 目录下")
		jsp2classExists = true
	}

	return procyonExists && fernflowerExists && jsp2classExists
}

// DownloadDecompilers 下载反编译器到tools目录
func DownloadDecompilers() error {
	if CheckDecompileInstalled() {
		return nil
	}

	// 创建tools目录
	toolsDir := "./tools"
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		return fmt.Errorf("创建tools目录失败: %v", err)
	}

	// 下载procyon-decompiler
	procyonPath := filepath.Join(toolsDir, "procyon-decompiler-0.6.0.jar")
	if _, err := os.Stat(procyonPath); os.IsNotExist(err) {
		fmt.Println("开始下载procyon-decompiler-0.6.0.jar...")
		procyonURL := "https://raw.githubusercontent.com/yezere/codeql_n1ght_dp/refs/heads/main/procyon-decompiler-0.6.0.jar"
		if err := Common.DownloadFile(procyonURL, procyonPath); err != nil {
			return fmt.Errorf("下载procyon-decompiler-0.6.0.jar失败: %v", err)
		}
		fmt.Printf("procyon-decompiler-0.6.0.jar下载完成: %s\n", procyonPath)
	}

	// 下载java-decompiler (fernflower)
	fernflowerPath := filepath.Join(toolsDir, "java-decompiler.jar")
	if _, err := os.Stat(fernflowerPath); os.IsNotExist(err) {
		fmt.Println("开始下载java-decompiler.jar...")
		fernflowerURL := "https://raw.githubusercontent.com/yezere/codeql_n1ght_dp/refs/heads/main/java-decompiler.jar"
		if err := Common.DownloadFile(fernflowerURL, fernflowerPath); err != nil {
			return fmt.Errorf("下载java-decompiler.jar失败: %v", err)
		}
		fmt.Printf("java-decompiler.jar下载完成: %s\n", fernflowerPath)
	}

	// 下载jsp2class.jar
	jsp2classPath := filepath.Join(toolsDir, "jsp2class.jar")
	if _, err := os.Stat(jsp2classPath); os.IsNotExist(err) {
		fmt.Println("开始下载jsp2class.jar...")
		jsp2classURL := "https://raw.githubusercontent.com/yezere/codeql_n1ght_dp/refs/heads/main/jsp2class.jar"
		if err := Common.DownloadFile(jsp2classURL, jsp2classPath); err != nil {
			return fmt.Errorf("下载jsp2class.jar失败: %v", err)
		}
		fmt.Printf("jsp2class.jar下载完成: %s\n", jsp2classPath)
	}

	return nil
}

// DownloadProcyon 保持向后兼容性
func DownloadProcyon() error {
	return DownloadDecompilers()
}
