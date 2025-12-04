package boot

import (
	"os"

	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/spring-boot/v5/boot"
)

type WebApplicationType struct {
	Logger bard.Logger
}

// Reference implementation from:
// "Paketo Buildpack for Spring Boot"
// https://github.com/paketo-buildpacks/spring-boot/blob/main/boot/web_application_type.go
func (wat WebApplicationType) Execute() (map[string]string, error) {
	webAppType, err := ResolveWebAppType()
	if err != nil {
		return nil, err
	}

	threadCount := "200"
	switch webAppType {
	case boot.None:
		wat.Logger.Info("Non-web application detected")
		threadCount = "50"
	case boot.Reactive:
		wat.Logger.Info("Reactive web application detected")
		threadCount = "50"
	case boot.Servlet:
		wat.Logger.Info("Servlet web application detected")
		threadCount = "250"
	}

	return map[string]string{"BPL_JVM_THREAD_COUNT": threadCount}, nil
}

var ResolveWebAppType = func() (boot.ApplicationType, error) {
	resolver, err := boot.NewWebApplicationResolver(os.Getenv("APPLICATION_CLASSES_PATH"), os.Getenv("APPLICATION_LIB_PATH"))
	if err != nil {
		return boot.None, err
	}
	return resolver.Resolve(), nil
}
