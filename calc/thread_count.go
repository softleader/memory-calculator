package calc

import (
	"os"
	"strconv"
)

const (
	DefaultThreadCount = ThreadCount(200)
	FlagThreadCount    = "thread-count"
	EnvThreadCount     = "BPL_JVM_THREAD_COUNT"
	UsageThreadCount   = "the number of user threads"
)

type ThreadCount int

func NewThreadCount() *ThreadCount {
	tc := DefaultThreadCount
	if val, ok := os.LookupEnv(EnvThreadCount); ok {
		i, _ := strconv.Atoi(val)
		tc = ThreadCount(i)
	}
	return &tc
}

func (tc *ThreadCount) Set(s string) error {
	f, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*tc = ThreadCount(f)
	return nil
}

func (tc *ThreadCount) Type() string {
	return "int"
}

func (tc *ThreadCount) String() string {
	return strconv.FormatInt(int64(*tc), 10)
}

func (tc *ThreadCount) Contribute() error {
	if err := os.Setenv(EnvThreadCount, tc.String()); err != nil {
		return err
	}
	return nil
}
