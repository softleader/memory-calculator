package prep

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/magiconair/properties"
	"github.com/mattn/go-shellwords"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type Jre struct {
	Logger bard.Logger
}

func NewJrePreparer(logger bard.Logger) Jre {
	jre := Jre{
		Logger: logger,
	}
	return jre
}

func (jre Jre) Prepare() error {
	var cacertsPath string

	envJavaHome, ok := os.LookupEnv("JAVA_HOME")
	if !ok {
		return fmt.Errorf("JAVA_HOME not set")
	}

	cacertsPath = filepath.Join(envJavaHome, "lib", "security", "cacerts")
	ok, err := sherpa.FileExists(cacertsPath)
	if !ok || err != nil {
		cacertsPath = ""
	}
	if err := os.Setenv("BPI_JVM_CACERTS", cacertsPath); err != nil {
		return err
	}

	var file = filepath.Join(envJavaHome, "conf", "security", "java.security")

	p, err := properties.LoadFile(file, properties.UTF8)
	if err != nil {
		return fmt.Errorf("unable to read properties file %s\n%w", file, err)
	}
	p = p.FilterStripPrefix("security.provider.")

	var providers []string
	for k, v := range p.Map() {
		providers = append(providers, fmt.Sprintf("%s|%s", k, v))
	}
	sort.Strings(providers)
	if err := os.Setenv("BPI_JVM_SECURITY_PROVIDERS", strings.Join(providers, " ")); err != nil {
		return err
	}

	var values []string
	s, ok := os.LookupEnv("JAVA_TOOL_OPTIONS")
	if ok {
		values = append(values, s)
	}

	if p, err := shellwords.Parse(s); err != nil {
		return fmt.Errorf("unable to parse $JAVA_TOOL_OPTIONS\n%w", err)
	} else {
		var hasExitOnOutOfMemoryError bool
		for _, s := range p {
			if strings.HasPrefix(s, "-XX:+ExitOnOutOfMemoryError") {
				hasExitOnOutOfMemoryError = true
			}
		}
		if !hasExitOnOutOfMemoryError {
			values = append(values, "-XX:+ExitOnOutOfMemoryError")
			if err := os.Setenv("JAVA_TOOL_OPTIONS", strings.Join(values, " ")); err != nil {
				return err
			}
		}
	}

	return nil
}
