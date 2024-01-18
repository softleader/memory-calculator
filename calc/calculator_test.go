package calc

import (
	"os"
	"testing"
)

type MockContributor struct {
	Called bool
}

func (m *MockContributor) Contribute() error {
	m.Called = true
	return nil
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

func TestCalculator_buildCommands_HasCaCerts(t *testing.T) {
	calculator := NewCalculator()

	os.Setenv(envBpiJvmCaCerts, "some-value")
	defer os.Unsetenv(envBpiJvmCaCerts)

	cmds, err := calculator.buildCommands()
	if err != nil {
		t.Fatalf("buildCommands returned an error: %v", err)
	}

	if _, ok := cmds["openssl-certificate-loader"]; !ok {
		t.Errorf("openssl-certificate-loader command not found in cmds")
	}
}

func TestCalculator_buildCommands_NoCaCerts(t *testing.T) {
	calculator := NewCalculator()

	os.Unsetenv(envBpiJvmCaCerts)
	defer os.Unsetenv(envBpiJvmCaCerts)

	cmds, err := calculator.buildCommands()
	if err != nil {
		t.Fatalf("buildCommands returned an error: %v", err)
	}

	if _, ok := cmds["openssl-certificate-loader"]; ok {
		t.Errorf("openssl-certificate-loader command should not be present")
	}
}

func TestCalculator_buildCommands_EnableNmt(t *testing.T) {
	calculator := NewCalculator()
	*calculator.EnableNmt = true

	cmds, err := calculator.buildCommands()
	if err != nil {
		t.Fatalf("buildCommands returned an error: %v", err)
	}

	if _, ok := cmds["nmt"]; !ok {
		t.Errorf("nmt command should be present")
	}
}

func TestCalculator_buildCommands_DisableNmt(t *testing.T) {
	calculator := NewCalculator()
	*calculator.EnableNmt = false

	cmds, err := calculator.buildCommands()
	if err != nil {
		t.Fatalf("buildCommands returned an error: %v", err)
	}

	if _, ok := cmds["nmt"]; ok {
		t.Errorf("nmt command should not be present")
	}
}
