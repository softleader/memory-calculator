package calc

import (
	"fmt"
	"os"
	"strings"
)

const (
	EnvJavaToolOptions = "JAVA_TOOL_OPTIONS"
	OptionsSeparator   = " "
)

// ContributeOptions 固定要貢獻 jvm 的參數, 這些參數可能是 libjvm 在 build image 時加的而非計算出來的, 或是我們自己想要加上去的都可以放
var ContributeOptions = []string{"-XX:+ExitOnOutOfMemoryError"}

type JavaToolOptions string

func BuildJavaToolOptions() *JavaToolOptions {
	o := ""
	if val, ok := os.LookupEnv(EnvJavaToolOptions); ok {
		o = val
	}
	for _, option := range ContributeOptions {
		if !containsOption(o, option) {
			if o == "" {
				o = option
			} else {
				o += OptionsSeparator + option
			}
		}
	}
	j := JavaToolOptions(o)
	return &j
}

func (j *JavaToolOptions) String() string {
	return string(*j)
}

func (j *JavaToolOptions) Print() {
	fmt.Printf("%v: %v\n", EnvJavaToolOptions, j.String())
}

func (j *JavaToolOptions) WriteFile(file string) error {
	out, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to create file %v: %v", file, err)
	}
	defer closeFile(out)

	_, err = out.WriteString(fmt.Sprintf("export %v='%s'\n", EnvJavaToolOptions, j.String()))
	if err != nil {
		return fmt.Errorf("failed to write file %v: %v", file, err)
	}
	return nil
}

func containsOption(existing, option string) bool {
	for _, tok := range strings.Fields(existing) {
		if tok == option {
			return true
		}
	}
	return false
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Printf("WARNING: failed to close file %v: %v\n", file.Name(), err)
	}
}
