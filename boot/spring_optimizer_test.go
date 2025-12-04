package boot

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
