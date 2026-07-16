package main

import (
	"io"
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

const (
	runChildEnv       = "MEMORY_CALCULATOR_RUN_CHILD"
	runRootEnv        = "MEMORY_CALCULATOR_RUN_ROOT"
	runInitialOptions = "-Dtest.option=present -XX:ActiveProcessorCount=2 -XX:+ExitOnOutOfMemoryError"
)

func TestRun_ExportAndLoadedClassCountPrecedence(t *testing.T) {
	root := newRunFixture(t)
	runContractChild(t, root)
	options := readExportedJavaToolOptions(t, root)

	for _, required := range []string{
		"-Dtest.option=present",
		"-XX:ActiveProcessorCount=2",
		"-Djava.security.properties=" + filepath.Join(root, "java-security.properties"),
		"-XX:+ExitOnOutOfMemoryError",
		"-Xmx",
		"-XX:MaxDirectMemorySize=",
		"-XX:MaxMetaspaceSize=",
		"-XX:ReservedCodeCacheSize=",
		"-Xss",
	} {
		if !strings.Contains(options, required) {
			t.Errorf("JAVA_TOOL_OPTIONS %q does not contain %q", options, required)
		}
	}
	if got := strings.Count(options, "-XX:+ExitOnOutOfMemoryError"); got != 1 {
		t.Errorf("ExitOnOutOfMemoryError count = %d, want 1 in %q", got, options)
	}
}

func TestRun_ChildProcess(t *testing.T) {
	if os.Getenv(runChildEnv) == "" {
		return
	}

	root := os.Getenv(runRootEnv)
	logger := bard.NewLogger(io.Discard)
	optimizer := boot.NewSpringOptimizer(logger)
	if err := optimizer.AppClassesPath.Set(filepath.Join(root, "classes")); err != nil {
		t.Fatal(err)
	}
	if err := optimizer.AppLibPath.Set(filepath.Join(root, "libs")); err != nil {
		t.Fatal(err)
	}

	calculator := calc.NewCalculator(logger)
	calculator.MemoryLimitPath.V1 = filepath.Join(root, "memory.limit")
	calculator.MemoryLimitPath.V2 = filepath.Join(root, "missing-memory-limit")
	if err := calculator.LoadedClassCount.Set("42"); err != nil {
		t.Fatal(err)
	}
	if err := calculator.ThreadCount.Set("10"); err != nil {
		t.Fatal(err)
	}
	if err := calculator.EnableJdwp.Set("false"); err != nil {
		t.Fatal(err)
	}

	c := config{
		output: filepath.Join(root, "memory.env"),
		logger: logger,
		prep: prep.PreparerManager{
			Logger: logger,
			Preparers: []prep.Preparer{
				prep.NewJavaSecurityProperties(logger, root),
				prep.NewJrePreparer(logger),
			},
		},
		boot: optimizer,
		calc: calculator,
	}
	if err := run(c); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv(calc.EnvLoadedClassCount); got != "42" {
		t.Fatalf("%s = %q, want flag value 42", calc.EnvLoadedClassCount, got)
	}
}

func newRunFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, dir := range []string{
		filepath.Join(root, "classes"),
		filepath.Join(root, "libs"),
		filepath.Join(root, "java-home", "conf", "security"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(
		filepath.Join(root, "java-home", "conf", "security", "java.security"),
		[]byte("security.provider.1=SUN\nsecurity.provider.2=SunRsaSign\n"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "memory.limit"), []byte("10g"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func runContractChild(t *testing.T, root string) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestRun_ChildProcess$")
	cmd.Env = []string{
		runChildEnv + "=1",
		runRootEnv + "=" + root,
		"HOME=" + root,
		"JAVA_HOME=" + filepath.Join(root, "java-home"),
		"JAVA_TOOL_OPTIONS=" + runInitialOptions,
		calc.EnvLoadedClassCount + "=41",
	}
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("run child: %v\n%s", err, output)
	}
}

func readExportedJavaToolOptions(t *testing.T, root string) string {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(root, "memory.env"))
	if err != nil {
		t.Fatal(err)
	}
	const prefix = "export JAVA_TOOL_OPTIONS='"
	const suffix = "'\n"
	text := string(contents)
	if !strings.HasPrefix(text, prefix) || !strings.HasSuffix(text, suffix) {
		t.Fatalf("export = %q, want %q...%q", text, prefix, suffix)
	}
	return strings.TrimSuffix(strings.TrimPrefix(text, prefix), suffix)
}
