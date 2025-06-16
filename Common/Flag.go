package Common

import (
	"flag"
)

func InitFlag() {
	flag.BoolVar(&IsInstall, "install", false, "一键安装环境")
	flag.StringVar(&JDKDownloadURL, "jdk", "", "指定JDK下载地址（例子：-install -jdk http://xxx）")
	flag.StringVar(&AntDownloadURL, "ant", "", "指定Apache Ant下载地址（例子：-ant -jdk http://xxx）")
	flag.StringVar(&CodeQLDownloadURL, "codeql", "", "指定CodeQL下载地址（例子：-codeql -jdk http://xxx）")
	flag.StringVar(&CreateJar, "database", "", "通过jar包一键生成数据库")
	flag.StringVar(&DecompilerType, "decompiler", "procyon", "选择反编译器类型 (procyon|fernflower)，java版本可能和java-decompile.jar不符合，可以手动替换tools/java-decompile.jar")
	flag.Parse()
	// 注意：安装逻辑已移至main.go中处理，避免循环导入
}
