package calc

import (
	"fmt"
	"os"
)

const (
	DefaultJVMClassAdj = JVMClassAdj("")
	FlagJVMClassAdj    = "jvm-class-adj"
	EnvJVMClassAdj     = "BPL_JVM_CLASS_ADJUSTMENT"
	UsageJVMClassAdj   = "the adjustment for the number or percentage of JVM classes"
)

type JVMClassAdj string

func NewJVMClassAdj() *JVMClassAdj {
	j := DefaultJVMClassAdj
	if val, ok := os.LookupEnv(EnvJVMClassAdj); ok {
		j = JVMClassAdj(val)
	}
	return &j
}

func (j *JVMClassAdj) Set(s string) error {
	*j = JVMClassAdj(s)
	return nil
}

func (j *JVMClassAdj) String() string {
	return fmt.Sprintf("%s", *j)
}

func (j *JVMClassAdj) Type() string {
	return "string"
}

func (j *JVMClassAdj) Contribute() error {
	if s := j.String(); s != "" {
		if err := os.Setenv(EnvJVMClassAdj, s); err != nil {
			return err
		}
	}
	return nil
}
