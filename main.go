package main

import (
	"fmt"
	"github.com/paketo-buildpacks/libjvm/count"
	"github.com/paketo-buildpacks/libjvm/helper"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	DefaultMemoryLimitPathV2Fix = "/sys/fs/cgroup/memory.max_usage_in_bytes"
)

func main() {
	// 準備需要的環境變數
	if _, ok := os.LookupEnv("JAVA_TOOL_OPTIONS"); !ok {
		os.Setenv("JAVA_TOOL_OPTIONS", "")
	}
	if _, ok := os.LookupEnv("BPL_JVM_HEAD_ROOM"); !ok {
		os.Setenv("BPL_JVM_HEAD_ROOM", "0")
	}
	if _, ok := os.LookupEnv("BPL_JVM_THREAD_COUNT"); !ok {
		os.Setenv("BPL_JVM_THREAD_COUNT", "250")
	}
	if _, ok := os.LookupEnv("BPI_APPLICATION_PATH"); !ok {
		os.Setenv("BPI_APPLICATION_PATH", "/app")
	}
	// 計算JVM本身的Class數量
	if JavaHome, ok := os.LookupEnv("JAVA_HOME"); ok {
		jvmClassCount, err := count.Classes(JavaHome)
		if err != nil {
			log.Fatal(err)
		} else {
			os.Setenv("BPI_JVM_CLASS_COUNT", strconv.Itoa(jvmClassCount))
		}
	}

	// 修正部分記憶體限制檔案位置不一致問題
	var memoryLimitPathV2 string
	memoryLimitPathV2, ok := os.LookupEnv("BPI_MEMORY_LIMIT_PATH_V2")
	if !ok {
		memoryLimitPathV2 = DefaultMemoryLimitPathV2Fix
	}

	var (
		l = bard.NewLogger(os.Stdout)

		a = helper.ActiveProcessorCount{Logger: l}
		j = helper.JavaOpts{Logger: l}
		m = helper.MemoryCalculator{
			Logger:            l,
			MemoryLimitPathV1: helper.DefaultMemoryLimitPathV1,
			MemoryLimitPathV2: memoryLimitPathV2,
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
		if values, err := cmd.Execute(); err == nil {
			for k, v := range values {
				v = strings.TrimSpace(v)
				os.Setenv(k, v) // update golang environment variable
			}
		}
	}

	// 因為 Golang 無法直接對系統環境變數修改，所以需要輸出檔案
	file, err := os.Create("/tmp/.env")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var javaToolOptions = os.Getenv("JAVA_TOOL_OPTIONS")
	if _, err := file.WriteString(fmt.Sprintf("export JAVA_TOOL_OPTIONS='%s'\n", javaToolOptions)); err != nil {
		log.Fatal(err)
	}

	log.Printf("JAVA_TOOL_OPTIONS: %v", javaToolOptions)
}
