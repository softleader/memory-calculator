package flags

import (
	"fmt"
	"os"
)

const (
	DefaultJVMOptions = JVMOptions("")
	FlagJVMOptions    = "jvm-options"
	EnvJVMOptions     = "JAVA_OPTS"
	UsageJVMOptions   = "vm options, typically JAVA_OPTS"
)

type JVMOptions string

func NewJVMOptions() *JVMOptions {
	j := DefaultJVMOptions
	if val, ok := os.LookupEnv(EnvJVMOptions); ok {
		j = JVMOptions(val)
	}
	return &j
}

func (j *JVMOptions) Set(s string) error {
	*j = JVMOptions(s)
	return nil
}

func (j *JVMOptions) String() string {
	return fmt.Sprintf("%s", *j)
}

func (j *JVMOptions) Type() string {
	return "string"
}

func (j *JVMOptions) Contribute() error {
	if s := j.String(); s != "" {
		if err := os.Setenv(EnvJVMOptions, s); err != nil {
			return err
		}
	}
	return nil
}
