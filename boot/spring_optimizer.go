package boot

import (
	"os"
	"strings"

	"github.com/paketo-buildpacks/libpak/bard"
	boot "github.com/softleader/memory-calculator/boot/helper"
)

type SpringOptimizer struct {
	Logger         bard.Logger
	AppClassesPath *AppClassesPath
	AppLibPath     *AppLibPath
	Verbose        bool
}

func NewSpringOptimizer(logger bard.Logger) SpringOptimizer {

	so := SpringOptimizer{
		Logger:         logger,
		AppClassesPath: NewAppClassesPath(),
		AppLibPath:     NewAppLibPath(),
	}
	return so
}

func (so *SpringOptimizer) Execute() error {
	if err := so.contribute(); err != nil {
		return err
	}

	wat := boot.WebApplicationType{Logger: so.Logger}
	values, err := wat.Execute()
	if err != nil {
		return err
	}
	for k, v := range values {
		if err = os.Setenv(k, strings.TrimSpace(v)); err != nil {
			return err
		}
	}
	return nil
}

func (so *SpringOptimizer) contribute() error {
	return contribute(
		so.AppClassesPath,
		so.AppLibPath,
	)
}

func contribute(cs ...Contributor) error {
	for _, c := range cs {
		if err := c.Contribute(); err != nil {
			return err
		}
	}
	return nil
}
