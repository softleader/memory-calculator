//go:build integration

package calc_test

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/softleader/memory-calculator/calc"
)

func TestJavaToolOptions_WriteFileShellRoundTrip(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		value string
	}{
		{name: "ordinary and order", value: `-Xmx512M -Dfirst=one -Dsecond=two -XX:+ExitOnOutOfMemoryError`},
		{name: "apostrophe", value: `-Dlabel="O'Reilly"`},
		{name: "adjacent apostrophes", value: `''-Dlabel='one''two'`},
		{name: "double quotes and spaces", value: `-Dlabel="two words"`},
		{name: "backslashes", value: `-Dpath="C:\folder\file" -Dpattern=\d+`},
		{name: "literal variable", value: `-Dvalue=$VAR`},
		{name: "command substitution", value: `-Dvalue=$(printf expanded)`},
		{name: "backtick substitution", value: "-Dvalue=`printf expanded`"},
		{name: "apostrophe and substitutions", value: `-Dlabel="O'Reilly" -Dvalue=$VAR -Dcommand=$(printf expanded)`},
		{name: "whitespace and unicode", value: "\t-Dlabel=\"繁體 中文\"  -Xmx512M\n"},
		{name: "empty", value: ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(t.TempDir(), "options with spaces.env")
			options := calc.JavaToolOptions(tc.value)
			if err := options.WriteFile(path); err != nil {
				t.Fatal(err)
			}
			// Read from a child shell to verify that the value is exported, too.
			cmd := exec.CommandContext(t.Context(), "/bin/sh", "-eu", "-c",
				`. "$1"; exec /bin/sh -eu -c 'printf "%s" "$JAVA_TOOL_OPTIONS"'`, "read-options", path)
			cmd.Env = []string{"PATH=/usr/bin:/bin", "VAR=expanded"}
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("load exported options: %v\n%s", err, output)
			}
			if string(output) != tc.value {
				t.Errorf("shell read back %q, want %q", output, tc.value)
			}
		})
	}
}
