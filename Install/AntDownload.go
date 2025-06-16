package Install

import (
	"codeql_n1ght/Common"
	"fmt"
	"os"
	"path/filepath"
)

// CheckAntInstalled 检查Apache Ant是否已安装在tools目录下
func CheckAntInstalled() bool {
	toolsDir := "./tools"
	antPath := filepath.Join(toolsDir, "ant")

	if _, err := os.Stat(antPath); err == nil {
		fmt.Println("Apache Ant 已经安装在 ./tools/ant 目录下")
		return true
	}
	return false
}

// DownloadAnt 下载并安装Apache Ant到tools目录
func DownloadAnt() error {
	if CheckAntInstalled() {
		return nil
	}

	fmt.Println("开始下载Apache Ant...")

	// 创建tools目录
	toolsDir := "./tools"
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		return fmt.Errorf("创建tools目录失败: %v", err)
	}

	// Apache Ant下载链接（跨平台通用）// 下载Apache Ant
	var downloadURL string
	if Common.AntDownloadURL != "" {
		downloadURL = Common.AntDownloadURL
		fmt.Printf("使用用户指定的Apache Ant下载地址: %s\n", downloadURL)
	} else {
		downloadURL = "https://archive.apache.org/dist/ant/binaries/apache-ant-1.10.14-bin.zip"
		fmt.Printf("使用默认Apache Ant下载地址: %s\n", downloadURL)
	}

	fileName := "apache-ant-1.10.14-bin.zip"
	filePath := filepath.Join(toolsDir, fileName)
	if err := Common.DownloadFile(downloadURL, filePath); err != nil {
		return fmt.Errorf("下载Apache Ant失败: %v", err)
	}

	fmt.Printf("Apache Ant下载完成: %s\n", filePath)

	// 自动解压
	antDir := filepath.Join(toolsDir, "ant")
	if err := ExtractInstallZip(filePath, antDir); err != nil {
		return fmt.Errorf("解压Apache Ant失败: %v", err)
	}

	fmt.Println("Apache Ant解压完成")

	// 删除下载的压缩包
	Common.RemoveFile(filePath)

	return nil
}

// InstallAllTools 安装所有工具的便捷函数
func InstallAllTools() error {
	fmt.Println("=== 开始安装开发工具 ===")

	// 检查并安装JDK8
	fmt.Println("\n1. 检查JDK8...")
	if err := DownloadJDK(); err != nil {
		fmt.Printf("JDK安装失败: %v\n", err)
	}

	// 检查并安装CodeQL
	fmt.Println("\n2. 检查CodeQL...")
	if err := DownloadCodeQL(); err != nil {
		fmt.Printf("CodeQL安装失败: %v\n", err)
	}

	// 检查并安装Apache Ant
	fmt.Println("\n3. 检查Apache Ant...")
	if err := DownloadAnt(); err != nil {
		fmt.Printf("Apache Ant安装失败: %v\n", err)
	}

	// 检查并安装Procyon
	fmt.Println("\n4. 检查Procyon...")
	if err := DownloadProcyon(); err != nil {
		fmt.Printf("Procyon安装失败: %v\n", err)
	}

	// 检查并安装Apache Tomcat
	fmt.Println("\n5. 检查Apache Tomcat...")
	if err := DownloadTomcat(); err != nil {
		fmt.Printf("Apache Tomcat安装失败: %v\n", err)
	}

	fmt.Println("\n=== 工具安装检查完成 ===")
	return nil
}
