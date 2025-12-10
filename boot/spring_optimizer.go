package boot

import (
	"os"
	"strings"

	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
	boot "github.com/softleader/memory-calculator/boot/helper"
)

const (
	helperWebApplicationType = "web-application-type"
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

	hs, err := so.buildHelpers()
	if err != nil {
		return err
	}

	inOrder := []string{
		helperWebApplicationType,
	}

	// 按照指定順序執行
	for _, key := range inOrder {
		h, ok := hs[key]
		if !ok {
			continue
		}
		values, err := h.Execute()
		if err != nil {
			return err
		}
		for k, v := range values {
			v = strings.TrimSpace(v)
			if err = os.Setenv(k, v); err != nil { // update golang environment variable
				return err
			}
		}
	}

	return nil
}

func (so *SpringOptimizer) buildHelpers() (h map[string]sherpa.ExecD, err error) {
	var (
		l   = so.Logger
		wat = boot.WebApplicationType{Logger: l}
	)

	h = map[string]sherpa.ExecD{
		helperWebApplicationType: wat,
	}

	return h, nil
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
