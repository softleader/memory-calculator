package calc

import (
	"fmt"
	"os"
	"strconv"
)

const (
	DefaultThreadCount = ThreadCount(200)
	FlagThreadCount    = "thread-count"
	EnvThreadCount     = "BPL_JVM_THREAD_COUNT"
	UsageThreadCount   = "positive user thread count; overrides BPL_JVM_THREAD_COUNT and skips detection (non-web/reactive: 50, servlet: 250; detection failure fallback: 200 with Warning)"
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
	return tc.parse(s, "--"+FlagThreadCount)
}

func (tc *ThreadCount) parse(s, source string) error {
	value, err := strconv.Atoi(s)
	if err != nil || value <= 0 {
		return fmt.Errorf("%s value %q must be a representable positive integer", source, s)
	}
	*tc = ThreadCount(value)
	return nil
}

// Resolve selects the explicit source before invoking the detection adapter.
// The adapter exports its result to EnvThreadCount for the JVM calculator.
func (tc *ThreadCount) Resolve(cliExplicit bool, detect func() error) error {
	if cliExplicit {
		return nil
	}
	if value, exists := os.LookupEnv(EnvThreadCount); exists {
		return tc.parse(value, EnvThreadCount)
	}
	if err := detect(); err != nil {
		return err
	}
	return tc.parse(os.Getenv(EnvThreadCount), "application detection")
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
