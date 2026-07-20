package calc_test

import (
	"bytes"
	"crypto/sha256"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/softleader/memory-calculator/calc"
)

const (
	certificateChildEnv = "MEMORY_CALCULATOR_CERTIFICATE_CHILD"
	certificateRootEnv  = "MEMORY_CALCULATOR_CERTIFICATE_ROOT"
)

func TestCalculator_Execute_CertificateLoaderContract(t *testing.T) {
	if os.Getenv(certificateChildEnv) != "" {
		runCertificateLoaderChild(t)
		return
	}

	root := t.TempDir()
	truststore := filepath.Join(root, "truststore.p12")
	contents, err := os.ReadFile(filepath.Join("testdata", "truststore.p12"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(truststore, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "certificates"), 0o755); err != nil {
		t.Fatal(err)
	}
	certificate, err := filepath.Abs(filepath.Join("testdata", "certificate.pem"))
	if err != nil {
		t.Fatal(err)
	}
	certificateContents, err := os.ReadFile(certificate)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(certificateContents, []byte("BEGIN CERTIFICATE")) || bytes.Contains(certificateContents, []byte("PRIVATE KEY")) {
		t.Fatal("certificate fixture must contain a public certificate and no private key")
	}

	before := sha256.Sum256(contents)
	cmd := exec.Command(os.Args[0], "-test.run=^TestCalculator_Execute_CertificateLoaderContract$")
	cmd.Env = []string{
		certificateChildEnv + "=1",
		certificateRootEnv + "=" + root,
		"BPI_JVM_CACERTS=" + truststore,
		"HOME=" + root,
		"JAVA_TOOL_OPTIONS=-XX:ActiveProcessorCount=2",
		"SSL_CERT_FILE=" + certificate,
		"SSL_CERT_DIR=" + filepath.Join(root, "certificates"),
	}
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("certificate loader child: %v\n%s", err, output)
	}
	after, err := os.ReadFile(truststore)
	if err != nil {
		t.Fatal(err)
	}
	if before == sha256.Sum256(after) {
		t.Fatal("truststore digest did not change")
	}
}

func runCertificateLoaderChild(t *testing.T) {
	t.Helper()
	root := os.Getenv(certificateRootEnv)
	logger := bard.NewLogger(io.Discard)
	calculator := calc.NewCalculator(logger)
	calculator.MemoryLimitPath.V1 = filepath.Join(root, "memory.limit")
	calculator.MemoryLimitPath.V2 = filepath.Join(root, "missing-memory-limit")
	if err := os.WriteFile(calculator.MemoryLimitPath.V1, []byte("10g"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := calculator.LoadedClassCount.Set("42"); err != nil {
		t.Fatal(err)
	}
	if err := calculator.ThreadCount.Set("10"); err != nil {
		t.Fatal(err)
	}
	if err := calculator.EnableJdwp.Set("false"); err != nil {
		t.Fatal(err)
	}
	if _, err := calculator.Execute(); err != nil {
		t.Fatal(err)
	}
}
