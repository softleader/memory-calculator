package boot_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	boot "github.com/softleader/memory-calculator/boot/helper"

	"github.com/paketo-buildpacks/libpak/bard"
	springboot "github.com/paketo-buildpacks/spring-boot/v5/boot"
)

func testWebApplicationType(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = gomega.NewWithT(t).Expect

		buffer *bytes.Buffer
		wat    boot.WebApplicationType
	)

	it.Before(func() {
		buffer = &bytes.Buffer{}
		logger := bard.NewLogger(buffer)
		wat = boot.WebApplicationType{Logger: logger}
	})

	context("when web application type is None", func() {
		it.Before(func() {
			boot.ResolveWebAppType = func() (springboot.ApplicationType, error) {
				return springboot.None, nil
			}
		})

		it("returns the correct thread count and logs", func() {
			result, err := wat.Execute()
			Expect(err).NotTo(gomega.HaveOccurred())
			Expect(result).To(gomega.Equal(map[string]string{"BPL_JVM_THREAD_COUNT": "50"}))
			Expect(buffer.String()).To(gomega.ContainSubstring("Non-web application detected"))
		})
	})

	context("when web application type is Reactive", func() {
		it.Before(func() {
			boot.ResolveWebAppType = func() (springboot.ApplicationType, error) {
				return springboot.Reactive, nil
			}
		})

		it("returns the correct thread count and logs", func() {
			result, err := wat.Execute()
			Expect(err).NotTo(gomega.HaveOccurred())
			Expect(result).To(gomega.Equal(map[string]string{"BPL_JVM_THREAD_COUNT": "50"}))
			Expect(buffer.String()).To(gomega.ContainSubstring("Reactive web application detected"))
		})
	})

	context("when web application type is Servlet", func() {
		it.Before(func() {
			boot.ResolveWebAppType = func() (springboot.ApplicationType, error) {
				return springboot.Servlet, nil
			}
		})

		it("returns the correct thread count and logs", func() {
			result, err := wat.Execute()
			Expect(err).NotTo(gomega.HaveOccurred())
			Expect(result).To(gomega.Equal(map[string]string{"BPL_JVM_THREAD_COUNT": "250"}))
			Expect(buffer.String()).To(gomega.ContainSubstring("Servlet web application detected"))
		})
	})

	context("when resolving web application type fails", func() {
		it.Before(func() {
			boot.ResolveWebAppType = func() (springboot.ApplicationType, error) {
				return springboot.None, errors.New("test-error")
			}
		})

		it("returns an error", func() {
			_, err := wat.Execute()
			Expect(err).To(gomega.MatchError("test-error"))
		})
	})
}
