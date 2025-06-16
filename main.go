package main

import (
	"codeql_n1ght/Common"
	"codeql_n1ght/Database"
	"codeql_n1ght/Install"
)

func main() {
	Common.Start()

	// 解析命令行参数
	Common.InitFlag()

	// 根据参数决定是否安装工具
	if Common.IsInstall {
		// 安装必要的工具
		Install.InstallAllTools()
		// 设置环境变量
		Common.SetupEnvironment()
		// 显示工具版本信息
		Common.PrintToolVersions()
	}
	if Common.CreateJar != "" {
		Database.Init(Common.CreateJar)
	}
}
