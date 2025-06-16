package Install

import (
	"codeql_n1ght/Common"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// CheckJDKInstalled 检查JDK是否已安装在tools目录下
func CheckJDKInstalled() bool {
	toolsDir := "./tools"
	jdkPath := filepath.Join(toolsDir, "jdk")

	if _, err := os.Stat(jdkPath); err == nil {
		fmt.Println("JDK 已经安装在 ./tools/jdk 目录下")
		return true
	}
	return false
}

// DownloadJDK 下载并安装JDK8到tools目录
func DownloadJDK() error {
	if CheckJDKInstalled() {
		return nil
	}

	fmt.Println("开始下载JDK8...")

	// 创建tools目录
	toolsDir := "./tools"
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		return fmt.Errorf("创建tools目录失败: %v", err)
	}

	// 根据操作系统选择下载链接
	var downloadURL string
	var fileName string

	// 优先使用用户指定的URL
	if Common.JDKDownloadURL != "" {
		downloadURL = Common.JDKDownloadURL
		fmt.Printf("使用用户指定的JDK下载地址: %s\n", downloadURL)
	} else {
		// 使用默认URL
		switch runtime.GOOS {
		case "windows":
			// Windows JDK17
			downloadURL = "https://github.com/adoptium/temurin8-binaries/releases/download/jdk8u392-b08/OpenJDK8U-jdk_x64_windows_hotspot_8u392b08.zip"
		case "linux":
			downloadURL = "https://github.com/adoptium/temurin8-binaries/releases/download/jdk8u392-b08/OpenJDK8U-jdk_x64_linux_hotspot_8u392b08.tar.gz"
		case "darwin":
			downloadURL = "https://github.com/adoptium/temurin8-binaries/releases/download/jdk8u392-b08/OpenJDK8U-jdk_x64_mac_hotspot_8u392b08.tar.gz"
		default:
			return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
		}
		fmt.Printf("使用默认JDK下载地址: %s\n", downloadURL)
	}

	switch runtime.GOOS {
	case "windows":
		fileName = "OpenJDK8U-jdk_x64_windows_hotspot_8u392b08.zip"
	case "linux":
		fileName = "OpenJDK8U-jdk_x64_linux_hotspot_8u392b08.tar.gz"
	case "darwin":
		fileName = "OpenJDK8U-jdk_x64_mac_hotspot_8u392b08.tar.gz"
	}

	// 下载文件
	filePath := filepath.Join(toolsDir, fileName)
	if err := Common.DownloadFile(downloadURL, filePath); err != nil {
		return fmt.Errorf("下载JDK失败: %v", err)
	}

	fmt.Printf("JDK下载完成: %s\n", filePath)

	// 自动解压
	jdkDir := filepath.Join(toolsDir, "jdk")
	if err := ExtractInstallZip(filePath, jdkDir); err != nil {
		return fmt.Errorf("解压JDK失败: %v", err)
	}

	fmt.Println("JDK解压完成")

	// 删除下载的压缩包
	Common.RemoveFile(filePath)

	return nil
}
