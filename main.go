package main

import (
	"fmt"
	"github.com/paketo-buildpacks/libjvm/count"
	"github.com/paketo-buildpacks/libjvm/helper"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

const (
	envJavaHome            = "JAVA_HOME"
	envJavaToolOptions     = "JAVA_TOOL_OPTIONS"
	envBplJvmHeadRoom      = "BPL_JVM_HEAD_ROOM"
	envBplJvmThreadCount   = "BPL_JVM_THREAD_COUNT"
	envBpiApplicationPath  = "BPI_APPLICATION_PATH"
	envBpiJvmClassCount    = "BPI_JVM_CLASS_COUNT"
	defaultJvmOptions      = ""
	defaultHeadRoom        = helper.DefaultHeadroom
	defaultThreadCount     = 200
	defaultApplicationPath = "/app"
	desc                   = `This command calculate the JVM memory for applications to run smoothly and stay within the memory limits of the container.
During the computation process, numerous parameters are required, which must be obtained in a specific order and logic.
The sequence and explanations of these parameters are as follows:

  1. Percentage of reserved space allocated by Memory Calculation tool:
     - First, determine if '--head-room' is passed through args.
     - If not, check the OS environment variable $BPL_JVM_HEAD_ROOM.
     - If neither is available, the default value is 0.

  2. Number of classes loaded at runtime:
     - First, determine if '--loaded-class-count' is passed through args.
     - If not, check the OS environment variable $BPL_JVM_LOADED_CLASS_COUNT.
     - If neither is available, dynamically calculate 35% of the total number of classes in the App directory.

  3. Number of user threads at runtime:
     - First, determine if '--thread-count' is passed through args.
     - If not, check the OS environment variable $BPL_JVM_THREAD_COUNT.
     - If neither is available, the default value is 200.

  4. App directory:
     - First, determine if '--application-path' is passed through args.
     - If not, the default directory is /app.

  5. VM creation parameters:
     - First, determine if '--jvm-options' is passed through args.
     - If not, check the OS environment variable $JAVA_TOOL_OPTIONS.

  6. Java startup parameters:
     - Only check the OS environment variable $JAVA_OPTS.

  7. Java home:
     - Only check the OS environment variable $JAVA_HOME.

Examples:
  # Use ZGC and output to /tmp/.env
  memory-calculator --jvm-options '-XX:+UseZGC' -o '/tmp/.env'

  # Print the version and exit
  memory-calculator --version
`
)

var version = "<unknown>"

type Config struct {
	jvmOptions        string
	headRoom          int
	threadCount       int
	loadedClassCount  int
	applicationPath   string
	memoryLimitPathV2 string
	output            string
	version           bool
	verbose           bool
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
				fmt.Println(version)
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
	flags.BoolVarP(&c.verbose, "verbose", "v", c.verbose, "enable verbose output")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newConfig() Config {
	c := Config{
		jvmOptions:      defaultJvmOptions,
		headRoom:        defaultHeadRoom,
		threadCount:     defaultThreadCount,
		applicationPath: defaultApplicationPath,
	}
	if val, ok := os.LookupEnv(envJavaToolOptions); ok {
		c.jvmOptions = val
	}
	if val, ok := os.LookupEnv(envBplJvmHeadRoom); ok {
		c.headRoom, _ = strconv.Atoi(val)
	}
	if val, ok := os.LookupEnv(envBplJvmThreadCount); ok {
		c.threadCount, _ = strconv.Atoi(val)
	}
	if val, ok := os.LookupEnv(envBpiJvmClassCount); ok {
		c.loadedClassCount, _ = strconv.Atoi(val)
	}
	if val, ok := os.LookupEnv(envBpiApplicationPath); ok {
		c.applicationPath = val
	}
	return c
}

func (c *Config) newLogger() bard.Logger {
	if c.verbose {
		return bard.NewLoggerWithOptions(os.Stdout, bard.WithDebug(os.Stdout))
	}
	return bard.NewLogger(os.Stdout)
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
		l = c.newLogger()
		a = helper.ActiveProcessorCount{Logger: l}
		j = helper.JavaOpts{Logger: l}
		m = helper.MemoryCalculator{
			Logger:            l,
			MemoryLimitPathV1: helper.DefaultMemoryLimitPathV1, // cgroup v1 的記憶體上限路徑
			MemoryLimitPathV2: helper.DefaultMemoryLimitPathV2, // cgroup v2 的記憶體上限路徑
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
		l.Infof("%v: %v\n", envJavaToolOptions, javaToolOptions)
		return nil
	}

	file, err := os.Create(c.output)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			l.Infof("WARNING: failed to close file %v: %v\n", file.Name(), err)
		}
	}(file)
	_, err = file.WriteString(fmt.Sprintf("export %v='%s'\n", envJavaToolOptions, javaToolOptions))
	return err
}
