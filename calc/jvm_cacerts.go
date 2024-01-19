package calc

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultJVMCacerts = JVMCacerts("")
	FlagJVMCacerts    = "jvm-cacerts"
	EnvJVMCacerts     = "BPI_JVM_CACERTS"
	cacertsSubPath    = "/lib/security/cacerts"
	UsageJVMCacerts   = "path to jvm cacerts, typically 'JAVA_HOME" + cacertsSubPath + "'"
)

type JVMCacerts string

func NewJVMCacerts() *JVMCacerts {
	j := DefaultJVMCacerts
	if val, ok := os.LookupEnv(EnvJVMCacerts); ok {
		j = JVMCacerts(val)
	}
	return &j
}

func (j *JVMCacerts) Set(s string) error {
	*j = JVMCacerts(s)
	return nil
}

func (j *JVMCacerts) String() string {
	return fmt.Sprintf("%s", *j)
}

func (j *JVMCacerts) Type() string {
	return "string"
}

func (j *JVMCacerts) Contribute() error {
	if s := j.String(); s == "" {
		if javaHome, ok := os.LookupEnv(envJavaHome); ok {
			cacert := filepath.Join(javaHome, cacertsSubPath)
			if exist, err := isFileExist(cacert); exist && err == nil {
				*j = JVMCacerts(cacert)
			}
		}
	}
	if *j != "" {
		if err := os.Setenv(EnvJVMCacerts, j.String()); err != nil {
			return err
		}
	}
	return nil
}

func isFileExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !fileInfo.IsDir(), nil
}
