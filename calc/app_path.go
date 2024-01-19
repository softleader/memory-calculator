package calc

import (
	"fmt"
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
	return &ap
}

func (ap *AppPath) Set(s string) error {
	*ap = AppPath(s)
	return nil
}

func (ap *AppPath) String() string {
	return fmt.Sprintf("%s", *ap)
}

func (ap *AppPath) Type() string {
	return "string"
}

func (ap *AppPath) Contribute() error {
	if err := os.Setenv(EnvAppPath, ap.String()); err != nil {
		return err
	}
	return nil
}
