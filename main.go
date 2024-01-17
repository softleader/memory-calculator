package main

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/paketo-buildpacks/libjvm"
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
	envEnabled                = "MEM_CALC_ENABLED"
	envJavaHome               = "JAVA_HOME"
	envJavaOpts               = "JAVA_OPTS"
	envJavaToolOptions        = "JAVA_TOOL_OPTIONS"
	envBplJvmHeadRoom         = "BPL_JVM_HEAD_ROOM"
	envBplJvmThreadCount      = "BPL_JVM_THREAD_COUNT"
	envBpiApplicationPath     = "BPI_APPLICATION_PATH"
	envBpiJvmLoadedClassCount = "BPL_JVM_LOADED_CLASS_COUNT"
	envBpiJvmCaCerts          = "BPI_JVM_CACERTS"
	envBplJavaNmtEnabled      = "BPL_JAVA_NMT_ENABLED"
	envBplJfrEnabled          = "BPL_JFR_ENABLED"
	envBplJmxEnabled          = "BPL_JMX_ENABLED"
	envBplDebugEnabled        = "BPL_DEBUG_ENABLED"
	envBplDebugPort           = "BPL_DEBUG_PORT"
	envBpLogLevel             = "BP_LOG_LEVEL"
	defaultJvmOptions         = ""
	defaultHeadRoom           = helper.DefaultHeadroom
	defaultThreadCount        = 200
	defaultAppPath            = "/app"
	defaultDebugPort          = 5005
	defaultEnabledNmt         = false
	defaultEnableJfr          = false
	defaultEnableJmx          = false
	defaultEnableJdwp         = true
	desc                      = `This command calculate the JVM memory for applications to run smoothly and stay within the memory limits of the container.
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
     - First, determine if '--app-path' is passed through args.
     - If not, the default directory is /app.

  5. Java startup parameters:
     - First, determine if '--jvm-options' is passed through args.
     - If not, check the OS environment variable $JAVA_OPTS.

  6. Java home:
     - Only check the OS environment variable $JAVA_HOME.

Examples:
  # Use ZGC and output to /tmp/.env
  memory-calculator --jvm-options '-XX:+UseZGC' -o '/tmp/.env'

  # Print the version and exit
  memory-calculator --version
`
)

var (
	version = "<unknown>"
	// 固定要加上的參數, 這些參數可能是 libjvm 在 build image 時加的而非計算出來的, 或是我們自己想要加上去的都可以放
	contributeOptions = []string{"-XX:+ExitOnOutOfMemoryError"}
)

type Config struct {
	enabled           bool // 整個機制是否啟用
	jvmOptions        string
	headRoom          int
	threadCount       int
	loadedClassCount  int
	appPath           string
	memoryLimitPathV2 string
	output            string
	version           bool
	verbose           bool
	enabledNmt        bool
	enableJfr         bool
	enableJmx         bool
	enableJdwp        bool
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
			if !c.enabled {
				fmt.Printf("%v is disabled\n", cmd.Short)
				return nil
			}
			return run(c)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&c.jvmOptions, "jvm-options", c.jvmOptions, "vm options, typically JAVA_OPTS")
	flags.IntVar(&c.headRoom, "head-room", c.headRoom, "percentage of total memory available which will be left unallocated to cover JVM overhead")
	flags.IntVar(&c.threadCount, "thread-count", c.threadCount, "the number of user threads")
	flags.IntVar(&c.loadedClassCount, "loaded-class-count", c.loadedClassCount, "the number of classes that will be loaded when the app is running")
	flags.StringVar(&c.appPath, "app-path", c.appPath, "the directory on the container where the app's contents are placed")
	flags.StringVarP(&c.output, "output", "o", c.output, "write to a file, instead of STDOUT")
	flags.BoolVar(&c.version, "version", c.version, "print version and exit")
	flags.BoolVarP(&c.verbose, "verbose", "v", c.verbose, "enable verbose output")
	flags.BoolVar(&c.enabledNmt, "enable-nmt", c.enabledNmt, "enable Native Memory Tracking (NMT)")
	flags.BoolVar(&c.enableJfr, "enable-jfr", c.enableJfr, "enable Java Flight Recorder (JFR)")
	flags.BoolVar(&c.enableJmx, "enable-jmx", c.enableJmx, "enable Java Management Extensions (JMX)")
	flags.BoolVar(&c.enableJdwp, "enable-jdwp", c.enableJdwp, "enable Java Debug Wire Protocol (JDWP)")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newConfig() Config {
	c := Config{
		enabled:     true,
		jvmOptions:  defaultJvmOptions,
		headRoom:    defaultHeadRoom,
		threadCount: defaultThreadCount,
		appPath:     defaultAppPath,
		enabledNmt:  defaultEnabledNmt,
		enableJfr:   defaultEnableJfr,
		enableJmx:   defaultEnableJmx,
		enableJdwp:  defaultEnableJdwp,
	}
	if val, ok := os.LookupEnv(envEnabled); ok {
		if b, err := strconv.ParseBool(val); err == nil {
			c.enabled = b
		}
	}
	if val, ok := os.LookupEnv(envJavaOpts); ok {
		c.jvmOptions = val
	}
	if val, ok := os.LookupEnv(envBplJvmHeadRoom); ok {
		c.headRoom, _ = strconv.Atoi(val)
	}
	if val, ok := os.LookupEnv(envBplJvmThreadCount); ok {
		c.threadCount, _ = strconv.Atoi(val)
	}
	if val, ok := os.LookupEnv(envBpiJvmLoadedClassCount); ok {
		c.loadedClassCount, _ = strconv.Atoi(val)
	}
	if val, ok := os.LookupEnv(envBpiApplicationPath); ok {
		c.appPath = val
	}
	if val, ok := os.LookupEnv(envBpLogLevel); ok {
		c.verbose = val == "DEBUG"
	}
	if val, ok := os.LookupEnv(envBplJavaNmtEnabled); ok {
		if b, err := strconv.ParseBool(val); err == nil {
			c.enabledNmt = b
		}
	}
	if val, ok := os.LookupEnv(envBplJfrEnabled); ok {
		if b, err := strconv.ParseBool(val); err == nil {
			c.enableJfr = b
		}
	}
	if val, ok := os.LookupEnv(envBplJmxEnabled); ok {
		if b, err := strconv.ParseBool(val); err == nil {
			c.enableJmx = b
		}
	}
	if val, ok := os.LookupEnv(envBplDebugEnabled); ok {
		if b, err := strconv.ParseBool(val); err == nil {
			c.enableJdwp = b
		}
	}
	return c
}

func (c *Config) prepareLibJvmEnv() (err error) {
	if c.verbose && os.Setenv(envBpLogLevel, "DEBUG") != nil {
		return err
	}
	if c.jvmOptions != "" {
		if err = os.Setenv(envJavaOpts, c.jvmOptions); err != nil {
			return err
		}
	}
	if err = os.Setenv(envBplJvmHeadRoom, strconv.Itoa(c.headRoom)); err != nil {
		return err
	}
	if err = os.Setenv(envBplJvmThreadCount, strconv.Itoa(c.threadCount)); err != nil {
		return err
	}
	if err = os.Setenv(envBpiApplicationPath, c.appPath); err != nil {
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
	if err = os.Setenv(envBpiJvmLoadedClassCount, strconv.Itoa(c.loadedClassCount)); err != nil {
		return err
	}
	if err = os.Setenv(envBplJmxEnabled, strconv.FormatBool(c.enableJmx)); err != nil {
		return err
	}
	if err = os.Setenv(envBplJavaNmtEnabled, strconv.FormatBool(c.enabledNmt)); err != nil {
		return err
	}
	if err = os.Setenv(envBplJfrEnabled, strconv.FormatBool(c.enableJfr)); err != nil {
		return err
	}
	if err = os.Setenv(envBplDebugEnabled, strconv.FormatBool(c.enableJdwp)); err != nil {
		return err
	}
	if err = os.Setenv(envBplDebugPort, strconv.Itoa(defaultDebugPort)); err != nil {
		return err
	}
	return nil
}

func run(c Config) (err error) {
	if err = c.prepareLibJvmEnv(); err != nil {
		return err
	}
	cmds, err := c.buildCommands()
	if err != nil {
		return err
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

	javaToolOptions := getJavaToolOptions()
	if c.output == "" {
		fmt.Printf("%v: %v\n", envJavaToolOptions, javaToolOptions)
		return nil
	}
	return writeFile(c.output, javaToolOptions)
}

func getJavaToolOptions() string {
	var javaToolOptions = os.Getenv(envJavaToolOptions)
	for _, option := range contributeOptions {
		if !strings.Contains(javaToolOptions, option) {
			javaToolOptions += " " + option
		}
	}
	return javaToolOptions
}

// 這邊基本上是從底層 libjvm 套件中複製過來, 我們只支援 Java 9+ 的計算
// https://github.com/paketo-buildpacks/libjvm/blob/main/cmd/helper/main.go
// https://github.com/paketo-buildpacks/libjvm/blob/main/build.go#L274
func (c *Config) buildCommands() (cmds map[string]sherpa.ExecD, err error) {
	var (
		l  = bard.NewLogger(os.Stdout)
		cl = libjvm.NewCertificateLoader()

		a   = helper.ActiveProcessorCount{Logger: l}
		spc = helper.SecurityProvidersConfigurer{Logger: l}
		d   = helper.LinkLocalDNS{Logger: l}
		j   = helper.JavaOpts{Logger: l}
		jh  = helper.JVMHeapDump{Logger: l}
		m   = helper.MemoryCalculator{
			Logger:            l,
			MemoryLimitPathV1: helper.DefaultMemoryLimitPathV1, // cgroup v1 的記憶體上限路徑
			MemoryLimitPathV2: helper.DefaultMemoryLimitPathV2, // cgroup v2 的記憶體上限路徑
			MemoryInfoPath:    helper.DefaultMemoryInfoPath,
		}
		o  = helper.OpenSSLCertificateLoader{CertificateLoader: cl, Logger: l}
		s9 = helper.SecurityProvidersClasspath9{Logger: l}
		d9 = helper.Debug9{Logger: l}
		jm = helper.JMX{Logger: l}
		n  = helper.NMT{Logger: l}
		jf = helper.JFR{Logger: l}
	)

	file := "/etc/resolv.conf"
	d.Config, err = dns.ClientConfigFromFile(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read DNS client configuration from %s\n%w", file, err)
	}

	cmds = map[string]sherpa.ExecD{
		"active-processor-count":         a,
		"java-opts":                      j,
		"jvm-heap":                       jh,
		"link-local-dns":                 d,
		"memory-calculator":              m,
		"openssl-certificate-loader":     o,
		"security-providers-classpath-9": s9,
		"security-providers-configurer":  spc,
		"debug-9":                        d9,
		"jmx":                            jm,
		"nmt":                            n,
		"jfr":                            jf,
	}

	// 底層的實作中要求若開啟 jvm-cacert 則必須要設定相關的系統參數, 否則會報錯, 所以針對這個改成沒設定就不要跑了
	if _, ok := os.LookupEnv(envBpiJvmCaCerts); !ok {
		delete(cmds, "openssl-certificate-loader")
	}
	// 由於關閉 nmt 底層會印出一些關閉的 log, 我不想要看到那些, 所以針對這個改成沒開啟就不要跑了
	if !c.enabledNmt {
		delete(cmds, "nmt")
	}
	return cmds, nil
}

func writeFile(file string, content string) (err error) {
	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("WARNING: failed to close file %v: %v\n", file.Name(), err)
		}
	}(out)
	_, err = out.WriteString(fmt.Sprintf("export %v='%s'\n", envJavaToolOptions, content))
	return err
}
