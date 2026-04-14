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
	subPathCacerts    = "/lib/security/cacerts"
	UsageJVMCacerts   = "path to jvm cacerts, typically 'JAVA_HOME" + subPathCacerts + "'"
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
			cacert := filepath.Join(javaHome, subPathCacerts)
			f, err := os.Open(cacert)
			if err == nil {
				defer func(f *os.File) {
					_ = f.Close()
				}(f)
				*j = JVMCacerts(cacert)
			}
		}
	}
	if s := j.String(); s != "" {
		if err := os.Setenv(EnvJVMCacerts, s); err != nil {
			return err
		}
	}
	return nil
}
