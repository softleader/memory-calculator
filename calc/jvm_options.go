package calc

import (
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
	return string(*j)
}

func (j *JVMOptions) Type() string {
	return "string"
}

func (j *JVMOptions) Contribute() error {
	if s := j.String(); s != "" {
		return os.Setenv(EnvJVMOptions, s)
	}
	return nil
}
