package main

import (
	"fmt"
	"github.com/paketo-buildpacks/libjvm/count"
	"github.com/paketo-buildpacks/libjvm/helper"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	envJavaHome                 = "JAVA_HOME"
	envJavaToolOptions          = "JAVA_TOOL_OPTIONS"
	envBplJvmHeadRoom           = "BPL_JVM_HEAD_ROOM"
	envBplJvmThreadCount        = "BPL_JVM_THREAD_COUNT"
	envBpiApplicationPath       = "BPI_APPLICATION_PATH"
	envBpiJvmClassCount         = "BPI_JVM_CLASS_COUNT"
	envBpiMemoryLimitPathV2     = "BPI_MEMORY_LIMIT_PATH_V2"
	defaultMemoryLimitPathV2Fix = "/sys/fs/cgroup/memory/memory.max_usage_in_bytes"
	defaultJvmOptions           = ""
	defaultHeadRoom             = 0
	defaultThreadCount          = 200
	defaultApplicationPath      = "/app"
	desc                        = `This command calculate the JVM memory for applications to run smoothly and stay within the memory limits of the container.
In order to perform this calculation, the Memory Calculator requires the following input:

  --loaded-class-count: the number of classes that will be loaded when the application is running
  --thread-count: the number of user threads
  --jvm-options: VM Options, typically JAVA_TOOL_OPTIONS
  --head-room: percentage of total memory available which will be left unallocated to cover JVM overhead
  --application-path: the directory on the container where the app's contents are placed
  --output: write to a file, instead of STDOUT

Examples:
  # Use ZGC and output to /tmp/.env 
  memory-calculator --jvm-options '-XX:+UseZGC' -o '/tmp/.env'
`
)

var Version = "<unknown>"

type Config struct {
	jvmOptions        string
	headRoom          int
	threadCount       int
	loadedClassCount  int
	applicationPath   string
	memoryLimitPathV2 string
	output            string
	version           bool
}

func main() {
	c := newConfig()
	cmd := &cobra.Command{
		Use:          "memory-calculator",
		Short:        "JVM Memory Calculator",
		Long:         desc,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if c.version {
				fmt.Println(Version)
				return nil
			}
			return run(c)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&c.jvmOptions, "jvm-options", c.jvmOptions, "vm options, typically JAVA_TOOL_OPTIONS")
	flags.IntVar(&c.headRoom, "head-room", c.headRoom, "percentage of total memory available which will be left unallocated to cover JVM overhead")
	flags.IntVar(&c.threadCount, "thread-count", c.threadCount, "the number of user threads")
	flags.IntVar(&c.loadedClassCount, "loaded-class-count", c.loadedClassCount, "the number of classes that will be loaded when the application is running")
	flags.StringVar(&c.applicationPath, "application-path", c.applicationPath, "the directory on the container where the app's contents are placed")
	flags.StringVarP(&c.output, "output", "o", c.output, "write to a file, instead of STDOUT")
	flags.BoolVar(&c.version, "version", c.version, "print version and exit")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newConfig() Config {
	c := Config{
		jvmOptions:        defaultJvmOptions,
		headRoom:          defaultHeadRoom,
		threadCount:       defaultThreadCount,
		applicationPath:   defaultApplicationPath,
		memoryLimitPathV2: defaultMemoryLimitPathV2Fix,
	}
	c.jvmOptions, _ = os.LookupEnv(envJavaToolOptions)
	if val, ok := os.LookupEnv(envBplJvmHeadRoom); ok {
		c.headRoom, _ = strconv.Atoi(val)
	}
	if val, ok := os.LookupEnv(envBplJvmThreadCount); ok {
		c.threadCount, _ = strconv.Atoi(val)
	}
	c.applicationPath, _ = os.LookupEnv(envBpiApplicationPath)
	if val, ok := os.LookupEnv(envBpiJvmClassCount); ok {
		c.loadedClassCount, _ = strconv.Atoi(val)
	}
	// 修正部分記憶體限制檔案位置不一致問題
	c.memoryLimitPathV2, _ = os.LookupEnv(envBpiMemoryLimitPathV2)
	return c
}

func (c *Config) prepareLibJvmEnv() (err error) {
	if err = os.Setenv(envBplJvmThreadCount, c.jvmOptions); err != nil {
		return err
	}
	if err = os.Setenv(envBplJvmHeadRoom, strconv.Itoa(c.headRoom)); err != nil {
		return err
	}
	if err = os.Setenv(envBplJvmThreadCount, strconv.Itoa(c.threadCount)); err != nil {
		return err
	}
	if err = os.Setenv(envBpiApplicationPath, c.applicationPath); err != nil {
		return err
	}
	// 計算JVM本身的Class數量
	if c.loadedClassCount == 0 {
		if javaHome, ok := os.LookupEnv(envJavaHome); !ok {
			return fmt.Errorf("failed to lookup %v env", envJavaHome)
		} else {
			jvmClassCount, err := count.Classes(javaHome)
			if err != nil {
				return err
			}
			c.loadedClassCount = jvmClassCount
		}
	}
	if err = os.Setenv(envBpiJvmClassCount, strconv.Itoa(c.loadedClassCount)); err != nil {
		return err
	}
	return nil
}

func run(c Config) (err error) {
	if err = c.prepareLibJvmEnv(); err != nil {
		return err
	}
	var (
		l = bard.NewLogger(os.Stdout)
		a = helper.ActiveProcessorCount{Logger: l}
		j = helper.JavaOpts{Logger: l}
		m = helper.MemoryCalculator{
			Logger:            l,
			MemoryLimitPathV1: helper.DefaultMemoryLimitPathV1,
			MemoryLimitPathV2: c.memoryLimitPathV2,
			MemoryInfoPath:    helper.DefaultMemoryInfoPath,
		}
	)

	cmds := map[string]sherpa.ExecD{
		"active-processor-count": a,
		"java-opts":              j,
		"memory-calculator":      m,
	}

	// 依序執行 helper
	for _, cmd := range cmds {
		values, err := cmd.Execute()
		if err != nil {
			return err
		}
		for k, v := range values {
			v = strings.TrimSpace(v)
			if err = os.Setenv(k, v); err != nil { // update golang environment variable
				return err
			}
		}
	}

	var javaToolOptions = os.Getenv(envJavaToolOptions)

	if c.output == "" {
		log.Printf("%v: %v\n", envJavaToolOptions, javaToolOptions)
		return nil
	}

	file, err := os.Create(c.output)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("export %v='%s'\n", envJavaToolOptions, javaToolOptions))
	return err
}
