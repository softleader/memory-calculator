package boot

import (
	"os"
)

const (
	DefaultAppClassesPath = AppClassesPath("/app/classes")
	FlagAppClassesPath    = "app-classes-path"
	EnvAppClassesPath     = "APPLICATION_CLASSES_PATH"
	UsageAppClassesPath   = ""
)

type AppClassesPath string

func NewAppClassesPath() *AppClassesPath {
	acp := DefaultAppClassesPath
	if val, ok := os.LookupEnv(EnvAppClassesPath); ok {
		acp = AppClassesPath(val)
	}
	return &acp
}

func (acp *AppClassesPath) Set(s string) error {
	*acp = AppClassesPath(s)
	return nil
}

func (acp *AppClassesPath) String() string {
	return string(*acp)
}

func (acp *AppClassesPath) Type() string {
	return "string"
}

func (acp *AppClassesPath) Contribute() error {
	return os.Setenv(EnvAppClassesPath, acp.String())
}
