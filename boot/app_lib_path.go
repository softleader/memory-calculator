package boot

import (
	"fmt"
	"os"
)

const (
	DefaultAppLibPath = AppLibPath("/app/libs")
	FlagAppLibPath    = "app-lib-path"
	EnvAppLibPath     = "APPLICATION_LIB_PATH"
	UsageAppLibPath   = ""
)

type AppLibPath string

func NewAppLibPath() *AppLibPath {
	alp := DefaultAppLibPath
	if val, ok := os.LookupEnv(EnvAppLibPath); ok {
		alp = AppLibPath(val)
	}
	return &alp
}

func (alp *AppLibPath) Set(s string) error {
	*alp = AppLibPath(s)
	return nil
}

func (alp *AppLibPath) String() string {
	return fmt.Sprintf("%s", *alp)
}

func (alp *AppLibPath) Type() string {
	return "string"
}

func (alp *AppLibPath) Contribute() error {
	if err := os.Setenv(EnvAppLibPath, alp.String()); err != nil {
		return err
	}
	return nil
}
