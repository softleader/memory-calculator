package calc

import (
	"os"
	"strconv"
)

const (
	DefaultEnableJdwp = EnableJdwp(true)
	FlagEnableJdwp    = "enable-jdwp"
	EnvEnableJdwp     = "BPL_DEBUG_ENABLED"
	UsageEnableJdwp   = "enables Java Debug Wire Protocol (JDWP)"
	envBplDebugPort   = "BPL_DEBUG_PORT"
	defaultDebugPort  = 5005
)

type EnableJdwp bool

func NewEnableJdwp() *EnableJdwp {
	jdwp := DefaultEnableJdwp
	if val, ok := os.LookupEnv(EnvEnableJdwp); ok {
		if b, err := strconv.ParseBool(val); err == nil {
			jdwp = EnableJdwp(b)
		}
	}
	return &jdwp
}

func (jdwp *EnableJdwp) Set(s string) error {
	f, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*jdwp = EnableJdwp(f)
	return nil
}

func (jdwp *EnableJdwp) Type() string {
	return "bool"
}

func (jdwp *EnableJdwp) String() string {
	return strconv.FormatBool(bool(*jdwp))
}

func (jdwp *EnableJdwp) Contribute() error {
	if err := os.Setenv(EnvEnableJdwp, jdwp.String()); err != nil {
		return err
	}
	if err := os.Setenv(envBplDebugPort, strconv.Itoa(defaultDebugPort)); err != nil {
		return err
	}
	return nil
}
