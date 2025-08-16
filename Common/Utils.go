package Common

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// LogLevel 日志级别
type LogLevel int

const (
	LogLevelInfo LogLevel = iota
	LogLevelWarn
	LogLevelError
)

// LogMessage 统一的日志输出函数
func LogMessage(level LogLevel, format string, args ...interface{}) {
	switch level {
	case LogLevelInfo:
		color.White(format, args...)
	case LogLevelWarn:
		color.Yellow(format, args...)
	case LogLevelError:
		color.Red(format, args...)
	}
}

// LogError 记录错误信息
func LogError(format string, args ...interface{}) {
	LogMessage(LogLevelError, format, args...)
}

// LogWarn 记录警告信息
func LogWarn(format string, args ...interface{}) {
	LogMessage(LogLevelWarn, format, args...)
}

// LogInfo 记录信息
func LogInfo(format string, args ...interface{}) {
	LogMessage(LogLevelInfo, format, args...)
}

// ValidateFile 验证文件是否存在且可读
func ValidateFile(filePath string) error {
	if !FileExists(filePath) {
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	return nil
}

// IsDirectory 检查路径是否为目录
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// DownloadFile 下载文件的通用函数（带进度条）
func DownloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 获取文件大小和文件名
	fileSize := resp.ContentLength
	filename := filepath[strings.LastIndex(filepath, "\\")+1:]

	if fileSize > 0 {
		// 使用第三方进度条库显示下载进度
		bar := progressbar.DefaultBytes(
			fileSize,
			fmt.Sprintf("下载 %s", filename),
		)
		_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	} else {
		// 如果无法获取文件大小，使用spinner模式
		bar := progressbar.DefaultBytes(-1, fmt.Sprintf("下载 %s", filename))
		_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	}

	return err
}

// ExtractZip 解压ZIP文件的通用函数，完全还原原始结构
func ExtractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	os.MkdirAll(dest, 0755)

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// 完全保留原始路径结构，不移除顶层目录
		path := f.Name
		fpath := filepath.Join(dest, path)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.FileInfo().Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// FileExists 检查文件或目录是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CreateDirIfNotExists 如果目录不存在则创建
func CreateDirIfNotExists(dir string) error {
	if !FileExists(dir) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// RemoveFile 安全删除文件
func RemoveFile(path string) error {
	if !FileExists(path) {
		return nil
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		// 如果是目录，递归删除整个目录
		return os.RemoveAll(path)
	} else {
		// 如果是文件，删除文件
		return os.Remove(path)
	}
}

// SafeExecute 安全执行函数，捕获panic
func SafeExecute(fn func() error, errorMsg string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("发生panic: %v", r)
		}
	}()

	if err := fn(); err != nil {
		LogError("%s: %v", errorMsg, err)
		return err
	}
	return nil
}

// CopyFile 复制单个文件
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 确保目标目录存在
	if err := CreateDirIfNotExists(filepath.Dir(dst)); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// 复制文件权限
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// CopyDirectory 递归复制整个目录
func CopyDirectory(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("源路径不是目录: %s", src)
	}

	// 创建目标目录
	if err := CreateDirIfNotExists(dst); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 递归复制子目录
			if err := CopyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 复制文件
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyExtraSourceToSrc1 将额外源码目录复制到src1目录
func CopyExtraSourceToSrc1(extraSourceDir, src1Dir string) error {
	if extraSourceDir == "" {
		return nil // 没有指定额外源码目录
	}

	LogInfo("开始复制额外源码目录: %s -> %s", extraSourceDir, src1Dir)

	// 复制所有文件和子目录到src1
	err := CopyDirectory(extraSourceDir, src1Dir)
	if err != nil {
		return fmt.Errorf("复制额外源码失败: %v", err)
	}

	LogInfo("额外源码复制完成")
	return nil
}
