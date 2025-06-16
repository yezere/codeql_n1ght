package Database

import (
	"codeql_n1ght/Common"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// Init 初始化数据库创建流程
func Init(jar string) {
	pwd, _ := os.Getwd()
	jar = filepath.Join(pwd, jar)
	if !Common.FileExists(jar) {
		color.Red("Jar file not found")
		return
	}
	location := filepath.Dir(jar)
	color.Green("Jar file found")

	// 清理旧文件
	cleanupOldFiles(location)

	// 解压jar包
	Common.ExtractZip(jar, filepath.Join(location, "output"))
	Common.SetupEnvironment()
	color.Green("解压完成")

	// 设置创建数据库目录
	setupDatabaseDirectory(location)

	// 生成构建文件
	err := GenerateBuildXML(filepath.Join(location, "createdabase"))
	if err != nil {
		color.Red("Generate build.xml failed: %v", err)
		return
	}

	// 创建必要的目录
	createDirectories(location)

	// 检查是否为war包，如果是则特殊处理
	if filepath.Ext(jar) == ".war" {
		// 对于war包，直接反编译classes目录和JSP文件
		outputDir := filepath.Join(location, "output")
		src1Dir := filepath.Join(location, "createdabase", "src1")

		// 反编译BOOT-INF/classes目录
		classesDir := filepath.Join(outputDir, "BOOT-INF", "classes")
		if _, err := os.Stat(classesDir); err == nil {
			color.Green("开始反编译BOOT-INF/classes目录")
			err := DecompileJava("-cp", "tools/java-decompiler.jar",
				"org.jetbrains.java.decompiler.main.decompiler.ConsoleDecompiler",
				"-dgs=true", "-hdc=0", "-dgs=1", "-rsy=1", "-rbr=1", "-lit=1", "-nls=1", "-mpm=60",
				classesDir, src1Dir)
			if err != nil {
				color.Red("BOOT-INF/classes目录反编译失败: %v", err)
				return
			}
			color.Green("BOOT-INF/classes目录反编译完成")
		}

		// 检查传统WAR包的WEB-INF/classes目录
		webInfClassesDir := filepath.Join(outputDir, "WEB-INF", "classes")
		if _, err := os.Stat(webInfClassesDir); err == nil {
			color.Green("开始反编译WEB-INF/classes目录")
			err := DecompileJava("-cp", "tools/java-decompiler.jar",
				"org.jetbrains.java.decompiler.main.decompiler.ConsoleDecompiler",
				"-dgs=true", "-hdc=0", "-dgs=1", "-rsy=1", "-rbr=1", "-lit=1", "-nls=1", "-mpm=60",
				webInfClassesDir, src1Dir)
			if err != nil {
				color.Red("WEB-INF/classes目录反编译失败: %v", err)
				return
			}
			color.Green("WEB-INF/classes目录反编译完成")
		}

		// 反编译JSP文件
		color.Green("反编译JSP文件: ")
		err := DecompileJava("-jar", "tools/jsp2class.jar", outputDir, src1Dir)
		if err != nil {
			color.Red("JSP文件反编译失败 %s: %v", err)
			// JSP反编译失败不影响整体流程，继续执行
		}
		color.Green("反编译JSP文件: 完成")
	} else {
		// 对于普通jar包，使用原有逻辑
		DecompileJava("-jar", "tools/procyon-decompiler-0.6.0.jar", jar, "-o", filepath.Join(location, "createdabase", "src1"))
	}

	// 反编译依赖到src1
	DecompileLibraries(location)

	// 清理可能导致编译失败的文件
	cleanupProblematicFiles(location)

	// 创建数据库
	Createdatabase(filepath.Join(location, "createdabase"))

	// 移动和清理文件
	finalizeDatabaseCreation(location)
}

// cleanupOldFiles 清理旧文件
func cleanupOldFiles(location string) {
	color.Red("删除output")
	color.Red("删除src")
	Common.RemoveFile(filepath.Join(location, "output"))
	Common.RemoveFile(filepath.Join(location, "src"))
}

// setupDatabaseDirectory 设置数据库目录
func setupDatabaseDirectory(location string) {
	if Common.FileExists(filepath.Join(location, "createdabase")) {
		color.Red("已删除存在的createdatabase")
		os.Remove(filepath.Join(location, "createdabase"))
	}
	os.Mkdir(filepath.Join(location, "createdabase"), 0755)
}

// createDirectories 创建必要的目录
func createDirectories(location string) {
	os.Mkdir(filepath.Join(location, "createdabase", "src1"), 0755)
	os.Mkdir(filepath.Join(location, "createdabase", "src2"), 0755)
	os.Mkdir(filepath.Join(location, "createdabase", "build_classes"), 0755)
}

// cleanupProblematicFiles 清理可能导致编译失败的文件
func cleanupProblematicFiles(location string) {
	src1Dir := filepath.Join(location, "createdabase", "src1")
	
	// 删除所有.kt文件
	ktFiles, err := filepath.Glob(filepath.Join(src1Dir, "**", "*.kt"))
	if err == nil {
		for _, ktFile := range ktFiles {
			os.Remove(ktFile)
			color.Yellow("删除Kotlin文件: %s", filepath.Base(ktFile))
		}
	}
	
	// 递归查找并删除所有.kt文件
	filepath.Walk(src1Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".kt" {
			os.Remove(path)
			color.Yellow("删除Kotlin文件: %s", path)
		}
		return nil
	})
	
	// 删除所有module-info.java文件
	filepath.Walk(src1Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && info.Name() == "module-info.java" {
			os.Remove(path)
			color.Yellow("删除模块信息文件: %s", path)
		}
		return nil
	})
	
	color.Green("清理问题文件完成")
}

// finalizeDatabaseCreation 完成数据库创建的最后步骤
func finalizeDatabaseCreation(location string) {
	Common.RemoveFile(filepath.Join(location, "temp"))

	err := os.Rename(filepath.Join(location, "createdabase", "temp"), filepath.Join(location, "temp"))
	if err != nil {
		color.Red("移动失败:", err)
	} else {
		color.Green("数据库移动成功")
	}
	color.Green("数据库生成完成，删除output和createdatabase")
	Common.RemoveFile(filepath.Join(location, "createdabase"))
	Common.RemoveFile(filepath.Join(location, "output"))
	color.Green("删除成功")
}
