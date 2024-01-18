package flags

import (
	"os"
	"strconv"
)

const (
	DefaultEnableJfr = EnableJfr(false)
	FlagEnableJfr    = "enable-jfr"
	EnvEnableJfr     = "BPL_JFR_ENABLED"
	UsageEnableJfr   = "enable Java Flight Recorder (JFR)"
)

type EnableJfr bool

func NewEnableJfr() *EnableJfr {
	jfr := DefaultEnableJfr
	if val, ok := os.LookupEnv(EnvEnableJfr); ok {
		if b, err := strconv.ParseBool(val); err == nil {
			jfr = EnableJfr(b)
		}
	}
	return &jfr
}

func (jfr *EnableJfr) Set(s string) error {
	f, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*jfr = EnableJfr(f)
	return nil
}

func (jfr *EnableJfr) Type() string {
	return "bool"
}

func (jfr *EnableJfr) String() string {
	return strconv.FormatBool(bool(*jfr))
}

func (jfr *EnableJfr) Contribute() error {
	if err := os.Setenv(EnvEnableJfr, jfr.String()); err != nil {
		return err
	}
	return nil
}
