package calc

import (
	"os"
	"strconv"
	"strings"
)

const (
	DefaultVerbose   = Verbose(false)
	FlagVerbose      = "verbose"
	FlagShortVerbose = "v"
	EnvVerbose       = "BP_LOG_LEVEL"
	EnvDebug         = "BP_DEBUG"
	UsageVerbose     = "enable verbose output"
	levelDebug       = "DEBUG"
)

type Verbose bool

// NewVerbose 判斷邏輯對齊 libpak bard.Logger: $BP_DEBUG 有設定(不論值)或 $BP_LOG_LEVEL 不分大小寫等於 DEBUG 即開啟
func NewVerbose() *Verbose {
	v := DefaultVerbose
	if _, ok := os.LookupEnv(EnvDebug); ok {
		v = true
	} else if val, ok := os.LookupEnv(EnvVerbose); ok {
		v = Verbose(strings.EqualFold(val, levelDebug))
	}
	return &v
}

func (v *Verbose) Set(s bool) {
	*v = Verbose(s)
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
