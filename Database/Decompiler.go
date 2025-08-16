package Database

import (
	"archive/zip"
	"bufio"
	"codeql_n1ght/Common"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

// DecompileJava 反编译Java文件
func DecompileJava(args ...string) error {
	cmd := exec.Command("java", args...)
	// 获取标准输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("获取 StdoutPipe 失败: %v", err)
	}
	// 获取标准错误管道（可选）
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("获取 StderrPipe 失败: %v", err)
	}
	// 启动命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动命令失败: %v", err)
	}
	// 创建协程并发读取标准输出
	go streamOutput(stdout, "STDOUT")
	// 创建协程并发读取标准错误
	go streamOutput(stderr, "STDERR")
	// 等待命令结束
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("命令执行异常: %v", err)
	}
	return nil
}

// copyDir 复制目录
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算目标路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			// 创建目录
			return os.MkdirAll(dstPath, info.Mode())
		} else {
			// 复制文件
			return copyFile(path, dstPath)
		}
	})
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	// 确保目标目录存在
	err := os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return err
	}

	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	return err
}



// DecompileLibraries 反编译依赖库，允许用户选择
func DecompileLibraries(location string) {
	// 优先检查BOOT-INF/lib目录（Spring Boot结构）
	libDir := filepath.Join(location, "output", "BOOT-INF", "lib")
	
	// 如果BOOT-INF/lib不存在，检查传统的WEB-INF/lib目录
	if _, err := os.Stat(libDir); os.IsNotExist(err) {
		libDir = filepath.Join(location, "output", "WEB-INF", "lib")
		if _, err := os.Stat(libDir); os.IsNotExist(err) {
			// 最后检查根目录下的lib目录
			libDir = filepath.Join(location, "output", "lib")
			if _, err := os.Stat(libDir); os.IsNotExist(err) {
				fmt.Println("No lib directory found (checked BOOT-INF/lib, WEB-INF/lib, and lib), skipping jar decompilation.")
				return
			}
		}
	}
	
	fmt.Printf("Using lib directory: %s\n", libDir)

	// 查找所有.jar文件
	jarFiles, err := filepath.Glob(filepath.Join(libDir, "*.jar"))
	if err != nil {
		fmt.Printf("Error searching for jar files: %v\n", err)
		return
	}

	if len(jarFiles) == 0 {
		fmt.Println("No jar files found in lib directory.")
		return
	}

	// 准备选项列表（只显示文件名）
	options := make([]string, len(jarFiles))
	for i, jarFile := range jarFiles {
		options[i] = filepath.Base(jarFile)
	}

	// 清空当前输出
	// clearScreen()

	// 使用survey的MultiSelect进行交互式选择
	var selectedFiles []string
	prompt := &survey.MultiSelect{
		Message:  "Select jar files to decompile (use arrow keys to navigate, space to select/deselect, enter to confirm):",
		Options:  options,
		PageSize: 40, // 每页显示40个选项
	}

	err = survey.AskOne(prompt, &selectedFiles)
	if err != nil {
		fmt.Printf("Error during selection: %v\n", err)
		return
	}

	if len(selectedFiles) == 0 {
		fmt.Println("No files selected, skipping jar decompilation.")
		return
	}

	// 反编译选中的文件
	fmt.Printf("\nDecompiling %d selected jar files...\n", len(selectedFiles))
	
	if Common.UseGoroutine {
		// 使用goroutine并发反编译
		decompileWithGoroutines(selectedFiles, jarFiles, location)
	} else {
		// 串行反编译
		for _, selectedFile := range selectedFiles {
			// 找到完整路径
			for _, jarFile := range jarFiles {
				if filepath.Base(jarFile) == selectedFile {
					fmt.Printf("Decompiling %s...\n", selectedFile)
					outputDir := filepath.Join(location, "createdabase", "src1")
					decompileJarFile(jarFile, outputDir, selectedFile)
					break
				}
			}
		}
	}
	fmt.Println("Jar decompilation completed.")
}

// 实时流式打印输出
func streamOutput(reader io.ReadCloser, prefix string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", prefix, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("[%s] 读取出错: %v\n", prefix, err)
	}
}

// clearScreen 清空控制台屏幕
func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default: // linux, darwin, etc.
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// decompileWithProcyon 使用Procyon反编译器
func decompileWithProcyon(jarFile, outputDir string) error {
	return DecompileJava("-jar", "tools/procyon-decompiler-0.6.0.jar", jarFile, "-o", outputDir)
}

// decompileWithFernflower 使用Fernflower反编译器
func decompileWithFernflower(jarFile, outputDir string) error {
	// 使用fernflower反编译jar文件
	err := DecompileJava("-cp", "tools/java-decompiler.jar",
		"org.jetbrains.java.decompiler.main.decompiler.ConsoleDecompiler",
		"-dgs=true",
		jarFile, outputDir)
	if err != nil {
		return err
	}

	// 反编译完成后，解压生成的jar文件
	decompiledJar := filepath.Join(outputDir, filepath.Base(jarFile))
	if _, err := os.Stat(decompiledJar); err == nil {
		extractDir := filepath.Join(outputDir)
		if err := extractJar(decompiledJar, extractDir); err != nil {
			return fmt.Errorf("解压反编译后的jar文件失败: %v", err)
		} else {
			// 删除反编译生成的jar文件
			os.Remove(decompiledJar)
			fmt.Printf("反编译和解压完成，源码位于: %s\n", extractDir)
		}
	}
	return nil
}

// extractJar 解压jar文件
func extractJar(jarFile, destDir string) error {
	// 创建目标目录
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 打开zip文件
	r, err := zip.OpenReader(jarFile)
	if err != nil {
		return fmt.Errorf("打开jar文件失败: %v", err)
	}
	defer r.Close()

	// 解压文件
	for _, f := range r.File {
		// 构建完整的文件路径
		path := filepath.Join(destDir, f.Name)

		// 确保路径安全
		if !strings.HasPrefix(path, filepath.Clean(destDir)+string(os.PathSeparator)) {
			continue
		}

		if f.FileInfo().IsDir() {
			// 创建目录
			os.MkdirAll(path, f.FileInfo().Mode())
			continue
		}

		// 创建文件的父目录
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("创建父目录失败: %v", err)
		}

		// 打开zip文件中的文件
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("打开zip文件中的文件失败: %v", err)
		}

		// 创建目标文件
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
		if err != nil {
			rc.Close()
			return fmt.Errorf("创建目标文件失败: %v", err)
		}

		// 复制文件内容
		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return fmt.Errorf("复制文件内容失败: %v", err)
		}
	}

	fmt.Printf("jar文件解压完成: %s\n", destDir)
	return nil
}

// decompileJarFile 反编译单个jar文件
func decompileJarFile(jarFile, outputDir, selectedFile string) {
	var err error
	// 根据反编译器类型选择不同的反编译方式
	switch Common.DecompilerType {
	case "fernflower":
		err = decompileWithFernflower(jarFile, outputDir)
		if err != nil {
			color.Red("Fernflower反编译失败: %v，切换到Procyon反编译器\n", err)
			err = decompileWithProcyon(jarFile, outputDir)
			if err != nil {
				color.Red("Procyon反编译也失败: %v\n", err)
			} else {
				fmt.Printf("使用Procyon反编译器成功完成 %s\n", selectedFile)
			}
		}
	default: // procyon
		err = decompileWithProcyon(jarFile, outputDir)
		if err != nil {
			color.Red("Procyon反编译失败: %v，切换到Fernflower反编译器\n", err)
			err = decompileWithFernflower(jarFile, outputDir)
			if err != nil {
				color.Red("Fernflower反编译也失败: %v\n", err)
			} else {
				fmt.Printf("使用Fernflower反编译器成功完成 %s\n", selectedFile)
			}
		}
	}
}

// decompileWithGoroutines 使用goroutine并发反编译
func decompileWithGoroutines(selectedFiles, jarFiles []string, location string) {
	// 创建工作队列
	type DecompileTask struct {
		jarFile      string
		selectedFile string
		outputDir    string
	}

	tasks := make(chan DecompileTask, len(selectedFiles))
	var wg sync.WaitGroup

	// 启动worker goroutines
	maxWorkers := Common.MaxGoroutines
	if maxWorkers <= 0 {
		maxWorkers = 4 // 默认值
	}

	fmt.Printf("启动 %d 个goroutine进行并发反编译...\n", maxWorkers)

	// 启动worker
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for task := range tasks {
				fmt.Printf("[Worker %d] Decompiling %s...\n", workerID, task.selectedFile)
				decompileJarFile(task.jarFile, task.outputDir, task.selectedFile)
				fmt.Printf("[Worker %d] Completed %s\n", workerID, task.selectedFile)
			}
		}(i)
	}

	// 发送任务到队列
	for _, selectedFile := range selectedFiles {
		// 找到完整路径
		for _, jarFile := range jarFiles {
			if filepath.Base(jarFile) == selectedFile {
				outputDir := filepath.Join(location, "createdabase", "src1")
				tasks <- DecompileTask{
					jarFile:      jarFile,
					selectedFile: selectedFile,
					outputDir:    outputDir,
				}
				break
			}
		}
	}

	// 关闭任务队列
	close(tasks)

	// 等待所有goroutine完成
	wg.Wait()
	fmt.Println("所有goroutine反编译任务完成")
}
