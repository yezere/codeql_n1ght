package Common

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// CommandExecutor 命令执行器结构体
type CommandExecutor struct {
	ToolsPath string // tools目录路径
}

// NewCommandExecutor 创建新的命令执行器
func NewCommandExecutor(toolsPath string) *CommandExecutor {
	return &CommandExecutor{
		ToolsPath: toolsPath,
	}
}

// GetExecutablePath 从环境变量或tools目录获取可执行文件路径
func (ce *CommandExecutor) GetExecutablePath(envVar, toolSubPath, executableName string) (string, error) {
	// 首先尝试从环境变量获取
	if envPath := os.Getenv(envVar); envPath != "" {
		execPath := filepath.Join(envPath, "bin", executableName)
		if FileExists(execPath) {
			return execPath, nil
		}
		// 如果bin目录下没有，直接在环境变量路径下查找
		execPath = filepath.Join(envPath, executableName)
		if FileExists(execPath) {
			return execPath, nil
		}
	}

	// 如果环境变量不存在或找不到可执行文件，从tools目录查找
	toolPath := filepath.Join(ce.ToolsPath, toolSubPath, "bin", executableName)
	if FileExists(toolPath) {
		return toolPath, nil
	}

	// 直接在工具子目录下查找
	toolPath = filepath.Join(ce.ToolsPath, toolSubPath, executableName)
	if FileExists(toolPath) {
		return toolPath, nil
	}

	return "", fmt.Errorf("无法找到 %s，请检查环境变量 %s 或确保工具已安装在 %s", executableName, envVar, toolSubPath)
}

// ExecuteCommand 执行命令并返回结果
func (ce *CommandExecutor) ExecuteCommand(executablePath string, args ...string) (string, error) {
	cmd := exec.Command(executablePath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("命令执行失败: %v\n输出: %s", err, string(output))
	}
	return string(output), nil
}

// getExecutableName 根据操作系统返回正确的可执行文件名
func getExecutableName(baseName string) string {
	if runtime.GOOS == "windows" {
		if baseName == "ant" {
			return "ant.bat"
		}
		if !strings.HasSuffix(baseName, ".exe") {
			return baseName + ".exe"
		}
		return baseName
	}
	return baseName
}

// getScriptName 根据操作系统返回正确的脚本文件名
func getScriptName(baseName string) string {
	if runtime.GOOS == "windows" {
		return baseName + ".bat"
	}
	return baseName + ".sh"
}

// ExecuteJavaCommand 执行Java命令
func (ce *CommandExecutor) ExecuteJavaCommand(args ...string) (string, error) {
	javaPath, err := ce.GetExecutablePath("JAVA_HOME", "jdk", getExecutableName("java"))
	if err != nil {
		return "", err
	}
	return ce.ExecuteCommand(javaPath, args...)
}

// ExecuteCodeQLCommand 执行CodeQL命令
func (ce *CommandExecutor) ExecuteCodeQLCommand(args ...string) (string, error) {
	codeqlPath, err := ce.GetExecutablePath("CODEQL_HOME", "codeql", getExecutableName("codeql"))
	if err != nil {
		return "", err
	}
	return ce.ExecuteCommand(codeqlPath, args...)
}

// ExecuteAntCommand 执行Ant命令
func (ce *CommandExecutor) ExecuteAntCommand(args ...string) (string, error) {
	antPath, err := ce.GetExecutablePath("ANT_HOME", "ant", getExecutableName("ant"))
	if err != nil {
		return "", err
	}
	return ce.ExecuteCommand(antPath, args...)
}

// GetJavaVersion 获取Java版本信息
func (ce *CommandExecutor) GetJavaVersion() (string, error) {
	output, err := ce.ExecuteJavaCommand("-version")
	if err != nil {
		return "", err
	}
	// Java版本信息通常在stderr中，但CombinedOutput会合并stdout和stderr
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "version") {
			return strings.TrimSpace(line), nil
		}
	}
	return output, nil
}

// GetCodeQLVersion 获取CodeQL版本信息
func (ce *CommandExecutor) GetCodeQLVersion() (string, error) {
	output, err := ce.ExecuteCodeQLCommand("version")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// GetAntVersion 获取Ant版本信息
func (ce *CommandExecutor) GetAntVersion() (string, error) {
	output, err := ce.ExecuteAntCommand("-version")
	if err != nil {
		return "", err
	}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Apache Ant") {
			return strings.TrimSpace(line), nil
		}
	}
	return output, nil
}
func (ce *CommandExecutor) GetProcyonVersion() (string, error) {
	output, err := ce.ExecuteJavaCommand("-jar", "./tools/procyon-decompiler-0.6.0.jar", "--version")
	if err != nil {
		return "", err
	}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		return strings.TrimSpace(line), nil
	}
	return output, nil
}

// GetTomcatVersion 获取Tomcat版本信息
func (ce *CommandExecutor) GetTomcatVersion() (string, error) {
	// 检查tools目录下的tomcat
	tomcatPath := filepath.Join(ce.ToolsPath, "tomcat")

	// 首先尝试执行version脚本
	versionScript := filepath.Join(tomcatPath, "bin", getScriptName("version"))
	if FileExists(versionScript) {
		output, err := ce.ExecuteCommand(versionScript)
		if err == nil {
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				if strings.Contains(line, "Server version:") {
					return strings.TrimSpace(strings.Split(line, ":")[1]), nil
				}
			}
		}
	}

	// 尝试从RELEASE-NOTES文件读取版本
	releaseNotes := filepath.Join(tomcatPath, "RELEASE-NOTES")
	if data, err := os.ReadFile(releaseNotes); err == nil {
		content := string(data)
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Apache Tomcat Version") {
				parts := strings.Split(line, "Version")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
		}
	}

	// 检查是否是已知版本目录结构
	tomcatVersionPath := filepath.Join(tomcatPath, "apache-tomcat-9.0.27")
	if _, err := os.Stat(tomcatVersionPath); err == nil {
		return "9.0.27", nil
	}

	// 检查tomcat目录是否存在
	if _, err := os.Stat(tomcatPath); err == nil {
		return "已安装 (版本未知)", nil
	}

	return "", fmt.Errorf("无法获取Tomcat版本信息")
}

// CheckToolAvailability 检查工具是否可用
func (ce *CommandExecutor) CheckToolAvailability(envVar, toolSubPath, executableName string) bool {
	_, err := ce.GetExecutablePath(envVar, toolSubPath, executableName)
	return err == nil
}
