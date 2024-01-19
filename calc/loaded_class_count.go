package calc

import (
	"os"
	"strconv"
)

const (
	DefaultLoadedClassCount = LoadedClassCount(0)
	FlagLoadedClassCount    = "loaded-class-count"
	EnvLoadedClassCount     = "BPL_JVM_LOADED_CLASS_COUNT"
	UsageLoadedClassCount   = "the number of classes that will be loaded when the app is running"
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

func (lcc *LoadedClassCount) HasValue() bool {
	return *lcc > 0
}

func (lcc *LoadedClassCount) Contribute() error {
	if err := os.Setenv(EnvLoadedClassCount, lcc.String()); err != nil {
		return err
	}
	return nil
}
