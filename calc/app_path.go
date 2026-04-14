package calc

import (
	"os"
)

const (
	DefaultAppPath = AppPath("/app")
	FlagAppPath    = "app-path"
	EnvAppPath     = "BPI_APPLICATION_PATH"
	UsageAppPath   = "the directory on the container where the app's contents are placed"
)

type AppPath string

func NewAppPath() *AppPath {
	ap := DefaultAppPath
	if val, ok := os.LookupEnv(EnvAppPath); ok {
		ap = AppPath(val)
	}
	return &ap
}

func (ap *AppPath) Set(s string) error {
	*ap = AppPath(s)
	return nil
}

func (ap *AppPath) String() string {
	return string(*ap)
}

func (ap *AppPath) Type() string {
	return "string"
}

func (ap *AppPath) Contribute() error {
	return os.Setenv(EnvAppPath, ap.String())
}
