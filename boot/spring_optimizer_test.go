package boot

import (
	"bytes"
	"os"
	"testing"

	"github.com/paketo-buildpacks/libpak/bard"
	springboot "github.com/paketo-buildpacks/spring-boot/v5/boot"
	boot "github.com/softleader/memory-calculator/boot/helper"
)

type MockContributor struct {
	Called bool
}

func (m *MockContributor) Contribute() error {
	m.Called = true
	return nil
}

func TestSpringOptimizer_Execute(t *testing.T) {
	original := boot.ResolveWebAppType
	boot.ResolveWebAppType = func() (springboot.ApplicationType, error) {
		return springboot.Servlet, nil
	}
	t.Cleanup(func() { boot.ResolveWebAppType = original })

	envBplJvmThreadCount := "BPL_JVM_THREAD_COUNT"
	os.Unsetenv(EnvAppClassesPath)
	os.Unsetenv(EnvAppLibPath)
	os.Unsetenv(envBplJvmThreadCount)
	t.Cleanup(func() {
		os.Unsetenv(EnvAppClassesPath)
		os.Unsetenv(EnvAppLibPath)
		os.Unsetenv(envBplJvmThreadCount)
	})

	so := NewSpringOptimizer(bard.NewLogger(&bytes.Buffer{}))
	if err := so.Execute(); err != nil {
		t.Fatalf("Execute returned an error: %v", err)
	}

	if v, ok := os.LookupEnv(envBplJvmThreadCount); !ok || v != "250" {
		t.Errorf("Expected %s to be '250', got '%v'", envBplJvmThreadCount, v)
	}
	if v, ok := os.LookupEnv(EnvAppClassesPath); !ok || v != string(DefaultAppClassesPath) {
		t.Errorf("Expected %s to be '%s', got '%v'", EnvAppClassesPath, DefaultAppClassesPath, v)
	}
	if v, ok := os.LookupEnv(EnvAppLibPath); !ok || v != string(DefaultAppLibPath) {
		t.Errorf("Expected %s to be '%s', got '%v'", EnvAppLibPath, DefaultAppLibPath, v)
	}
}

func TestCalculator_Contribute(t *testing.T) {
	mockContributor := &MockContributor{}

	err := contribute(mockContributor)
	if err != nil {
		t.Fatalf("contribute returned an error: %v", err)
	}

	if !mockContributor.Called {
		t.Errorf("Contributor's Contribute method was not called")
	}
}
