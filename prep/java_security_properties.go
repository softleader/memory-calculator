package prep

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/libpak/bard"
)

// JavaSecurityProperties ensures that a default java.security.properties file is configured if not already present.
type JavaSecurityProperties struct {
	// Path is the directory where the default java-security.properties will be created if needed.
	Path   string
	Logger bard.Logger
}

// NewJavaSecurityProperties creates a new instance of JavaSecurityProperties.
func NewJavaSecurityProperties(logger bard.Logger, path string) JavaSecurityProperties {
	jsp := JavaSecurityProperties{Logger: logger, Path: path}
	return jsp
}

// Prepare executes the logic to ensure a java.security.properties file is configured.
func (jsp JavaSecurityProperties) Prepare() error {
	// 1. Check JAVA_TOOL_OPTIONS for -Djava.security.properties=
	if pathFromToolOptions, found, err := findJavaSecurityProperties("JAVA_TOOL_OPTIONS"); err != nil {
		return err
	} else if found {
		jsp.Logger.Debugf("Found '-Djava.security.properties=' in 'JAVA_TOOL_OPTIONS', setting JAVA_SECURITY_PROPERTIES.")
		if err := os.Setenv("JAVA_SECURITY_PROPERTIES", pathFromToolOptions); err != nil {
			return fmt.Errorf("unable to set JAVA_SECURITY_PROPERTIES from JAVA_TOOL_OPTIONS: %w", err)
		}
		return nil // Done if found in JAVA_TOOL_OPTIONS
	}

	// 2. If not found in JAVA_TOOL_OPTIONS, create a default file.
	if jsp.Path == "" {
		return fmt.Errorf("JavaSecurityProperties.Path is required to create default properties file")
	}

	defaultFilePath := filepath.Join(jsp.Path, "java-security.properties")

	jsp.Logger.Debugf("No '-Djava.security.properties=' found in 'JAVA_TOOL_OPTIONS', creating default at %s.", defaultFilePath)

	if err := os.WriteFile(defaultFilePath, []byte{}, 0644); err != nil {
		return fmt.Errorf("unable to create default security properties file at %s: %w", defaultFilePath, err)
	}

	// Set JAVA_SECURITY_PROPERTIES to the new default file.
	if err := os.Setenv("JAVA_SECURITY_PROPERTIES", defaultFilePath); err != nil {
		return fmt.Errorf("unable to set JAVA_SECURITY_PROPERTIES: %w", err)
	}

	// Also update JAVA_TOOL_OPTIONS to include this new default file.
	newOption := "-Djava.security.properties=" + defaultFilePath
	currentOpts := os.Getenv("JAVA_TOOL_OPTIONS")
	var newOpts string
	if currentOpts == "" {
		newOpts = newOption
	} else {
		newOpts = currentOpts + " " + newOption // Append to the end.
	}

	if err := os.Setenv("JAVA_TOOL_OPTIONS", newOpts); err != nil {
		return fmt.Errorf("unable to update JAVA_TOOL_OPTIONS: %w", err)
	}

	return nil
}

// findJavaSecurityProperties is a shared helper function to find the property in a given environment variable.
func findJavaSecurityProperties(envVar string) (string, bool, error) {
	s, ok := os.LookupEnv(envVar)
	if !ok {
		return "", false, nil
	}

	p := strings.Fields(s)

	for _, item := range p {
		if strings.HasPrefix(item, "-Djava.security.properties=") {
			return strings.TrimPrefix(item, "-Djava.security.properties="), true, nil
		}
	}

	return "", false, nil
}
