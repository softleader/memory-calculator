//go:build linux && amd64

package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_LinuxResolverContract(t *testing.T) {
	root := newRunFixture(t)
	runContractChild(t, root)

	want := strings.Join([]string{
		runInitialOptions,
		"-Djava.security.properties=" + filepath.Join(root, "java-security.properties"),
		"-XX:MaxDirectMemorySize=10M",
		"-Xmx10205610K",
		"-XX:MaxMetaspaceSize=13909K",
		"-XX:ReservedCodeCacheSize=240M",
		"-Xss1M",
	}, " ")
	if got := readExportedJavaToolOptions(t, root); got != want {
		t.Fatalf("JAVA_TOOL_OPTIONS = %q, want %q", got, want)
	}
}
