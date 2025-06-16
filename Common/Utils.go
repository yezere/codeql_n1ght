package Common

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

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
