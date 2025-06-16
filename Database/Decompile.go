package Database

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// UnzipJar 解压JAR文件到指定目录
func UnzipJar(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// 创建目标目录
	os.MkdirAll(dest, 0755)

	// 解压文件
	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)
		
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.FileInfo().Mode())
			continue
		}

		// 创建文件目录
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		// 打开压缩文件中的文件
		rc, err := f.Open()
		if err != nil {
			return err
		}

		// 创建目标文件
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
		if err != nil {
			rc.Close()
			return err
		}

		// 复制文件内容
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	fmt.Printf("Successfully extracted %s to %s\n", src, dest)
	return nil
}