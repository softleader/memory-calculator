package calc

import (
	"os"
	"strconv"
)

const (
	DefaultJVMClassAdj = JVMClassAdj(0)
	FlagJVMClassAdj    = "jvm-class-adj"
	EnvJVMClassAdj     = "BPL_JVM_CLASS_ADJUSTMENT"
	UsageJVMClassAdj   = "the adjustment for JVM classes number"
)

type JVMClassAdj int

func NewJVMClassAdj() *JVMClassAdj {
	lcc := DefaultJVMClassAdj
	if val, ok := os.LookupEnv(EnvJVMClassAdj); ok {
		f, _ := strconv.Atoi(val)
		lcc = JVMClassAdj(f)
	}
	return &lcc
}

func (jca *JVMClassAdj) Set(s string) error {
	f, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*jca = JVMClassAdj(f)
	return nil
}

func (jca *JVMClassAdj) Type() string {
	return "int"
}

func (jca *JVMClassAdj) String() string {
	return strconv.FormatInt(int64(*jca), 10)
}

func (jca *JVMClassAdj) Contribute() error {
	if *jca > 0 {
		if err := os.Setenv(EnvJVMClassAdj, jca.String()); err != nil {
			return err
		}
	}
	return nil
}
