package calc

import (
	"fmt"
	"github.com/paketo-buildpacks/libjvm/count"
	"os"
	"strconv"
)

const (
	DefaultLoadedClassCount = LoadedClassCount(0)
	FlagLoadedClassCount    = "loaded-class-count"
	EnvLoadedClassCount     = "BPL_JVM_LOADED_CLASS_COUNT"
	UsageLoadedClassCount   = "the number of classes that will be loaded when the app is running"
	envJavaHome             = "JAVA_HOME"
)

type LoadedClassCount int

func NewLoadedClassCount() *LoadedClassCount {
	lcc := DefaultLoadedClassCount
	if val, ok := os.LookupEnv(EnvLoadedClassCount); ok {
		f, _ := strconv.Atoi(val)
		lcc = LoadedClassCount(f)
	}
	return &lcc
}

func (lcc *LoadedClassCount) Set(s string) error {
	f, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*lcc = LoadedClassCount(f)
	return nil
}

func (lcc *LoadedClassCount) Type() string {
	return "int"
}

func (lcc *LoadedClassCount) String() string {
	return strconv.FormatInt(int64(*lcc), 10)
}

func (lcc *LoadedClassCount) Contribute() error {
	if int64(*lcc) == 0 {
		if javaHome, ok := os.LookupEnv(envJavaHome); !ok {
			return fmt.Errorf("failed to lookup %v env", envJavaHome)
		} else {
			jvmClassCount, err := count.Classes(javaHome)
			if err != nil {
				return err
			}
			*lcc = LoadedClassCount(jvmClassCount)
		}
	}
	return nil
}
