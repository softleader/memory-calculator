package prep

import (
	"fmt"

	"github.com/paketo-buildpacks/libpak/bard"
)

const (
	DefaultJavaSecurityPropertiesPath = "/tmp"
)

type PreparerManager struct {
	Logger            bard.Logger
	JavaSecurityProps JavaSecurityProperties
	Jre               Jre
}

func NewPreparerManager(logger bard.Logger) PreparerManager {
	pm := PreparerManager{
		Logger:            logger,
		JavaSecurityProps: NewJavaSecurityProperties(logger, DefaultJavaSecurityPropertiesPath),
		Jre:               NewJrePreparer(logger),
	}
	return pm
}

func (p PreparerManager) PrepareAll() error {
	steps := []Preparer{
		p.JavaSecurityProps,
		p.Jre,
	}

	for i, step := range steps {
		if err := step.Prepare(); err != nil {
			return fmt.Errorf("failed to run preparer step %d: %w", i+1, err)
		}
	}

	return nil
}
