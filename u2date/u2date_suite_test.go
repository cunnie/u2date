package main_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var pathToU2datetCLI string

var _ = BeforeSuite(func() {
	err := os.Setenv("TZ", "America/Los_Angeles")
	Ω(err).ShouldNot(HaveOccurred())
	pathToU2datetCLI, err = gexec.Build("github.com/cunnie/u2date/u2date")
	Ω(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func TestU2date(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "U2date Suite")
}
