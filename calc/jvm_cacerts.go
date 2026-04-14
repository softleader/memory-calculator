package calc

import (
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
	return string(*j)
}

func (j *JVMCacerts) Type() string {
	return "string"
}

func (j *JVMCacerts) Contribute() error {
	cacert := j.String()
	if cacert == "" {
		if javaHome, ok := os.LookupEnv(envJavaHome); ok {
			path := filepath.Join(javaHome, subPathCacerts)
			if f, err := os.Open(path); err == nil {
				f.Close()
				*j = JVMCacerts(path)
				cacert = path
			}
		}
	}
	if cacert != "" {
		return os.Setenv(EnvJVMCacerts, cacert)
	}
	return nil
}
