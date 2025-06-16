package Install

import (
	"codeql_n1ght/Common"
	"fmt"
	"os"
	"path/filepath"
)

// CheckTomcatInstalled 检查Apache Tomcat是否已安装在tools目录下
func CheckTomcatInstalled() bool {
	toolsDir := "./tools"
	tomcatPath := filepath.Join(toolsDir, "tomcat")

	if _, err := os.Stat(tomcatPath); err == nil {
		fmt.Println("Apache Tomcat 已经安装在 ./tools/tomcat 目录下")
		return true
	}
	return false
}

// DownloadTomcat 下载并安装Apache Tomcat到tools目录
func DownloadTomcat() error {
	if CheckTomcatInstalled() {
		return nil
	}

	fmt.Println("开始下载Apache Tomcat...")

	// 创建tools目录
	toolsDir := "./tools"
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		return fmt.Errorf("创建tools目录失败: %v", err)
	}

	// Apache Tomcat 9.0.27下载链接
	downloadURL := "https://archive.apache.org/dist/tomcat/tomcat-9/v9.0.27/bin/apache-tomcat-9.0.27.zip"
	fileName := "apache-tomcat-9.0.27.zip"

	// 下载文件
	filePath := filepath.Join(toolsDir, fileName)
	if err := Common.DownloadFile(downloadURL, filePath); err != nil {
		return fmt.Errorf("下载Apache Tomcat失败: %v", err)
	}

	fmt.Printf("Apache Tomcat下载完成: %s\n", filePath)

	// 自动解压
	tomcatDir := filepath.Join(toolsDir, "tomcat")
	if err := ExtractInstallZip(filePath, tomcatDir); err != nil {
		return fmt.Errorf("解压Apache Tomcat失败: %v", err)
	}

	fmt.Printf("Apache Tomcat解压完成: %s\n", tomcatDir)

	// 删除下载的zip文件
	if err := os.Remove(filePath); err != nil {
		fmt.Printf("警告: 删除下载文件失败: %v\n", err)
	} else {
		fmt.Println("已删除下载的zip文件")
	}

	fmt.Println("Apache Tomcat安装完成")
	return nil
}

// InstallTomcat 安装Apache Tomcat的便捷函数
func InstallTomcat() error {
	fmt.Println("=== 安装Apache Tomcat ===")
	return DownloadTomcat()
}

// CheckTomcatAvailability 检查Tomcat可用性
func CheckTomcatAvailability() bool {
	toolsDir := "./tools"
	tomcatPath := filepath.Join(toolsDir, "tomcat")

	// 检查tomcat目录是否存在
	if _, err := os.Stat(tomcatPath); err != nil {
		return false
	}

	// 检查关键文件是否存在
	tomcatVersionPath := filepath.Join(tomcatPath, "apache-tomcat-9.0.27")
	if _, err := os.Stat(tomcatVersionPath); err == nil {
		// 检查bin目录和startup脚本
		binPath := filepath.Join(tomcatVersionPath, "bin")
		if _, err := os.Stat(binPath); err == nil {
			return true
		}
	}

	return false
}

// GetTomcatPath 获取Tomcat安装路径
func GetTomcatPath() string {
	toolsDir := "./tools"
	tomcatVersionPath := filepath.Join(toolsDir, "tomcat", "apache-tomcat-9.0.27")

	if _, err := os.Stat(tomcatVersionPath); err == nil {
		absPath, _ := filepath.Abs(tomcatVersionPath)
		return absPath
	}

	return ""
}
