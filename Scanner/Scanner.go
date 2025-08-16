package Scanner

import (
	"codeql_n1ght/Common"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

// ScanResult 扫描结果结构
type ScanResult struct {
	QueryFile string
	Success   bool
	Output    string
	Error     error
	Duration  time.Duration
}

// RunScan 执行CodeQL扫描
func RunScan() error {
	// 显示带边框的扫描开始提示
	displayScanHeader()

	// 清理之前的结果文件
	if err := cleanupPreviousResults(); err != nil {
		Common.LogWarn("清理之前的结果文件失败: %v", err)
	}

	// 检查并解压源码文件
	if err := extractSourceFiles(); err != nil {
		Common.LogWarn("解压源码文件失败: %v", err)
	}

	// 显示扫描配置信息
	Common.LogInfo("数据库路径: %s", Common.DatabasePath)
	Common.LogInfo("QL库路径: %s", Common.QLLibsPath)

	// 验证扫描相关目录
	if err := validateScanDirectory(); err != nil {
		return err
	}

	// 获取目录下的所有.ql文件
	qlFiles, err := findQLFiles()
	if err != nil {
		return err
	}

	if len(qlFiles) == 0 {
		Common.LogWarn("未找到任何.ql文件")
		return fmt.Errorf("QL库目录 %s 中未找到.ql文件", Common.QLLibsPath)
	}

	Common.LogInfo("找到 %d 个查询文件", len(qlFiles))

	// 执行查询
	results := make([]ScanResult, 0, len(qlFiles))
	if Common.UseGoroutine {
		results = executeConcurrentQueries(qlFiles)
	} else {
		results = executeSequentialQueries(qlFiles)
	}

	// 显示扫描总结
	displayScanSummary(results)



	return nil
}

// displayScanHeader 显示扫描开始的界面
func displayScanHeader() {
	scannerFigure := figure.NewFigure("Starting Scan", "", true)
	asciiArt := scannerFigure.String()
	lines := strings.Split(asciiArt, "\n")

	// 计算最长行长度
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// 打印带边框的绿色提示
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Println(strings.Repeat("-", maxLen+4))
	for _, line := range lines {
		fmt.Printf("| %s%s |\n", green(line), strings.Repeat(" ", maxLen-len(line)))
	}
	fmt.Println(strings.Repeat("-", maxLen+4))
	fmt.Println()
}

// validateScanDirectory 验证扫描相关目录
func validateScanDirectory() error {
	// 验证数据库路径
	if !Common.IsDirectory(Common.DatabasePath) {
		return fmt.Errorf("指定的数据库路径不是有效目录: %s", Common.DatabasePath)
	}

	// 验证QL库路径
	if !Common.IsDirectory(Common.QLLibsPath) {
		return fmt.Errorf("指定的QL库路径不是有效目录: %s", Common.QLLibsPath)
	}

	return nil
}

// findQLFiles 查找所有.ql文件
func findQLFiles() ([]string, error) {
	var qlFiles []string
	err := filepath.Walk(Common.QLLibsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".ql") {
			qlFiles = append(qlFiles, path)
		}
		return nil
	})
	return qlFiles, err
}

// executeConcurrentQueries 并发执行查询
func executeConcurrentQueries(qlFiles []string) []ScanResult {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, Common.MaxGoroutines)
	results := make(chan ScanResult, len(qlFiles))

	Common.LogInfo("使用并发模式执行查询 (最大并发数: %d)", Common.MaxGoroutines)

	for _, qlFile := range qlFiles {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量
		go executeQuery(qlFile, &wg, semaphore, results)
	}

	// 等待所有查询完成
	wg.Wait()
	close(results)

	// 收集结果
	var scanResults []ScanResult
	for result := range results {
		scanResults = append(scanResults, result)
	}

	return scanResults
}

// executeSequentialQueries 顺序执行查询
func executeSequentialQueries(qlFiles []string) []ScanResult {
	var results []ScanResult

	Common.LogInfo("使用顺序模式执行查询")

	for _, qlFile := range qlFiles {
		startTime := time.Now()
		result := ScanResult{
			QueryFile: qlFile,
			Success:   false,
		}

		Common.LogInfo("正在执行查询: %s", filepath.Base(qlFile))
		Common.SetupEnvironment()

		// 构建CodeQL命令
		cmd := exec.Command("codeql", "database", "analyze",
			Common.DatabasePath, // 数据库路径
			qlFile,              // 查询文件
			fmt.Sprintf("--threads=%d", Common.CodeQLThreads),
			"--format=sarifv2.1.0",
			"--output=results.sarif",
		)

		// 执行命令并获取输出
		output, err := cmd.CombinedOutput()
		result.Duration = time.Since(startTime)
		result.Output = string(output)

		if err != nil {
			result.Error = err
			Common.LogError("执行查询 %s 失败 (耗时: %v): %v", filepath.Base(qlFile), result.Duration, err)
			if len(result.Output) > 0 {
				Common.LogError("错误输出: %s", result.Output)
			}
			showPackInstallHint()
		} else {
			result.Success = true
			Common.LogInfo("查询 %s 完成 (耗时: %v)", filepath.Base(qlFile), result.Duration)
			if len(result.Output) > 0 {
				color.White("查询结果:\n%s", result.Output)
			}
		}

		results = append(results, result)
	}

	return results
}

// executeQuery 执行单个查询
func executeQuery(qlFile string, wg *sync.WaitGroup, semaphore chan struct{}, results chan<- ScanResult) {
	defer wg.Done()
	defer func() { <-semaphore }() // 释放信号量

	startTime := time.Now()
	result := ScanResult{
		QueryFile: qlFile,
		Success:   false,
	}

	Common.LogInfo("正在执行查询: %s", filepath.Base(qlFile))
	Common.SetupEnvironment()
	// 构建CodeQL命令
	cmd := exec.Command("codeql", "database", "analyze",
		Common.DatabasePath, // 数据库路径
		qlFile,              // 查询文件
		fmt.Sprintf("--threads=%d", Common.CodeQLThreads),
		"--format=sarifv2.1.0",
		"--output=results.sarif",
	)

	// 执行命令并获取输出
	output, err := cmd.CombinedOutput()
	result.Duration = time.Since(startTime)
	result.Output = string(output)

	if err != nil {
		result.Error = err
		Common.LogError("执行查询 %s 失败 (耗时: %v): %v", filepath.Base(qlFile), result.Duration, err)
		if len(result.Output) > 0 {
			Common.LogError("错误输出: %s", result.Output)
		}
		showPackInstallHint()
	} else {
		result.Success = true
		Common.LogInfo("查询 %s 完成 (耗时: %v)", filepath.Base(qlFile), result.Duration)
		if len(result.Output) > 0 {
			color.White("查询结果:\n%s", result.Output)
		}
	}

	results <- result
}

// displayScanSummary 显示扫描总结
func displayScanSummary(results []ScanResult) {
	successCount := 0
	totalDuration := time.Duration(0)

	for _, result := range results {
		if result.Success {
			successCount++
		}
		totalDuration += result.Duration
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	Common.LogInfo("扫描总结:")
	Common.LogInfo("总查询数: %d", len(results))
	color.Green("成功: %d", successCount)
	color.Red("失败: %d", len(results)-successCount)
	Common.LogInfo("总耗时: %v", totalDuration)
	fmt.Println(strings.Repeat("=", 60))
}
