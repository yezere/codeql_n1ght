package Common

var IsInstall bool
var CreateJar string
var DecompilerType string
var UseGoroutine bool
var MaxGoroutines int
var KeepTempFiles bool
var CodeQLThreads int

// 用户指定的下载URL
var JDKDownloadURL string
var AntDownloadURL string
var CodeQLDownloadURL string

// 扫描相关配置
var ScanMode bool
var DatabasePath string
var QLLibsPath string

// 保持向后兼容的旧变量名
var ScanDirectory string

// 额外源码目录配置
var ExtraSourceDir string

// 清理缓存配置
var CleanCache bool
