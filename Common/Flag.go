package Common

import (
	"flag"
	"fmt"
	"os"
)

func InitFlag() {
	// 主要功能参数
	flag.BoolVar(&IsInstall, "install", false, "一键安装环境")
	flag.StringVar(&CreateJar, "database", "", "通过jar包一键生成数据库")
	flag.BoolVar(&ScanMode, "scan", false, "启用扫描模式")

	// 安装模式专用参数（只能与-install一起使用）
	flag.StringVar(&JDKDownloadURL, "jdk", "", "指定JDK下载地址（仅限-install模式）")
	flag.StringVar(&AntDownloadURL, "ant", "", "指定Apache Ant下载地址（仅限-install模式）")
	flag.StringVar(&CodeQLDownloadURL, "codeql", "", "指定CodeQL下载地址（仅限-install模式）")

	// 扫描模式专用参数（只能与-scan一起使用）
	flag.StringVar(&DatabasePath, "db", "./lib", "指定CodeQL数据库路径（仅限-scan模式）")
	flag.StringVar(&QLLibsPath, "ql", "./qlLibs", "指定QL查询库路径（仅限-scan模式）")
	flag.BoolVar(&CleanCache, "clean-cache", false, "扫描前清理缓存，确保修改的QL文件生效（仅限-scan模式）")

	// 保持向后兼容
	flag.StringVar(&ScanDirectory, "d", "", "【已弃用】指定要扫描的目录，请使用-db和-ql参数")

	// 数据库模式专用参数（只能与-database一起使用）
	flag.StringVar(&ExtraSourceDir, "dir", "", "指定额外的源码目录，将复制到src1中一起生成数据库（仅限-database模式）")
	// 新增：控制依赖选择模式（none=空依赖, all=全依赖；不指定则进入交互选择）
	flag.StringVar(&DependencySelection, "deps", "", "数据库生成时依赖选择：none=空依赖, all=全依赖；不指定进入交互选择（仅限-database模式）")

	// 通用配置参数
	flag.StringVar(&DecompilerType, "decompiler", "procyon", "选择反编译器类型 (procyon|fernflower)")
	flag.BoolVar(&UseGoroutine, "goroutine", false, "启用goroutine并发处理")
	flag.IntVar(&MaxGoroutines, "max-goroutines", 4, "最大goroutine数量（需要-goroutine）")
	flag.BoolVar(&KeepTempFiles, "keep-temp", false, "保留临时文件和目录")
	flag.IntVar(&CodeQLThreads, "threads", 20, "CodeQL处理时的线程数")

	// 自定义help信息
	flag.Usage = printUsage

	flag.Parse()
}

// printUsage 自定义使用说明
func printUsage() {
	fmt.Println("Usage: codeql_n1ght [options]")
	fmt.Println("\n主要功能：")
	fmt.Println("  -install                   一键安装环境")
	fmt.Println("  -database <jar>            通过jar包一键生成数据库")
	fmt.Println("  -scan                      启用扫描模式")

	fmt.Println("\n数据库模式参数（仅与 -database 一起使用）：")
	fmt.Println("  -dir <path>                指定额外源码目录，复制到src1中一起生成数据库")
	fmt.Println("  -deps <none|all>           依赖选择：none=空依赖, all=全依赖；不指定进入交互选择")

	fmt.Println("\n扫描模式参数（仅与 -scan 一起使用）：")
	fmt.Println("  -db <path>                 指定CodeQL数据库路径")
	fmt.Println("  -ql <path>                 指定QL查询库路径")
	fmt.Println("  -clean-cache               扫描前清理缓存，确保修改的QL文件生效")

	fmt.Println("\n安装模式参数（仅与 -install 一起使用）：")
	fmt.Println("  -jdk <url>                 指定JDK下载地址")
	fmt.Println("  -ant <url>                 指定Apache Ant下载地址")
	fmt.Println("  -codeql <url>              指定CodeQL下载地址")

	fmt.Println("\n通用配置：")
	fmt.Println("  -decompiler <type>         选择反编译器类型 (procyon|fernflower)")
	fmt.Println("  -goroutine                 启用goroutine并发处理")
	fmt.Println("  -max-goroutines <n>        最大goroutine数量（需要-goroutine）")
	fmt.Println("  -keep-temp                 保留临时文件和目录")
	fmt.Println("  -threads <n>               CodeQL处理时的线程数")

	fmt.Println("\n示例：")
	fmt.Println("  codeql_n1ght -database app.jar -deps none")
	fmt.Println("  codeql_n1ght -database app.jar -deps all")
	os.Exit(0)
}
