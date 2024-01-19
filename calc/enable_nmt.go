package calc

import (
	"os"
	"strconv"
)

const (
	DefaultEnableNmt = EnableNmt(false)
	FlagEnableNmt    = "enable-nmt"
	EnvEnableNmt     = "BPL_JAVA_NMT_ENABLED"
	UsageEnableNmt   = "enables Native Memory Tracking (NMT)"
)

type EnableNmt bool

func NewEnableNmt() *EnableNmt {
	nmt := DefaultEnableNmt
	if val, ok := os.LookupEnv(EnvEnableNmt); ok {
		if b, err := strconv.ParseBool(val); err == nil {
			nmt = EnableNmt(b)
		}
	}
	return &nmt
}

func (nmt *EnableNmt) Set(s string) error {
	f, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*nmt = EnableNmt(f)
	return nil
}

func (nmt *EnableNmt) Type() string {
	return "bool"
}

func (nmt *EnableNmt) String() string {
	return strconv.FormatBool(bool(*nmt))
}

func (nmt *EnableNmt) Contribute() error {
	if err := os.Setenv(EnvEnableNmt, nmt.String()); err != nil {
		return err
	}
	return nil
}
