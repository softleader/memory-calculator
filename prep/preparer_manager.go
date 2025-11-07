package prep

import (
	"fmt"

	"github.com/paketo-buildpacks/libpak/bard"
)

const (
	DefaultJavaSecurityPropertiesPath = "/tmp"
)

// PreparerManager is a coordinator for all preparation steps.
type PreparerManager struct {
	Logger    bard.Logger
	Preparers []Preparer
}

// NewPreparerManager creates a new instance of the PreparerManager coordinator,
// and populates it with the default, ordered list of preparers.
func NewPreparerManager(logger bard.Logger) PreparerManager {
	return PreparerManager{
		Logger: logger,
		Preparers: []Preparer{
			NewJavaSecurityProperties(logger, DefaultJavaSecurityPropertiesPath),
			NewJrePreparer(logger),
		},
	}
}

// Prepare executes all registered preparation steps in order.
func (p PreparerManager) Prepare() error {
	for i, step := range p.Preparers {
		if err := step.Prepare(); err != nil {
			return fmt.Errorf("failed to run preparer step %d: %w", i+1, err)
		}
	}

	return nil
}
