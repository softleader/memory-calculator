package prep_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/sclevine/spec"

	"github.com/softleader/memory-calculator/prep"
)

func TestJrePreparer(t *testing.T) {
	spec.Run(t, "JrePreparer", testJrePreparer)
}

func testJrePreparer(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect   = NewWithT(t).Expect
		tempDir  string
		javaHome string
		logger   bard.Logger
	)

	it.Before(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "jre-test")
		Expect(err).NotTo(HaveOccurred())

		javaHome = filepath.Join(tempDir, "java_home")
		Expect(os.MkdirAll(filepath.Join(javaHome, "lib", "security"), 0755)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(javaHome, "conf", "security"), 0755)).To(Succeed())

		// Create a dummy java.security file
		javaSecurityContent := `
security.provider.1=SunProvider
security.provider.2=BouncyCastleProvider
some.other.property=value
`
		Expect(os.WriteFile(filepath.Join(javaHome, "conf", "security", "java.security"), []byte(javaSecurityContent), 0644)).To(Succeed())

		logger = bard.NewLogger(io.Discard)

		// Clean up environment variables before each test
		Expect(os.Unsetenv("JAVA_HOME")).To(Succeed())
		Expect(os.Unsetenv("BPI_JVM_CACERTS")).To(Succeed())
		Expect(os.Unsetenv("BPI_JVM_SECURITY_PROVIDERS")).To(Succeed())
		Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
	})

	it.After(func() {
		// Clean up environment variables after each test
		Expect(os.Unsetenv("JAVA_HOME")).To(Succeed())
		Expect(os.Unsetenv("BPI_JVM_CACERTS")).To(Succeed())
		Expect(os.Unsetenv("BPI_JVM_SECURITY_PROVIDERS")).To(Succeed())
		Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
		Expect(os.RemoveAll(tempDir)).To(Succeed())
	})

	context("when JAVA_HOME is not set", func() {
		it("returns an error", func() {
			jsp := prep.NewJrePreparer(logger)
			Expect(jsp.Prepare()).To(MatchError("JAVA_HOME not set"))
		})
	})

	context("when JAVA_HOME is set", func() {
		it.Before(func() {
			Expect(os.Setenv("JAVA_HOME", javaHome)).To(Succeed())
		})

		context("when cacerts file is missing", func() {
			it("sets BPI_JVM_CACERTS to an empty string", func() {
				jsp := prep.NewJrePreparer(logger)
				Expect(jsp.Prepare()).To(Succeed())
				Expect(os.Getenv("BPI_JVM_CACERTS")).To(BeEmpty())
			})
		})

		context("when java.security file is missing", func() {
			it.Before(func() {
				Expect(os.Remove(filepath.Join(javaHome, "conf", "security", "java.security"))).To(Succeed())
			})

			it("returns an error", func() {
				jsp := prep.NewJrePreparer(logger)
				Expect(jsp.Prepare()).To(MatchError(ContainSubstring("unable to read properties file")))
			})
		})

		context("when all conditions are met", func() {
			it.Before(func() {
				// Create a dummy cacerts file
				Expect(os.WriteFile(filepath.Join(javaHome, "lib", "security", "cacerts"), []byte{}, 0644)).To(Succeed())
			})

			it("sets BPI_JVM_CACERTS, BPI_JVM_SECURITY_PROVIDERS and appends to JAVA_TOOL_OPTIONS", func() {
				jsp := prep.NewJrePreparer(logger)
				Expect(jsp.Prepare()).To(Succeed())

				// Verify BPI_JVM_CACERTS
				Expect(os.Getenv("BPI_JVM_CACERTS")).To(Equal(filepath.Join(javaHome, "lib", "security", "cacerts")))

				// Verify BPI_JVM_SECURITY_PROVIDERS
				Expect(os.Getenv("BPI_JVM_SECURITY_PROVIDERS")).To(Equal("1|SunProvider 2|BouncyCastleProvider"))

				// Verify JAVA_TOOL_OPTIONS
				Expect(os.Getenv("JAVA_TOOL_OPTIONS")).To(ContainSubstring("-XX:+ExitOnOutOfMemoryError"))
			})

			context("when -XX:+ExitOnOutOfMemoryError is already in JAVA_TOOL_OPTIONS", func() {
				it.Before(func() {
					Expect(os.Setenv("JAVA_TOOL_OPTIONS", "-Xmx1g -XX:+ExitOnOutOfMemoryError")).To(Succeed())
				})

				it("does not append it again", func() {
					jsp := prep.NewJrePreparer(logger)
					Expect(jsp.Prepare()).To(Succeed())

					// Verify JAVA_TOOL_OPTIONS contains it once
					opts := os.Getenv("JAVA_TOOL_OPTIONS")
					Expect(strings.Count(opts, "-XX:+ExitOnOutOfMemoryError")).To(Equal(1))
				})
			})
		})
	})
}
