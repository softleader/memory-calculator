package calc

import (
	"fmt"
	"github.com/paketo-buildpacks/libjvm/count"
	"os"
	"strconv"
)

const (
	DefaultJVMClassCount = JVMClassCount(0)
	FlagJVMClassCount    = "jvm-class-count"
	EnvJVMClassCount     = "BPI_JVM_CLASS_COUNT"
	UsageJVMClassCount   = "the number of JVM classes"
	envJavaHome          = "JAVA_HOME"
)

type JVMClassCount int

func NewJVMClassCount() *JVMClassCount {
	lcc := DefaultJVMClassCount
	if val, ok := os.LookupEnv(EnvJVMClassCount); ok {
		f, _ := strconv.Atoi(val)
		lcc = JVMClassCount(f)
	}
	return &lcc
}

func (lcc *JVMClassCount) Set(s string) error {
	f, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*lcc = JVMClassCount(f)
	return nil
}

func (lcc *JVMClassCount) Type() string {
	return "int"
}

func (lcc *JVMClassCount) String() string {
	return strconv.FormatInt(int64(*lcc), 10)
}

func (lcc *JVMClassCount) Contribute() error {
	if int64(*lcc) == 0 {
		javaHome, ok := os.LookupEnv(envJavaHome)
		if !ok {
			return fmt.Errorf("failed to lookup %v env", envJavaHome)
		}
		jvmClassCount, err := count.Classes(javaHome)
		if err != nil {
			return err
		}
		*lcc = JVMClassCount(jvmClassCount)
	}
	if err := os.Setenv(EnvJVMClassCount, lcc.String()); err != nil {
		return err
	}
	return nil
}
