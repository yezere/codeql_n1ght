package main

import (
	"codeql_n1ght/Common"
	"codeql_n1ght/Database"
	"codeql_n1ght/Install"
	"codeql_n1ght/Scanner"
	"fmt"
	"os"
)

func main() {
	// 显示启动界面
	Common.Start()

	// 解析命令行参数
	Common.InitFlag()

	// 检查参数合法性
	if err := validateArguments(); err != nil {
		Common.LogError("参数验证失败: %v", err)
		os.Exit(1)
	}

	// 执行相应的功能
	if err := executeCommand(); err != nil {
		Common.LogError("执行失败: %v", err)
		os.Exit(1)
	}

	Common.LogInfo("程序执行完成")
}

// validateArguments 验证命令行参数
func validateArguments() error {
	// 检查是否指定了操作
	if !Common.IsInstall && Common.CreateJar == "" && !Common.ScanMode {
		return fmt.Errorf("请指定要执行的操作: -install, -database 或 -scan")
	}

	// 验证下载URL参数只能在install模式下使用
	if !Common.IsInstall {
		if Common.JDKDownloadURL != "" {
			return fmt.Errorf("-jdk 参数只能在 -install 模式下使用")
		}
		if Common.AntDownloadURL != "" {
			return fmt.Errorf("-ant 参数只能在 -install 模式下使用")
		}
		if Common.CodeQLDownloadURL != "" {
			return fmt.Errorf("-codeql 参数只能在 -install 模式下使用")
		}
	}

	// 验证额外源码目录参数只能在database模式下使用
	if Common.ExtraSourceDir != "" && Common.CreateJar == "" {
		return fmt.Errorf("-dir 参数只能在 -database 模式下使用")
	}

	// 验证扫描模式参数
	if Common.ScanMode {
		// 向后兼容处理：如果使用了旧的-d参数，给出提示
		if Common.ScanDirectory != "" {
			Common.LogWarn("-d 参数已弃用，请使用 -db 指定数据库路径，-ql 指定查询库路径")
			// 为了向后兼容，将-d参数的值作为数据库路径
			if Common.DatabasePath == "./lib" { // 如果是默认值
				Common.DatabasePath = Common.ScanDirectory
			}
		}

		// 验证数据库路径
		if !Common.IsDirectory(Common.DatabasePath) {
			return fmt.Errorf("指定的数据库路径不是有效目录: %s", Common.DatabasePath)
		}

		// 验证QL库路径
		if !Common.IsDirectory(Common.QLLibsPath) {
			return fmt.Errorf("指定的QL库路径不是有效目录: %s", Common.QLLibsPath)
		}

		// 扫描模式下不能同时使用install或database
		if Common.IsInstall {
			return fmt.Errorf("扫描模式不能与安装模式同时使用")
		}
		if Common.CreateJar != "" {
			return fmt.Errorf("扫描模式不能与数据库创建模式同时使用")
		}
	}

	// 验证数据库模式参数
	if Common.CreateJar != "" {
		if err := Common.ValidateFile(Common.CreateJar); err != nil {
			return fmt.Errorf("JAR文件验证失败: %v", err)
		}

		// 验证额外源码目录
		if Common.ExtraSourceDir != "" {
			if !Common.IsDirectory(Common.ExtraSourceDir) {
				return fmt.Errorf("指定的额外源码路径不是有效目录: %s", Common.ExtraSourceDir)
			}
		}

		// 数据库模式下不能同时使用install
		if Common.IsInstall {
			return fmt.Errorf("数据库创建模式不能与安装模式同时使用")
		}
	}

	// 验证并发参数
	if Common.UseGoroutine && Common.MaxGoroutines <= 0 {
		return fmt.Errorf("最大goroutine数量必须大于0")
	}

	// 验证线程数参数
	if Common.CodeQLThreads <= 0 {
		return fmt.Errorf("线程数必须大于0")
	}

	return nil
}

// executeCommand 执行相应的命令
func executeCommand() error {
	// 安装工具
	if Common.IsInstall {
		if err := installTools(); err != nil {
			return err
		}
	}

	// 创建数据库
	if Common.CreateJar != "" {
		if err := createDatabase(); err != nil {
			return err
		}
	}

	// 执行扫描
	if Common.ScanMode {
		if err := runScan(); err != nil {
			return err
		}
	}

	return nil
}

// installTools 安装工具
func installTools() error {
	return Common.SafeExecute(func() error {
		Common.LogInfo("开始安装工具...")

		// 安装必要的工具
		if err := Install.InstallAllTools(); err != nil {
			return err
		}

		// 设置环境变量
		if err := Common.SetupEnvironment(); err != nil {
			return err
		}

		// 显示工具版本信息
		Common.PrintToolVersions()

		Common.LogInfo("工具安装完成")
		return nil
	}, "工具安装失败")
}

// createDatabase 创建数据库
func createDatabase() error {
	return Common.SafeExecute(func() error {
		Common.LogInfo("开始创建数据库: %s", Common.CreateJar)
		Database.Init(Common.CreateJar)
		Common.LogInfo("数据库创建完成")
		return nil
	}, "数据库创建失败")
}

// runScan 执行扫描
func runScan() error {
	return Common.SafeExecute(func() error {
		Common.LogInfo("开始扫描 - 数据库: %s, QL库: %s", Common.DatabasePath, Common.QLLibsPath)
		if err := Scanner.RunScan(); err != nil {
			return err
		}
		Common.LogInfo("扫描完成")
		return nil
	}, "扫描执行失败")
}
