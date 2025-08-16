package Scanner

import (
	"archive/zip"
	"codeql_n1ght/Common"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 全局变量存储源码根目录路径
var sourceRootPath string

// extractSourceFiles 检查并解压源码文件
func extractSourceFiles() error {
	srcZipPath := filepath.Join(Common.DatabasePath, "src.zip")
	srcDir := filepath.Join(Common.DatabasePath, "src")

	// 检查src.zip是否存在
	if _, err := os.Stat(srcZipPath); os.IsNotExist(err) {
		Common.LogInfo("未找到src.zip文件，跳过源码解压")
		return nil
	}

	// 检查src目录是否已存在
	if _, err := os.Stat(srcDir); err == nil {
		Common.LogInfo("src目录已存在，跳过解压")
		// 探测并缓存源码根目录路径
		detectSourceRootPath()
		return nil
	}

	Common.LogInfo("正在解压源码文件: %s", srcZipPath)

	// 打开zip文件
	reader, err := zip.OpenReader(srcZipPath)
	if err != nil {
		return fmt.Errorf("无法打开zip文件: %v", err)
	}
	defer reader.Close()

	// 创建目标目录
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		return fmt.Errorf("无法创建目录: %v", err)
	}

	// 解压文件
	for _, file := range reader.File {
		if err := extractFile(file, srcDir); err != nil {
			return fmt.Errorf("解压文件 %s 失败: %v", file.Name, err)
		}
	}

	Common.LogInfo("源码解压完成到: %s", srcDir)

	// 探测并缓存源码根目录路径
	detectSourceRootPath()
	return nil
}

// detectSourceRootPath 探测源码根目录路径
func detectSourceRootPath() {
	srcDir := filepath.Join(Common.DatabasePath, "src")

	// 递归查找包含src1目录的路径
	var findSrc1 func(string) string
	findSrc1 = func(dir string) string {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return ""
		}

		for _, entry := range entries {
			if entry.IsDir() {
				entryPath := filepath.Join(dir, entry.Name())
				if entry.Name() == "src1" {
					// 找到src1目录，返回其父目录
					return filepath.Dir(entryPath)
				}
				// 递归搜索子目录
				if result := findSrc1(entryPath); result != "" {
					return result
				}
			}
		}
		return ""
	}

	if rootPath := findSrc1(srcDir); rootPath != "" {
		// 计算相对于src目录的路径
		if relPath, err := filepath.Rel(srcDir, rootPath); err == nil {
			sourceRootPath = relPath
			Common.LogInfo("检测到源码根目录: %s", filepath.Join(srcDir, sourceRootPath))
		}
	} else {
		Common.LogWarn("未能检测到源码根目录，使用默认路径")
		sourceRootPath = ""
	}
}

// extractFile 解压单个文件
func extractFile(file *zip.File, destDir string) error {
	// 构建目标路径
	destPath := filepath.Join(destDir, file.Name)

	// 确保目标路径在目标目录内（安全检查）
	if !strings.HasPrefix(destPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
		return fmt.Errorf("无效的文件路径: %s", file.Name)
	}

	// 如果是目录，创建目录
	if file.FileInfo().IsDir() {
		return os.MkdirAll(destPath, file.FileInfo().Mode())
	}

	// 创建文件的父目录
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// 打开zip中的文件
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 复制文件内容
	_, err = io.Copy(destFile, srcFile)
	return err
}

// GetSourceRootPath 获取源码根目录路径（供其他模块使用）
func GetSourceRootPath() string {
	return sourceRootPath
}
