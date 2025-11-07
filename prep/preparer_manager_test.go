package prep_test

import (
	"errors"
	"io"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/sclevine/spec"

	"github.com/softleader/memory-calculator/prep"
)

// MockPreparer is a test double for types that implement the prep.Preparer interface.
// It allows us to track the order of execution and simulate errors.
type MockPreparer struct {
	// A slice to record the name of the preparer when Prepare() is called.
	callOrder *[]string
	// The name of this specific mock, to be appended to callOrder.
	name string
	// If this error is non-nil, the Prepare() method will return it.
	err error
}

// Prepare records that it was called and returns a simulated error if one is set.
func (m MockPreparer) Prepare() error {
	*m.callOrder = append(*m.callOrder, m.name)
	if m.err != nil {
		return m.err
	}
	return nil
}

func TestPreparerManager(t *testing.T) {
	spec.Run(t, "PreparerManager", testPreparerManager)
}

func testPreparerManager(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
	)

	context("when all preparers succeed", func() {
		it("executes them in the correct order", func() {
			callOrder := &[]string{}

			// Create a PreparerManager and inject our mock preparers
			pm := prep.PreparerManager{
				Logger: bard.NewLogger(io.Discard),
				Preparers: []prep.Preparer{
					MockPreparer{callOrder: callOrder, name: "first"},
					MockPreparer{callOrder: callOrder, name: "second"},
				},
			}

			Expect(pm.PrepareAll()).To(Succeed())

			// Verify the execution order
			Expect(*callOrder).To(Equal([]string{"first", "second"}))
		})
	})

	context("when a preparer fails", func() {
		it("stops execution and returns the error", func() {
			callOrder := &[]string{}
			simulatedError := errors.New("preparer failed")

			pm := prep.PreparerManager{
				Logger: bard.NewLogger(io.Discard),
				Preparers: []prep.Preparer{
					MockPreparer{callOrder: callOrder, name: "first"},
					MockPreparer{callOrder: callOrder, name: "second", err: simulatedError},
					MockPreparer{callOrder: callOrder, name: "third"},
				},
			}

			Expect(pm.PrepareAll()).To(MatchError(ContainSubstring("preparer failed")))

			// Verify that execution stopped after the failing preparer
			Expect(*callOrder).To(Equal([]string{"first", "second"}))
		})
	})
}
