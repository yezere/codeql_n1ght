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

## 📖 详细用法

### 命令行参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `-install` | 一键安装环境 | `./codeql_n1ght -install` |
| `-database` | 指定要分析的 JAR/WAR 文件 | `./codeql_n1ght -database app.jar` |
| `-decompiler` | 选择反编译器 (procyon\|fernflower) | `./codeql_n1ght -database app.jar -decompiler fernflower` |
| `-jdk` | 自定义 JDK 下载地址 | `./codeql_n1ght -install -jdk https://example.com/jdk.zip` |
| `-ant` | 自定义 Apache Ant 下载地址 | `./codeql_n1ght -install -ant https://example.com/ant.zip` |
| `-codeql` | 自定义 CodeQL 下载地址 | `./codeql_n1ght -install -codeql https://example.com/codeql.zip` |

### 工作流程

1. **环境检查**：检查必要工具是否已安装
2. **文件解压**：解压 JAR/WAR 包到临时目录
3. **智能反编译**：
   - JAR 包：反编译所有 class 文件
   - WAR 包：分别处理 `BOOT-INF/classes`、`WEB-INF/classes` 和 JSP 文件
4. **构建配置**：生成 Apache Ant 构建文件
5. **数据库创建**：使用 CodeQL 创建分析数据库

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
├── tools/           # 工具目录（自动创建）
│   ├── ant/         # Apache Ant
│   ├── codeql/      # CodeQL CLI
│   └── jdk/         # JDK
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

---

⭐ 如果这个项目对你有帮助，请给它一个 Star！