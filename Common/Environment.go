package Common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// SetupEnvironment 设置所有工具的环境变量
func SetupEnvironment() error {
	fmt.Println("正在设置环境变量...")

	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	toolsDir := filepath.Join(cwd, "tools")

	// 设置JDK环境变量
	if err := setupJDKEnvironment(toolsDir); err != nil {
		fmt.Printf("设置JDK环境变量失败: %v\n", err)
	}

	// 设置CodeQL环境变量
	if err := setupCodeQLEnvironment(toolsDir); err != nil {
		fmt.Printf("设置CodeQL环境变量失败: %v\n", err)
	}

	// 设置Ant环境变量
	if err := setupAntEnvironment(toolsDir); err != nil {
		fmt.Printf("设置Ant环境变量失败: %v\n", err)
	}

	// 设置Tomcat环境变量
	if err := setupTomcatEnvironment(toolsDir); err != nil {
		fmt.Printf("设置Tomcat环境变量失败: %v\n", err)
	}

	fmt.Println("环境变量设置完成")
	return nil
}

// setupJDKEnvironment 设置JDK环境变量
func setupJDKEnvironment(toolsDir string) error {
	jdkPath := filepath.Join(toolsDir, "jdk")
	if _, err := os.Stat(jdkPath); err != nil {
		return fmt.Errorf("JDK未安装")
	}

	// 设置JAVA_HOME
	if err := os.Setenv("JAVA_HOME", jdkPath); err != nil {
		return fmt.Errorf("设置JAVA_HOME失败: %v", err)
	}
	// 设置JRE_HOME
	if err := os.Setenv("JRE_HOME", jdkPath); err != nil {
		return fmt.Errorf("设置JRE_HOME失败: %v", err)
	}
	// 添加到PATH
	jdkBinPath := filepath.Join(jdkPath, "bin")
	if err := addToPath(jdkBinPath); err != nil {
		return fmt.Errorf("添加JDK到PATH失败: %v", err)
	}

	fmt.Printf("JDK环境变量设置完成: JAVA_HOME=%s\n", jdkPath)
	return nil
}

// setupCodeQLEnvironment 设置CodeQL环境变量
func setupCodeQLEnvironment(toolsDir string) error {
	codeqlPath := filepath.Join(toolsDir, "codeql")
	if _, err := os.Stat(codeqlPath); err != nil {
		return fmt.Errorf("CodeQL未安装")
	}

	// 添加到PATH
	if err := addToPath(codeqlPath); err != nil {
		return fmt.Errorf("添加CodeQL到PATH失败: %v", err)
	}

	// 获取系统内存并计算一半
	memoryGB := getHalfSystemMemoryGB()
	color.Green("获取内存", memoryGB, "\n")
	jvmArgs := fmt.Sprintf("-Xmx%dg -Xms%dg", memoryGB, memoryGB)

	// 设置SEMMLE_JAVA_EXTRACTOR_JVM_ARGS
	if err := os.Setenv("SEMMLE_JAVA_EXTRACTOR_JVM_ARGS", jvmArgs); err != nil {
		return fmt.Errorf("设置SEMMLE_JAVA_EXTRACTOR_JVM_ARGS失败: %v", err)
	}

	fmt.Printf("CodeQL环境变量设置完成: PATH中已添加 %s\n", codeqlPath)
	fmt.Printf("SEMMLE_JAVA_EXTRACTOR_JVM_ARGS已设置为: %s\n", jvmArgs)
	return nil
}

// setupAntEnvironment 设置Ant环境变量
func setupAntEnvironment(toolsDir string) error {
	antPath := filepath.Join(toolsDir, "ant")
	if _, err := os.Stat(antPath); err != nil {
		return fmt.Errorf("Apache Ant未安装")
	}

	// 设置ANT_HOME
	if err := os.Setenv("ANT_HOME", antPath); err != nil {
		return fmt.Errorf("设置ANT_HOME失败: %v", err)
	}

	// 添加到PATH
	antBinPath := filepath.Join(antPath, "bin")
	if err := addToPath(antBinPath); err != nil {
		return fmt.Errorf("添加Ant到PATH失败: %v", err)
	}

	fmt.Printf("Ant环境变量设置完成: ANT_HOME=%s\n", antPath)
	return nil
}

// setupTomcatEnvironment 设置Tomcat环境变量
func setupTomcatEnvironment(toolsDir string) error {
	tomcatPath := filepath.Join(toolsDir, "tomcat")
	if _, err := os.Stat(tomcatPath); err != nil {
		return fmt.Errorf("Apache Tomcat未安装")
	}

	// 检查是否有apache-tomcat-9.0.27目录
	tomcatVersionPath := filepath.Join(tomcatPath, "apache-tomcat-9.0.27")
	if _, err := os.Stat(tomcatVersionPath); err == nil {
		// 使用具体版本目录作为CATALINA_HOME
		if err := os.Setenv("CATALINA_HOME", tomcatVersionPath); err != nil {
			return fmt.Errorf("设置CATALINA_HOME失败: %v", err)
		}
		fmt.Printf("Tomcat环境变量设置完成: CATALINA_HOME=%s\n", tomcatVersionPath)
	} else {
		// 使用tomcat目录作为CATALINA_HOME
		if err := os.Setenv("CATALINA_HOME", tomcatPath); err != nil {
			return fmt.Errorf("设置CATALINA_HOME失败: %v", err)
		}
		fmt.Printf("Tomcat环境变量设置完成: CATALINA_HOME=%s\n", tomcatPath)
	}

	// 添加Tomcat bin目录到PATH
	catalinaHome := os.Getenv("CATALINA_HOME")
	tomcatBinPath := filepath.Join(catalinaHome, "bin")
	if _, err := os.Stat(tomcatBinPath); err == nil {
		if err := addToPath(tomcatBinPath); err != nil {
			return fmt.Errorf("添加Tomcat到PATH失败: %v", err)
		}
	}

	return nil
}

// addToPath 将路径添加到PATH环境变量
func addToPath(newPath string) error {
	currentPath := os.Getenv("PATH")

	// 检查路径是否已存在
	var pathSeparator string
	if runtime.GOOS == "windows" {
		pathSeparator = ";"
	} else {
		pathSeparator = ":"
	}

	paths := strings.Split(currentPath, pathSeparator)
	for _, path := range paths {
		if strings.EqualFold(strings.TrimSpace(path), newPath) {
			// 路径已存在
			return nil
		}
	}

	// 添加新路径到PATH前面
	newPathValue := newPath + pathSeparator + currentPath
	return os.Setenv("PATH", newPathValue)
}

// GetToolVersions 获取工具版本信息
func GetToolVersions() map[string]string {
	versions := make(map[string]string)

	cwd, err := os.Getwd()
	if err != nil {
		return versions
	}

	toolsDir := filepath.Join(cwd, "tools")
	executor := NewCommandExecutor(toolsDir)

	// 检查JDK并获取版本
	if executor.CheckToolAvailability("JAVA_HOME", "jdk", "java.exe") {
		if version, err := executor.GetJavaVersion(); err == nil {
			versions["JDK"] = version
		} else {
			versions["JDK"] = "已安装 (版本获取失败)"
		}
	} else {
		versions["JDK"] = "未安装"
	}

	// 检查CodeQL并获取版本
	if executor.CheckToolAvailability("CODEQL_HOME", "codeql", "codeql.exe") {
		if version, err := executor.GetCodeQLVersion(); err == nil {
			versions["CodeQL"] = version
		} else {
			versions["CodeQL"] = "已安装 (版本获取失败)"
		}
	} else {
		versions["CodeQL"] = "未安装"
	}

	// 检查Ant并获取版本
	if executor.CheckToolAvailability("ANT_HOME", "ant", "ant.bat") {
		if version, err := executor.GetAntVersion(); err == nil {
			versions["Ant"] = version
		} else {
			versions["Ant"] = "已安装 (版本获取失败)"
		}
	} else {
		versions["Ant"] = "未安装"
	}

	// 检查Procyon并获取版本
	if version, err := executor.GetProcyonVersion(); err == nil {
		versions["Procyon"] = version
	} else {
		versions["Procyon"] = "未安装"
	}

	// 检查Tomcat并获取版本
	if version, err := executor.GetTomcatVersion(); err == nil {
		versions["Tomcat"] = version
	} else {
		versions["Tomcat"] = "未安装"
	}

	return versions
}

// PrintToolVersions 打印所有工具的版本信息
func PrintToolVersions() {
	fmt.Println("\n=== 工具版本信息 ===")
	versions := GetToolVersions()

	for tool, version := range versions {
		fmt.Printf("%s: %s\n", tool, version)
	}
	fmt.Println("===================")
}

// getHalfSystemMemoryGB 获取系统内存的一半（以GB为单位）
func getHalfSystemMemoryGB() int {
	// 获取系统总内存（字节）
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 使用runtime包获取系统内存信息
	// 注意：runtime.MemStats主要用于Go程序的内存统计，这里我们使用一个更直接的方法
	// 在Windows上，我们可以通过执行系统命令来获取内存信息
	if runtime.GOOS == "windows" {
		return getWindowsMemoryGB() / 2
	}

	// 对于其他系统，使用默认值
	return 16 // 默认16GB
}

// getWindowsMemoryGB 获取Windows系统的总内存（GB）
func getWindowsMemoryGB() int {
	// 使用wmic命令获取总物理内存
	cmd := fmt.Sprintf("wmic computersystem get TotalPhysicalMemory /value")
	exec := NewCommandExecutor(".")
	output, err := exec.ExecuteCommand("cmd", "/c", cmd)
	if err != nil {
		fmt.Printf("获取系统内存失败，使用默认值16GB: %v\n", err)
		return 16
	}

	// 解析输出
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "TotalPhysicalMemory=") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				memoryBytes := strings.TrimSpace(parts[1])
				if memoryInt, err := strconv.ParseInt(memoryBytes, 10, 64); err == nil {
					// 转换为GB（1GB = 1024^3 bytes）
					memoryGB := int(memoryInt / (1024 * 1024 * 1024))
					if memoryGB > 0 {
						return memoryGB
					}
				}
			}
		}
	}

	// 如果解析失败，返回默认值
	fmt.Println("解析系统内存失败，使用默认值32GB")
	return 32
}
