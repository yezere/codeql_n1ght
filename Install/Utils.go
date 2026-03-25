package Install

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractInstallZip 解压ZIP文件并移除顶层路径，专门用于Install模块
// 这个函数会自动移除ZIP文件中的顶层目录，将内容直接解压到目标目录
func ExtractInstallZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	os.MkdirAll(dest, 0755)

	// 找到顶层目录名
	var topLevelDir string
	for _, f := range r.File {
		if f.FileInfo().IsDir() && !strings.Contains(f.Name, "/") {
			topLevelDir = f.Name
			break
		}
	}

	// 如果没有找到顶层目录，尝试从第一个文件路径中提取
	if topLevelDir == "" && len(r.File) > 0 {
		firstPath := r.File[0].Name
		if idx := strings.Index(firstPath, "/"); idx != -1 {
			topLevelDir = firstPath[:idx+1]
		}
	}

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// 移除顶层目录路径
		path := f.Name
		if topLevelDir != "" && strings.HasPrefix(path, topLevelDir) {
			path = strings.TrimPrefix(path, topLevelDir)
			// 如果移除顶层目录后路径为空，跳过
			if path == "" {
				continue
			}
		}

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

// ExtractInstallZipWithProgress 带进度显示的解压函数
func ExtractInstallZipWithProgress(src, dest string) error {
	fmt.Printf("正在解压 %s 到 %s...\n", filepath.Base(src), dest)
	var err error
	if strings.HasSuffix(strings.ToLower(src), ".tar.gz") || strings.HasSuffix(strings.ToLower(src), ".tgz") {
		err = ExtractInstallTarGz(src, dest)
	} else {
		err = ExtractInstallZip(src, dest)
	}
	if err != nil {
		return fmt.Errorf("解压失败: %v", err)
	}
	// 给bin目录下的.sh文件添加执行权限
	binDir := filepath.Join(dest, "bin")
	if _, err := os.Stat(binDir); err == nil {
		filepath.Walk(binDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".sh") {
				os.Chmod(path, 0755)
			}
			return nil
		})
	}
	fmt.Println("解压完成")
	return nil
}

// ExtractInstallTarGz 解压tar.gz文件
func ExtractInstallTarGz(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	os.MkdirAll(dest, 0755)

	tr := tar.NewReader(gzr)

	var topLevelDir string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if topLevelDir == "" {
			if idx := strings.Index(hdr.Name, "/"); idx != -1 {
				topLevelDir = hdr.Name[:idx+1]
			}
		}

		path := hdr.Name
		if topLevelDir != "" && strings.HasPrefix(path, topLevelDir) {
			path = strings.TrimPrefix(path, topLevelDir)
			if path == "" {
				continue
			}
		}

		fpath := filepath.Join(dest, path)

		if hdr.Typeflag == tar.TypeDir {
			os.MkdirAll(fpath, os.FileMode(hdr.Mode))
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(hdr.Mode))
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, tr)
		outFile.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
