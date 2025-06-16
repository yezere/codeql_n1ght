package Install

import (
	"codeql_n1ght/Common"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// CheckCodeQLInstalled 检查CodeQL是否已安装在tools目录下
func CheckCodeQLInstalled() bool {
	toolsDir := "./tools"
	codeqlPath := filepath.Join(toolsDir, "codeql")

	if _, err := os.Stat(codeqlPath); err == nil {
		fmt.Println("CodeQL 已经安装在 ./tools/codeql 目录下")
		return true
	}
	return false
}

// DownloadCodeQL 下载并安装CodeQL到tools目录
func DownloadCodeQL() error {
	if CheckCodeQLInstalled() {
		return nil
	}

	fmt.Println("开始下载CodeQL...")

	// 创建tools目录
	toolsDir := "./tools"
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		return fmt.Errorf("创建tools目录失败: %v", err)
	}

	// 根据操作系统选择下载链接
	var downloadURL string
	var fileName string

	// 优先使用用户指定的URL
	if Common.CodeQLDownloadURL != "" {
		downloadURL = Common.CodeQLDownloadURL
		fmt.Printf("使用用户指定的CodeQL下载地址: %s\n", downloadURL)
		// 从URL中提取文件名
		fileName = filepath.Base(downloadURL)
	} else {
		// 使用默认URL
		switch runtime.GOOS {
		case "windows":
			downloadURL = "https://github.com/github/codeql-cli-binaries/releases/latest/download/codeql-win64.zip"
			fileName = "codeql-win64.zip"
		case "linux":
			downloadURL = "https://github.com/github/codeql-cli-binaries/releases/latest/download/codeql-linux64.zip"
			fileName = "codeql-linux64.zip"
		case "darwin":
			downloadURL = "https://github.com/github/codeql-cli-binaries/releases/latest/download/codeql-osx64.zip"
			fileName = "codeql-osx64.zip"
		default:
			downloadURL = "https://github.com/github/codeql-cli-binaries/releases/latest/download/codeql-osx64.zip"
			fileName = "codeql-osx64.zip"
		}
		fmt.Printf("使用默认CodeQL下载地址: %s\n", downloadURL)
	}

	// 下载文件
	filePath := filepath.Join(toolsDir, fileName)
	if err := Common.DownloadFile(downloadURL, filePath); err != nil {
		return fmt.Errorf("下载CodeQL失败: %v", err)
	}

	fmt.Printf("CodeQL下载完成: %s\n", filePath)

	// 自动解压
	codeqlDir := filepath.Join(toolsDir, "codeql")
	if err := ExtractInstallZip(filePath, codeqlDir); err != nil {
		return fmt.Errorf("解压CodeQL失败: %v", err)
	}

	fmt.Println("CodeQL解压完成")

	// 删除下载的压缩包
	Common.RemoveFile(filePath)

	return nil
}
