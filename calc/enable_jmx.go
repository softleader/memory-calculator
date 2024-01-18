package calc

import (
	"os"
	"strconv"
)

const (
	DefaultEnableJmx = EnableJmx(false)
	FlagEnableJmx    = "enable-jmx"
	EnvEnableJmx     = "BPL_JMX_ENABLED"
	UsageEnableJmx   = "enable Java Management Extensions (JMX)"
)

type EnableJmx bool

func NewEnableJmx() *EnableJmx {
	jmx := DefaultEnableJmx
	if val, ok := os.LookupEnv(EnvEnableJmx); ok {
		if b, err := strconv.ParseBool(val); err == nil {
			jmx = EnableJmx(b)
		}
	}
	return &jmx
}

func (jmx *EnableJmx) Set(s string) error {
	f, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*jmx = EnableJmx(f)
	return nil
}

func (jmx *EnableJmx) Type() string {
	return "bool"
}

func (jmx *EnableJmx) String() string {
	return strconv.FormatBool(bool(*jmx))
}

func (jmx *EnableJmx) Contribute() error {
	if err := os.Setenv(EnvEnableJmx, jmx.String()); err != nil {
		return err
	}
	return nil
}
