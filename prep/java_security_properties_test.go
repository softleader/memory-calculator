package prep_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/sclevine/spec"

	"github.com/softleader/memory-calculator/prep"
)

func TestJavaSecurityProperties(t *testing.T) {
	spec.Run(t, "JavaSecurityProperties", testJavaSecurityProperties)
}

func testJavaSecurityProperties(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect  = NewWithT(t).Expect
		tempDir string
	)

	it.Before(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "jsp-test")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(tempDir)).To(Succeed())
	})

	// This is the main test case, where no security properties are configured beforehand.
	context("when no security properties are configured", func() {
		it.Before(func() {
			Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
			Expect(os.Unsetenv("JAVA_SECURITY_PROPERTIES")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
			Expect(os.Unsetenv("JAVA_SECURITY_PROPERTIES")).To(Succeed())
		})

		it("creates a default file and sets environment variables", func() {
			jsp := prep.NewJavaSecurityProperties(bard.NewLogger(io.Discard), tempDir)

			Expect(jsp.Prepare()).To(Succeed())

			// 1. Check that the new file was created
			defaultFilePath := filepath.Join(tempDir, "java-security.properties")
			Expect(defaultFilePath).To(BeARegularFile())

			// 2. Check that both environment variables are set correctly
			Expect(os.Getenv("JAVA_SECURITY_PROPERTIES")).To(Equal(defaultFilePath))
			Expect(os.Getenv("JAVA_TOOL_OPTIONS")).To(ContainSubstring("-Djava.security.properties=" + defaultFilePath))
		})
	})

	// This tests the "back-fill" logic.
	context("when JAVA_TOOL_OPTIONS already defines the property", func() {
		var userDefinedPath string

		it.Before(func() {
			userDefinedPath = filepath.Join(tempDir, "user.properties")
			Expect(os.WriteFile(userDefinedPath, []byte("test=prop"), 0644)).To(Succeed())

			Expect(os.Setenv("JAVA_TOOL_OPTIONS", "-Djava.security.properties="+userDefinedPath)).To(Succeed())
			Expect(os.Unsetenv("JAVA_SECURITY_PROPERTIES")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
			Expect(os.Unsetenv("JAVA_SECURITY_PROPERTIES")).To(Succeed())
		})

		it("back-fills JAVA_SECURITY_PROPERTIES and does not create a new file", func() {
			jsp := prep.NewJavaSecurityProperties(bard.NewLogger(io.Discard), tempDir)
			// Note: We do not set jsp.Path here to ensure it's not used.

			Expect(jsp.Prepare()).To(Succeed())

			// 1. Check that JAVA_SECURITY_PROPERTIES is correctly back-filled
			Expect(os.Getenv("JAVA_SECURITY_PROPERTIES")).To(Equal(userDefinedPath))

			// 2. Check that no new default file was created
			defaultFilePath := filepath.Join(tempDir, "java-security.properties")
			Expect(defaultFilePath).NotTo(BeAnExistingFile())
		})
	})
}
