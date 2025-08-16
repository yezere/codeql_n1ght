# CodeQL N1ght

一个自用的 CodeQL 数据库自动化创建工具，支持 JAR/WAR 包的反编译和数据库生成。
没用使用maven去创建codeql数据库，使用了ant去编译这样，就不会因为报错而终止数据库生成。
可以手动在exe所在目录下，修改tools里面的工具。（~~其中基本所有代码都是trea ai写的，包括readme.md，为了不浪费我首月花费的3美元才写的~~）
~~可以加入我的qq群玩一玩：1027627836~~

## 🚀 功能特性

- **一键环境安装**：自动下载并配置 JDK、Apache Ant、CodeQL 等必要工具
- **智能反编译**：支持 JAR 和 WAR 包的自动反编译
- **多反编译器支持**：支持 Procyon 和 Fernflower 反编译器
- **WAR 包特殊处理**：针对 Spring Boot 和传统 WAR 包的智能路径处理
- **自动数据库创建**：一键生成 CodeQL 数据库用于安全分析
- **安全扫描功能**：集成 CodeQL 扫描引擎，支持并发扫描和报告生成
- **多格式报告**：生成 SARIF 和 HTML 格式的扫描报告
- **并发处理**：支持 Goroutine 并发反编译和扫描，提升处理效率

## 📋 系统要求

- Go 1.22.0 或更高版本
- 网络连接（用于下载工具）
- 足够的磁盘空间（建议至少 2GB）

## 🛠️ 安装

### 方法一：直接下载可执行文件

从 [Releases](https://github.com/yezere/codeql_n1ght/releases) 页面下载对应平台的可执行文件。

### 方法二：从源码编译

```bash
git clone https://github.com/yezere/codeql_n1ght.git
cd codeql_n1ght
go build -o codeql_n1ght
```

## 🎯 快速开始

### 1. 一键安装环境

```bash
# 安装所有必要工具（JDK、Apache Ant、CodeQL）
./codeql_n1ght -install

# 使用自定义下载地址安装
./codeql_n1ght -install -jdk https://your-jdk-url.zip -codeql https://your-codeql-url.zip
```

### 2. 创建 CodeQL 数据库

```bash
# 从 JAR 包创建数据库
./codeql_n1ght -database your-app.jar

# 从 WAR 包创建数据库
./codeql_n1ght -database your-webapp.war

# 指定反编译器类型
./codeql_n1ght -database your-app.jar -decompiler fernflower

# 反编译自己想要的lib，将jar包放入lib文件夹下，打包成zip
./codeql_n1ght -database your-zip.zip
```

### 3. 执行安全扫描

```bash
# 扫描数据库（使用默认路径）
./codeql_n1ght -scan

# 扫描数据库（指定路径）
./codeql_n1ght -scan -db ./mydb -ql ./myqueries

# 并发扫描（提升扫描速度）
./codeql_n1ght -scan -db ./mydb -ql ./myqueries -goroutine -max-goroutines 8

# 清理缓存后扫描（确保修改的QL文件生效）
./codeql_n1ght -scan -clean-cache
```

## 📖 详细用法

### 命令行参数

#### 基础功能参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `-install` | 一键安装环境 | `./codeql_n1ght -install` |
| `-database` | 指定要分析的 JAR/WAR 文件 | `./codeql_n1ght -database app.jar` |
| `-scan` | 执行 CodeQL 安全扫描 | `./codeql_n1ght -scan` |
| `-decompiler` | 选择反编译器 (procyon\|fernflower) | `./codeql_n1ght -database app.jar -decompiler fernflower` |

#### 扫描功能参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `-db` | 指定 CodeQL 数据库路径 | `./codeql_n1ght -scan -db ./mydb` |
| `-ql` | 指定 QL 查询文件或目录路径 | `./codeql_n1ght -scan -ql ./myqueries` |
| `-goroutine` | 启用并发扫描模式 | `./codeql_n1ght -scan -goroutine` |
| `-max-goroutines` | 设置最大并发数 | `./codeql_n1ght -scan -goroutine -max-goroutines 8` |
| `-threads` | 设置 CodeQL 扫描线程数 | `./codeql_n1ght -scan -threads 4` |
| `-clean-cache` | 清理 CodeQL 缓存 | `./codeql_n1ght -scan -clean-cache` |

#### 自定义下载参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `-jdk` | 自定义 JDK 下载地址 | `./codeql_n1ght -install -jdk https://example.com/jdk.zip` |
| `-ant` | 自定义 Apache Ant 下载地址 | `./codeql_n1ght -install -ant https://example.com/ant.zip` |
| `-codeql` | 自定义 CodeQL 下载地址 | `./codeql_n1ght -install -codeql https://example.com/codeql.zip` |

### 工作流程

#### 数据库创建流程

1. **环境检查**：检查必要工具是否已安装
2. **文件解压**：解压 JAR/WAR 包到临时目录
3. **智能反编译**：
   - JAR 包：反编译所有 class 文件
   - WAR 包：分别处理 `BOOT-INF/classes`、`WEB-INF/classes` 和 JSP 文件
4. **构建配置**：生成 Apache Ant 构建文件
5. **数据库创建**：使用 CodeQL 创建分析数据库

#### 安全扫描流程

1. **扫描准备**：验证数据库和查询文件路径
2. **清理环境**：清理之前的扫描结果和缓存
3. **源码提取**：从数据库中提取源码文件（如果需要）
4. **查询执行**：
   - 顺序模式：逐个执行 QL 查询文件
   - 并发模式：使用 Goroutine 并发执行查询
5. **结果生成**：生成 SARIF 和 HTML 格式的扫描报告
6. **报告展示**：显示扫描摘要和结果统计

### WAR 包特殊处理

本工具针对 WAR 包进行了特殊优化：

- **Spring Boot JAR/WAR**：自动处理 `BOOT-INF/classes` 和 `BOOT-INF/lib` 目录
- **传统 WAR**：兼容处理 `WEB-INF/classes` 和 `WEB-INF/lib` 目录
- **JSP 文件**：使用专用的 `jsp2class.jar` 进行反编译
- **智能路径检测**：自动识别不同的 WAR 包结构

## 📁 项目结构

```
codeql_n1ght/
├── Common/          # 公共工具模块
│   ├── CommandExecutor.go  # 命令执行器
│   ├── Config.go           # 配置管理
│   ├── Environment.go      # 环境变量设置
│   ├── Flag.go             # 命令行参数解析
│   ├── Start.go            # 启动界面
│   └── Utils.go            # 工具函数
├── Database/        # 数据库创建模块
│   ├── Builder.go          # CodeQL 数据库构建
│   ├── Decompile.go        # 反编译入口
│   ├── Decompiler.go       # 反编译器实现
│   ├── Initializer.go      # 初始化流程
│   └── Utils.go            # 数据库工具函数
├── Install/         # 工具安装模块
│   ├── AntDownload.go      # Apache Ant 下载
│   ├── CodeqlDownload.go   # CodeQL 下载
│   ├── DecompileDownload.go # 反编译器下载
│   ├── JDKDownload.go      # JDK 下载
│   ├── TomcatDownload.go   # Tomcat 下载
│   └── Utils.go            # 安装工具函数
├── Scanner/         # 安全扫描模块
│   ├── Scanner.go          # 扫描引擎核心
│   ├── cleanup.go          # 清理工具
│   ├── file_extractor.go   # 文件提取器
│   ├── hints.go            # 扫描提示
│   └── html_report.go      # HTML 报告生成
├── qlLibs/          # CodeQL 查询库（自动创建）
├── tools/           # 工具目录（自动创建）
│   ├── ant/         # Apache Ant
│   ├── codeql/      # CodeQL CLI
│   └── jdk/         # JDK
├── results.sarif    # SARIF 格式扫描结果
├── scan_report.html # HTML 格式扫描报告
└── main.go          # 主程序入口
```

## 🔧 配置说明

### 反编译器选择

- **Procyon**（默认）：Java 反编译效果较好，推荐用于大多数场景
- **Fernflower**：IntelliJ IDEA 内置反编译器，在某些复杂场景下表现更好

### 自定义工具版本

如果默认的工具版本与你的 Java 版本不兼容，可以手动替换：

```bash
# 替换 Java 反编译器
cp your-java-decompiler.jar tools/java-decompiler.jar

# 替换 JSP 反编译器
cp your-jsp2class.jar tools/jsp2class.jar
```

## 🐛 故障排除

### 常见问题

1. **下载失败**
   - 检查网络连接
   - 尝试使用自定义下载地址
   - 检查防火墙设置

2. **反编译失败**
   - 尝试切换反编译器：`-decompiler fernflower`
   - 检查 JAR/WAR 文件是否损坏
   - 确保有足够的磁盘空间

3. **数据库创建失败**
   - 检查 CodeQL 是否正确安装
   - 确保 Apache Ant 在 PATH 中
   - 检查内存设置（默认 51200MB）
   - 可能存在一些特殊问题

4. **扫描失败**
   - 确保数据库路径正确且数据库完整
   - 检查 QL 查询文件是否存在
   - 尝试运行 `codeql pack install` 安装依赖包
   - 使用 `-clean-cache` 参数清理缓存

5. **并发扫描问题**
   - 降低并发数：`-max-goroutines 4`
   - 检查系统内存是否充足
   - 尝试使用顺序扫描模式（不使用 `-goroutine` 参数）

6. **报告生成失败**
   - 检查当前目录的写入权限
   - 确保没有其他程序占用结果文件
   - 检查磁盘空间是否充足

---

⭐ 如果这个项目对你有帮助，请给它一个 Star！