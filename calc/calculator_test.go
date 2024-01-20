package calc

import (
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

func TestCalculator_buildHelpers_HasCaCerts(t *testing.T) {
	calculator := NewCalculator()
	calculator.JVMCacerts.Set("some-value")

	helpers, err := calculator.buildHelpers()
	if err != nil {
		t.Fatalf("buildHelpers returned an error: %v", err)
	}

	if _, ok := helpers[helperOpensslCertificateLoader]; !ok {
		t.Errorf(helperOpensslCertificateLoader + " helper not found in helpers")
	}
}

func TestCalculator_buildHelpers_NoCaCerts(t *testing.T) {
	calculator := NewCalculator()

	helpers, err := calculator.buildHelpers()
	if err != nil {
		t.Fatalf("buildHelpers returned an error: %v", err)
	}

	if _, ok := helpers[helperOpensslCertificateLoader]; ok {
		t.Errorf(helperOpensslCertificateLoader + " helper should not be present")
	}
}

func TestCalculator_buildHelpers_EnableNmt(t *testing.T) {
	calculator := NewCalculator()
	*calculator.EnableNmt = true

	helpers, err := calculator.buildHelpers()
	if err != nil {
		t.Fatalf("buildHelpers returned an error: %v", err)
	}

	if _, ok := helpers[helperNmt]; !ok {
		t.Errorf(helperNmt + " helper should be present")
	}
}

func TestCalculator_buildHelpers_DisableNmt(t *testing.T) {
	calculator := NewCalculator()
	*calculator.EnableNmt = false

	helpers, err := calculator.buildHelpers()
	if err != nil {
		t.Fatalf("buildHelpers returned an error: %v", err)
	}

	if _, ok := helpers[helperNmt]; ok {
		t.Errorf(helperNmt + " helper should not be present")
	}
}
