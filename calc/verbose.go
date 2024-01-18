package calc

import (
	"os"
	"strconv"
)

const (
	DefaultVerbose   = Verbose(false)
	FlagVerbose      = "verbose"
	FlagShortVerbose = "v"
	EnvVerbose       = "BP_LOG_LEVEL"
	UsageVerbose     = "enable verbose output"
	levelDebug       = "DEBUG"
)

type Verbose bool

func NewVerbose() *Verbose {
	v := DefaultVerbose
	if val, ok := os.LookupEnv(EnvEnableNmt); ok {
		v = val == levelDebug
	}
	return &v
}

func (v *Verbose) Set(s string) error {
	f, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*v = Verbose(f)
	return nil
}

func (v *Verbose) Type() string {
	return "bool"
}

func (v *Verbose) String() string {
	return strconv.FormatBool(bool(*v))
}

func (v *Verbose) Contribute() error {
	if *v {
		if err := os.Setenv(EnvVerbose, levelDebug); err != nil {
			return err
		}
	}
	return nil
}
