package calc

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/paketo-buildpacks/libjvm"
	"github.com/paketo-buildpacks/libjvm/helper"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
	"os"
	"strings"
)

const (
	envBpiJvmCaCerts                  = "BPI_JVM_CACERTS"
	helperActiveProcessorCount        = "active-processor-count"
	helperJavaOpts                    = "java-opts"
	helperJvmHeap                     = "jvm-heap"
	helperLinkLocalDns                = "link-local-dns"
	helperMemoryCalculator            = "memory-calculator"
	helperOpensslCertificateLoader    = "openssl-certificate-loader"
	helperSecurityProvidersClasspath9 = "security-providers-classpath-9"
	helperSecurityProvidersConfigurer = "security-providers-configurer"
	helperDebug9                      = "debug-9"
	helperJmx                         = "jmx"
	helperNmt                         = "nmt"
	helperJfr                         = "jfr"
)

type Calculator struct {
	JVMOptions       *JVMOptions
	HeadRoom         *HeadRoom
	ThreadCount      *ThreadCount
	LoadedClassCount *LoadedClassCount
	AppPath          *AppPath
	EnableNmt        *EnableNmt
	EnableJfr        *EnableJfr
	EnableJmx        *EnableJmx
	EnableJdwp       *EnableJdwp
	Verbose          *Verbose
}

func NewCalculator() Calculator {
	c := Calculator{
		JVMOptions:       NewJVMOptions(),
		HeadRoom:         NewHeadRoom(),
		ThreadCount:      NewThreadCount(),
		LoadedClassCount: NewLoadedClassCount(),
		AppPath:          NewAppPath(),
		EnableNmt:        NewEnableNmt(),
		EnableJfr:        NewEnableJfr(),
		EnableJmx:        NewEnableJmx(),
		EnableJdwp:       NewEnableJdwp(),
		Verbose:          NewVerbose(),
	}
	return c
}

func (c *Calculator) Execute() (*JavaToolOptions, error) {
	if err := contribute(
		c.JVMOptions,
		c.HeadRoom,
		c.ThreadCount,
		c.LoadedClassCount,
		c.AppPath,
		c.EnableNmt,
		c.EnableJfr,
		c.EnableJmx,
		c.EnableJdwp,
		c.Verbose,
	); err != nil {
		return nil, err
	}

	helpers, err := c.buildHelpers()
	if err != nil {
		return nil, err
	}

	helpersInOrder := []string{
		helperActiveProcessorCount,
		helperJavaOpts,
		helperJvmHeap,
		helperLinkLocalDns,
		helperMemoryCalculator,
		helperOpensslCertificateLoader,
		helperSecurityProvidersClasspath9,
		helperSecurityProvidersConfigurer,
		helperDebug9,
		helperJmx,
		helperNmt,
		helperJfr,
	}

	// 按照指定順序執行命令
	for _, key := range helpersInOrder {
		h, ok := helpers[key]
		if !ok {
			continue
		}
		values, err := h.Execute()
		if err != nil {
			return nil, err
		}
		for k, v := range values {
			v = strings.TrimSpace(v)
			if err = os.Setenv(k, v); err != nil { // update golang environment variable
				return nil, err
			}
		}
	}

	return BuildJavaToolOptions(), nil
}

// 這邊基本上是從底層 libjvm 套件中複製過來, 我們只支援 Java 9+ 的計算
// https://github.com/paketo-buildpacks/libjvm/blob/main/cmd/helper/main.go
// https://github.com/paketo-buildpacks/libjvm/blob/main/build.go#L274
func (c *Calculator) buildHelpers() (helpers map[string]sherpa.ExecD, err error) {
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

	helpers = map[string]sherpa.ExecD{
		helperActiveProcessorCount:        a,
		helperJavaOpts:                    j,
		helperJvmHeap:                     jh,
		helperLinkLocalDns:                d,
		helperMemoryCalculator:            m,
		helperOpensslCertificateLoader:    o,
		helperSecurityProvidersClasspath9: s9,
		helperSecurityProvidersConfigurer: spc,
		helperDebug9:                      d9,
		helperJmx:                         jm,
		helperNmt:                         n,
		helperJfr:                         jf,
	}

	// 底層的實作中要求若開啟 jvm-cacert 則必須要設定相關的系統參數, 否則會報錯, 所以針對這個改成沒設定就不要跑了
	if _, ok := os.LookupEnv(envBpiJvmCaCerts); !ok {
		delete(helpers, helperOpensslCertificateLoader)
	}
	// 由於關閉 nmt 底層會印出一些關閉的 log, 我不想要看到那些, 所以針對這個改成沒開啟就不要跑了
	if !*c.EnableNmt {
		delete(helpers, helperNmt)
	}
	return helpers, nil
}

func contribute(cs ...Contributor) error {
	for _, c := range cs {
		if err := c.Contribute(); err != nil {
			return err
		}
	}
	return nil
}
