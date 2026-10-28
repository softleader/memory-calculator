//go:build integration

package main

import (
	"errors"
	"fmt"
	springboot "github.com/paketo-buildpacks/spring-boot/v5/boot"
	helper "github.com/softleader/memory-calculator/boot/helper"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/softleader/memory-calculator/boot"
	"github.com/softleader/memory-calculator/calc"
	"github.com/softleader/memory-calculator/prep"
)

func TestRun_ThreadCount(t *testing.T) {
	type scenario struct {
		name, detection               string
		cli, env                      *string
		heap, errorSource, errorValue string
		warning                       bool
	}
	cases := []scenario{
		{name: "real_non_web", heap: "10164650"},
		{name: "reactive", detection: "reactive", heap: "10164650"},
		{name: "servlet", detection: "servlet", heap: "9959850"},
		{name: "fallback", detection: "error", heap: "10011050", warning: true},
		{name: "cli_over_env", detection: "forbidden", cli: new("80"), env: new("120"), heap: "10133930"},
		{name: "env", detection: "forbidden", env: new("120"), heap: "10092970"},
		{name: "explicit_cli_200", detection: "forbidden", cli: new("200"), env: new("120"), heap: "10011050"},
		{name: "explicit_env_200", detection: "forbidden", env: new("200"), heap: "10011050"},
		{name: "cli_over_invalid_env", detection: "forbidden", cli: new("80"), env: new("abc"), heap: "10133930"},
	}
	for _, invalid := range []string{"abc", "0", "-1", "", "999999999999999999999999999999", "1.5"} {
		cases = append(cases, scenario{name: "invalid_cli_" + invalid, detection: "forbidden", cli: new(invalid), env: new("120"), errorSource: "--thread-count", errorValue: invalid})
		cases = append(cases, scenario{name: "invalid_env_" + invalid, detection: "forbidden", env: new(invalid), errorSource: calc.EnvThreadCount, errorValue: invalid})
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := newRunFixture(t)
			cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestRun_ThreadCountChild$")
			cmd.Env = []string{runRootEnv + "=" + root, "HOME=" + root, "JAVA_HOME=" + filepath.Join(root, "java-home"),
				"JAVA_TOOL_OPTIONS=" + runInitialOptions + " -Djava.security.properties=" + filepath.Join(root, "java-security.properties"),
				"TEST_DETECTION=" + tc.detection}
			if tc.cli != nil {
				cmd.Env = append(cmd.Env, "TEST_CLI_THREAD_COUNT="+*tc.cli)
			}
			if tc.env != nil {
				cmd.Env = append(cmd.Env, calc.EnvThreadCount+"="+*tc.env)
			}
			output, err := cmd.CombinedOutput()
			text := string(output)
			if strings.Contains(text, "UNEXPECTED_DETECTION") {
				t.Fatalf("explicit input triggered detection: %s", text)
			}
			if tc.errorSource != "" {
				if err == nil {
					t.Fatalf("invalid input succeeded: %s", text)
				}
				for _, want := range []string{tc.errorSource, fmt.Sprintf("%q", tc.errorValue), "positive integer"} {
					if !strings.Contains(text, want) {
						t.Errorf("error lacks %q: %s", want, text)
					}
				}
				if _, err := os.Stat(filepath.Join(root, "memory.env")); !os.IsNotExist(err) {
					t.Errorf("invalid input created output: %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("child: %v\n%s", err, output)
				}
				// Fixed 10 GiB, 42 classes, direct 10 MiB, code cache 240 MiB,
				// metaspace 13909 KiB, and 1 MiB per stack; literal heaps are independently calculated.
				options := readExportedJavaToolOptions(t, root)
				for _, want := range []string{"-Xmx" + tc.heap + "K", "-Xss1M"} {
					if !strings.Contains(options, want) {
						t.Errorf("missing %s: %s", want, options)
					}
				}
			}
			if tc.warning {
				for _, want := range []string{"Warning", "fixture detection failure", "200"} {
					if !strings.Contains(text, want) {
						t.Errorf("fallback diagnostic lacks %q: %s", want, text)
					}
				}
			} else if strings.Contains(text, "Warning") {
				t.Errorf("unexpected fallback warning: %s", text)
			}
		})
	}
}

func TestRun_ThreadCountChild(t *testing.T) {
	root := os.Getenv(runRootEnv)
	if root == "" {
		return
	}
	switch os.Getenv("TEST_DETECTION") {
	case "forbidden":
		helper.ResolveWebAppType = func() (springboot.ApplicationType, error) {
			fmt.Println("UNEXPECTED_DETECTION")
			return springboot.None, errors.New("unexpected detection")
		}
	case "reactive":
		helper.ResolveWebAppType = func() (springboot.ApplicationType, error) { return springboot.Reactive, nil }
	case "servlet":
		helper.ResolveWebAppType = func() (springboot.ApplicationType, error) { return springboot.Servlet, nil }
	case "error":
		helper.ResolveWebAppType = func() (springboot.ApplicationType, error) {
			return springboot.None, errors.New("fixture detection failure")
		}
	}
	logger := bard.NewLogger(os.Stdout)
	calculator := calc.NewCalculator(logger)
	calculator.MemoryLimitPath.V1 = filepath.Join(root, "memory.limit")
	calculator.MemoryLimitPath.V2 = filepath.Join(root, "missing")
	c := config{logger: logger, prep: prep.NewPreparerManager(logger), boot: boot.NewSpringOptimizer(logger), calc: calculator}
	cmd := newCommand(&c)
	args := []string{"--loaded-class-count=42", "--enable-jdwp=false", "--app-classes-path=" + filepath.Join(root, "classes"), "--app-lib-path=" + filepath.Join(root, "libs"), "--output=" + filepath.Join(root, "memory.env")}
	if value, exists := os.LookupEnv("TEST_CLI_THREAD_COUNT"); exists {
		args = append(args, "--thread-count="+value)
	}
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}
