package Scanner

import (
	"codeql_n1ght/Common"
	"os"
	"path/filepath"
)

// cleanupPreviousResults 清理之前的结果文件
func cleanupPreviousResults() error {
	// 定义需要清理的文件列表
	filesToClean := []string{
		"results.sarif",
		"scan_report.html",
	}

	for _, file := range filesToClean {
		if _, err := os.Stat(file); err == nil {
			// 文件存在，尝试删除
			if err := os.Remove(file); err != nil {
				Common.LogWarn("无法删除文件 %s: %v", file, err)
			} else {
				Common.LogInfo("已删除之前的结果文件: %s", file)
			}
		}
	}

	// 清理CodeQL缓存，确保修改的QL文件能生效
	cleanupCodeQLCache()

	// 可选：清理之前解压的src目录（如果用户想要重新解压）
	// 注释掉下面的代码以保留之前解压的文件，加快后续扫描速度
	/*
		srcDir := filepath.Join(Common.DatabasePath, "src")
		if _, err := os.Stat(srcDir); err == nil {
			if err := os.RemoveAll(srcDir); err != nil {
				Common.LogWarn("无法删除src目录 %s: %v", srcDir, err)
			} else {
				Common.LogInfo("已清理之前的src目录: %s", srcDir)
			}
		}
	*/

	return nil
}

// cleanupCodeQLCache 清理CodeQL缓存文件
func cleanupCodeQLCache() {
	// 只有在用户明确指定时才清理缓存
	if !Common.CleanCache {
		return
	}

	Common.LogInfo("开始清理CodeQL缓存...")

	// 清理数据库缓存目录
	cacheDir := filepath.Join(Common.DatabasePath, "cache")
	if _, err := os.Stat(cacheDir); err == nil {
		if err := os.RemoveAll(cacheDir); err != nil {
			Common.LogWarn("无法删除缓存目录 %s: %v", cacheDir, err)
		} else {
			Common.LogInfo("已清理CodeQL缓存目录: %s", cacheDir)
		}
	}

	// 清理查询结果缓存目录（完全删除以确保重新运行）
	resultsDir := filepath.Join(Common.DatabasePath, "results")
	if _, err := os.Stat(resultsDir); err == nil {
		if err := os.RemoveAll(resultsDir); err != nil {
			Common.LogWarn("无法删除results目录 %s: %v", resultsDir, err)
		} else {
			Common.LogInfo("已清理查询结果缓存目录: %s", resultsDir)
		}
	}

	// 清理全局CodeQL缓存（如果存在）
	homeDir, err := os.UserHomeDir()
	if err == nil {
		globalCacheDir := filepath.Join(homeDir, ".codeql", "cache")
		if _, err := os.Stat(globalCacheDir); err == nil {
			Common.LogInfo("检测到全局CodeQL缓存目录: %s，建议手动清理以获得最佳效果", globalCacheDir)
		}
	}

	Common.LogInfo("缓存清理完成，这将确保修改的QL文件生效")
}
