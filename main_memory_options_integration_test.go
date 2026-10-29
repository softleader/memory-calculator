//go:build integration

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_UserConfiguredMemoryOptions(t *testing.T) {
	t.Parallel()

	root := newRunFixture(t)
	if err := os.WriteFile(filepath.Join(root, "memory.limit"), []byte("2g"), 0o644); err != nil {
		t.Fatal(err)
	}
	const initialOptions = "-XX:MaxDirectMemorySize=64M -XX:MaxMetaspaceSize=512m"
	runContractChildWithOptions(t, root, initialOptions)
	options := readExportedJavaToolOptions(t, root)
	if !strings.HasPrefix(options, initialOptions+" ") {
		t.Errorf("JAVA_TOOL_OPTIONS = %q, want original options first and in order", options)
	}

	for _, tt := range []struct {
		name   string
		prefix string
		want   string
	}{
		{
			name:   "direct memory",
			prefix: "-XX:MaxDirectMemorySize=",
			want:   "-XX:MaxDirectMemorySize=64M",
		},
		{
			name:   "metaspace",
			prefix: "-XX:MaxMetaspaceSize=",
			want:   "-XX:MaxMetaspaceSize=512m",
		},
		{
			name:   "heap accounts for both user settings",
			prefix: "-Xmx",
			// 2048 MiB - 64 MiB direct - 512 MiB metaspace - 240 MiB code cache - 10 × 1 MiB stacks.
			want: "-Xmx1222M",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var matches []string
			for option := range strings.FieldsSeq(options) {
				if strings.HasPrefix(option, tt.prefix) {
					matches = append(matches, option)
				}
			}
			if len(matches) != 1 || matches[0] != tt.want {
				t.Errorf("options matching %q = %q, want exactly [%q]", tt.prefix, matches, tt.want)
			}
		})
	}
}
