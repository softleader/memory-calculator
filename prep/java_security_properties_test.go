package prep_test

import (
	"errors"
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

func TestJavaSecurityProperties_PrepareExistingPath(t *testing.T) {
	for _, tc := range []struct {
		name     string
		filename string
		quote    string
		before   string
	}{
		{name: "ordinary", filename: "security.properties"},
		{name: "double quoted spaces", filename: "security policy.properties", quote: `"`},
		{name: "single quoted spaces", filename: "security policy.properties", quote: "'"},
		{name: "apostrophe in path", filename: "owner's security.properties", quote: `"`},
		{name: "literal variable", filename: "$VAR security.properties", quote: `"`},
		{name: "URL before property", filename: "security.properties", before: "-Durl=https://example.test/?a=1&b=2 "},
		{name: "shell operators before property", filename: "security.properties", before: "-Dliteral=one;two|three<four>five "},
		{name: "literal backslash", filename: `back\slash policy.properties`, quote: `"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, tc.filename)
			if err := os.WriteFile(path, []byte("# user configuration\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			options := "-Xmx1G " + tc.before + "-Djava.security.properties=" + tc.quote + path + tc.quote + " -Dafter=kept"
			t.Setenv("JAVA_TOOL_OPTIONS", options)
			t.Setenv("JAVA_SECURITY_PROPERTIES", "stale-path")
			t.Setenv("VAR", "expanded")
			jsp := prep.NewJavaSecurityProperties(bard.NewLogger(io.Discard), root)
			if err := jsp.Prepare(); err != nil {
				t.Fatal(err)
			}
			if got := os.Getenv("JAVA_SECURITY_PROPERTIES"); got != path {
				t.Errorf("JAVA_SECURITY_PROPERTIES = %q, want %q", got, path)
			}
			if got := os.Getenv("JAVA_TOOL_OPTIONS"); got != options {
				t.Errorf("JAVA_TOOL_OPTIONS = %q, want unchanged %q", got, options)
			}
			if _, err := os.Stat(filepath.Join(root, "java-security.properties")); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("default properties file must not be created; stat error = %v", err)
			}
		})
	}
}

func TestJavaSecurityProperties(t *testing.T) {
	spec.Run(t, "JavaSecurityProperties", testJavaSecurityProperties)
}

func TestJavaSecurityProperties_PrepareRejectsUnclosedQuote(t *testing.T) {
	root := t.TempDir()
	options := `-Djava.security.properties="unclosed path`
	t.Setenv("JAVA_TOOL_OPTIONS", options)
	t.Setenv("JAVA_SECURITY_PROPERTIES", "unchanged")
	jsp := prep.NewJavaSecurityProperties(bard.NewLogger(io.Discard), root)
	if err := jsp.Prepare(); err == nil || !strings.Contains(err.Error(), "JAVA_TOOL_OPTIONS") {
		t.Fatalf("want a parse error identifying JAVA_TOOL_OPTIONS, got %v", err)
	}
	if got := os.Getenv("JAVA_SECURITY_PROPERTIES"); got != "unchanged" {
		t.Errorf("JAVA_SECURITY_PROPERTIES changed on error: %q", got)
	}
	if got := os.Getenv("JAVA_TOOL_OPTIONS"); got != options {
		t.Errorf("JAVA_TOOL_OPTIONS changed on error: %q", got)
	}
	if _, err := os.Stat(filepath.Join(root, "java-security.properties")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("default properties file must not be created on parse error; stat error = %v", err)
	}
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
