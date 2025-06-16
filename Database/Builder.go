package Database

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"codeql_n1ght/Common"
)

// Createdatabase 创建CodeQL数据库
func Createdatabase(location string) {
	Common.SetupEnvironment()
	cmd := exec.Command(
		"codeql",
		"database", "create", "temp",
		"--language=java",
		"--command=ant -f build.xml",
		"--source-root", "./",
		"--overwrite",
		"--ram=51200",
		"--threads=400",
	)
	cmd.Dir = location
	// 获取标准输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("获取 StdoutPipe 失败:", err)
		return
	}
	// 获取标准错误管道（可选）
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("获取 StderrPipe 失败:", err)
		return
	}
	// 启动命令
	if err := cmd.Start(); err != nil {
		fmt.Println("启动命令失败:", err)
		return
	}
	// 创建协程并发读取标准输出
	go streamOutput(stdout, "STDOUT")
	// 创建协程并发读取标准错误
	go streamOutput(stderr, "STDERR")
	// 等待命令结束
	if err := cmd.Wait(); err != nil {
		fmt.Println("命令执行异常:", err)
	}
}

// GenerateBuildXML 生成Ant构建文件
func GenerateBuildXML(location string) error {
	buildxml := fmt.Sprintf(`
<project name="fax" basedir="." default="build">
  <property name="src.dir" value="src1"/>
  <property name="web.dir" value="default"/>
  <property name="build.dir" value="build_classes"/>
  <property name="tomcat.dir" value="%s"/>
  <path id="master-classpath">
    <pathelement path="${tomcat.dir}/lib"/>
    <fileset dir="${tomcat.dir}/lib">
      <include name="*.jar"/>
    </fileset>
    <fileset dir="${tomcat.dir}/bin">
      <include name="*.jar"/>
    </fileset>
  </path>
  <target name="build" description="Compile source tree java files">
    <mkdir dir="${build.dir}"/>
    <javac destdir="${build.dir}" source="8" target="8" fork="true" optimize="off" debug="on" failonerror="false">
      <src path="${src.dir}"/>
      <classpath refid="master-classpath"/>
    </javac>
  </target>
</project>
`, os.Getenv("CATALINA_HOME"))

	f, err := os.OpenFile(filepath.Join(location, "build.xml"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("create build.xml failed: %v", err)
	}
	defer f.Close()
	
	_, err = f.Write([]byte(buildxml))
	return err
}