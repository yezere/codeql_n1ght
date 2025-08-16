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
	fmt.Fprintf(os.Stderr, "CodeQL N1ght - CodeQL数据库自动化创建和扫描工具\n\n")
	fmt.Fprintf(os.Stderr, "使用方法:\n")
	fmt.Fprintf(os.Stderr, "  %s [选项]\n\n", os.Args[0])

	fmt.Fprintf(os.Stderr, "主要功能（必选其一）:\n")
	fmt.Fprintf(os.Stderr, "  -install\n")
	fmt.Fprintf(os.Stderr, "        一键安装环境（JDK、CodeQL、Apache Ant等）\n")
	fmt.Fprintf(os.Stderr, "  -database string\n")
	fmt.Fprintf(os.Stderr, "        通过jar/war包一键生成CodeQL数据库\n")
	fmt.Fprintf(os.Stderr, "  -scan\n")
	fmt.Fprintf(os.Stderr, "        启用扫描模式，对CodeQL数据库执行查询\n\n")

	fmt.Fprintf(os.Stderr, "安装模式专用选项（与-install配合使用）:\n")
	fmt.Fprintf(os.Stderr, "  -jdk string\n")
	fmt.Fprintf(os.Stderr, "        指定JDK下载地址\n")
	fmt.Fprintf(os.Stderr, "  -ant string\n")
	fmt.Fprintf(os.Stderr, "        指定Apache Ant下载地址\n")
	fmt.Fprintf(os.Stderr, "  -codeql string\n")
	fmt.Fprintf(os.Stderr, "        指定CodeQL下载地址\n\n")

	fmt.Fprintf(os.Stderr, "扫描模式专用选项（与-scan配合使用）:\n")
	fmt.Fprintf(os.Stderr, "  -db string\n")
	fmt.Fprintf(os.Stderr, "        指定CodeQL数据库路径 (默认 \"./lib\")\n")
	fmt.Fprintf(os.Stderr, "  -ql string\n")
	fmt.Fprintf(os.Stderr, "        指定QL查询库路径 (默认 \"./qlLibs\")\n")
	fmt.Fprintf(os.Stderr, "  -clean-cache\n")
	fmt.Fprintf(os.Stderr, "        扫描前清理缓存，确保修改的QL文件生效\n")
	fmt.Fprintf(os.Stderr, "  -d string\n")
	fmt.Fprintf(os.Stderr, "        【已弃用】指定要扫描的目录，请使用-db和-ql参数\n\n")

	fmt.Fprintf(os.Stderr, "数据库模式专用选项（与-database配合使用）:\n")
	fmt.Fprintf(os.Stderr, "  -dir string\n")
	fmt.Fprintf(os.Stderr, "        指定额外的源码目录，将复制到src1中一起生成数据库\n\n")

	fmt.Fprintf(os.Stderr, "通用选项:\n")
	fmt.Fprintf(os.Stderr, "  -decompiler string\n")
	fmt.Fprintf(os.Stderr, "        选择反编译器类型 (procyon|fernflower) (默认 \"procyon\")\n")
	fmt.Fprintf(os.Stderr, "  -goroutine\n")
	fmt.Fprintf(os.Stderr, "        启用goroutine并发处理\n")
	fmt.Fprintf(os.Stderr, "  -max-goroutines int\n")
	fmt.Fprintf(os.Stderr, "        最大goroutine数量 (默认 4)\n")
	fmt.Fprintf(os.Stderr, "  -keep-temp\n")
	fmt.Fprintf(os.Stderr, "        保留临时文件和目录\n")
	fmt.Fprintf(os.Stderr, "  -threads int\n")
	fmt.Fprintf(os.Stderr, "        CodeQL处理时的线程数 (默认 20)\n\n")

	fmt.Fprintf(os.Stderr, "使用示例:\n")
	fmt.Fprintf(os.Stderr, "  # 安装环境\n")
	fmt.Fprintf(os.Stderr, "  %s -install\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  # 使用自定义下载地址安装\n")
	fmt.Fprintf(os.Stderr, "  %s -install -jdk https://example.com/jdk.zip\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  # 创建数据库\n")
	fmt.Fprintf(os.Stderr, "  %s -database app.jar\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  # 创建数据库并包含额外源码\n")
	fmt.Fprintf(os.Stderr, "  %s -database app.jar -dir ./extra-sources\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  # 扫描数据库（使用默认路径）\n")
	fmt.Fprintf(os.Stderr, "  %s -scan\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  # 扫描数据库（指定路径）\n")
	fmt.Fprintf(os.Stderr, "  %s -scan -db ./mydb -ql ./myqueries\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  # 并发扫描\n")
	fmt.Fprintf(os.Stderr, "  %s -scan -db ./mydb -ql ./myqueries -goroutine -max-goroutines 8\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  # 清理缓存后扫描（确保修改的QL文件生效）\n")
	fmt.Fprintf(os.Stderr, "  %s -scan -clean-cache\n\n", os.Args[0])
}
