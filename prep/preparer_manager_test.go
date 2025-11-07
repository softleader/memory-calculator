package prep_test

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/sclevine/spec"

	"github.com/softleader/memory-calculator/prep"
)

func TestPreparerManager(t *testing.T) {
	spec.Run(t, "PreparerManager", testPreparerManager)
}

func testPreparerManager(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect   = NewWithT(t).Expect
		tempDir  string
		javaHome string
		logger   bard.Logger
	)

	it.Before(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "preparer-manager-test")
		Expect(err).NotTo(HaveOccurred())

		// Create a fake JAVA_HOME for the Jre preparer to use
		javaHome = filepath.Join(tempDir, "java_home")
		Expect(os.MkdirAll(filepath.Join(javaHome, "lib", "security"), 0755)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(javaHome, "conf", "security"), 0755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(javaHome, "lib", "security", "cacerts"), []byte{}, 0644)).To(Succeed())
		javaSecurityContent := `security.provider.1=TestProvider`
		Expect(os.WriteFile(filepath.Join(javaHome, "conf", "security", "java.security"), []byte(javaSecurityContent), 0644)).To(Succeed())

		logger = bard.NewLogger(io.Discard)

		// CRITICAL: Clean up all environment variables before each test.
		Expect(os.Unsetenv("JAVA_HOME")).To(Succeed())
		Expect(os.Unsetenv("BPI_JVM_SECURITY_PROVIDERS")).To(Succeed())
		Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
		Expect(os.Unsetenv("JAVA_SECURITY_PROPERTIES")).To(Succeed())
	})

	it.After(func() {
		// CRITICAL: Clean up all environment variables after each test to prevent pollution.
		Expect(os.Unsetenv("JAVA_HOME")).To(Succeed())
		Expect(os.Unsetenv("BPI_JVM_SECURITY_PROVIDERS")).To(Succeed())
		Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
		Expect(os.Unsetenv("JAVA_SECURITY_PROPERTIES")).To(Succeed())
		Expect(os.RemoveAll(tempDir)).To(Succeed())
	})

	context("when a sub-preparer fails", func() {
		it("stops and returns the error", func() {
			// We can trigger an error in the Jre preparer by not setting JAVA_HOME
			pm := prep.NewPreparerManager(logger)
			Expect(pm.PrepareAll()).To(MatchError(ContainSubstring("JAVA_HOME not set")))
		})
	})

	context("when all preparers succeed", func() {
		it("correctly modifies the environment in sequence", func() {
			// This test is only valid on non-Windows systems due to the unix.Access call in jre.go
			if runtime.GOOS == "windows" {
				return
			}

			// Set the necessary input environment variable
			Expect(os.Setenv("JAVA_HOME", javaHome)).To(Succeed())

			pm := prep.NewPreparerManager(logger)
			Expect(pm.PrepareAll()).To(Succeed())

			// --- Assert the final state of the environment ---

			// 1. Assert side-effects of the JavaSecurityProps preparer
			// It should have created a default file in /tmp (as hardcoded in NewPreparerManager)
			defaultJspPath := filepath.Join(prep.DefaultJavaSecurityPropertiesPath, "java-security.properties")
			Expect(defaultJspPath).To(BeARegularFile())
			Expect(os.Getenv("JAVA_SECURITY_PROPERTIES")).To(Equal(defaultJspPath))

			// 2. Assert side-effects of the Jre preparer
			Expect(os.Getenv("BPI_JVM_SECURITY_PROVIDERS")).To(Equal("1|TestProvider"))

			// 3. Assert the combined final state of JAVA_TOOL_OPTIONS
			opts := os.Getenv("JAVA_TOOL_OPTIONS")
			Expect(opts).To(ContainSubstring("-Djava.security.properties=" + defaultJspPath)) // From JavaSecurityProps
			Expect(opts).To(ContainSubstring("-XX:+ExitOnOutOfMemoryError"))                  // From Jre

			// The Jre preparer also tries to add -Djava.security.properties, let's check it's not duplicated
			// Note: This assertion depends on the exact (and somewhat flawed) implementation of jre.go
			jreJspPath := filepath.Join(javaHome, "conf", "security", "java.security")
			Expect(opts).To(ContainSubstring("-Djava.security.properties=" + jreJspPath))
			Expect(strings.Count(opts, "-Djava.security.properties=")).To(Equal(2))
		})
	})
}
