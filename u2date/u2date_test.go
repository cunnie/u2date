package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
)

var _ = Describe("U2date", func() {
	var pathToU2datetCLI string
	var stdin io.WriteCloser
	var err error
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		panic("Tests won't work under Windows; they depend on the TZ environment variable")
	}

	BeforeSuite(func() {
		os.Setenv("TZ", "America/Los_Angeles")
		pathToU2datetCLI, err = gexec.Build("github.com/cunnie/u2date/u2date")
		Ω(err).ShouldNot(HaveOccurred())
	})

	BeforeEach(func() {
		cmd = exec.Command(pathToU2datetCLI)
		stdin, err = cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
	})

	Describe("When passed a zero-length file", func() {
		It("return a zero-length file", func() {
			go writeToStdin(stdin, "")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(session.Wait().Out.Contents()).Should(Equal([]uint8("")))
		})
	})

	Describe("When passed a file with no carriage return on the last line", func() {
		It("appends a carriage return (this is a mostly-harmless but)", func() {
			go writeToStdin(stdin, "a")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(session.Wait().Out.Contents()).Should(Equal([]uint8("a\n")))
		})
	})

	Describe("When passed a file containing a recognizable timestamp", func() {
		It("converts it", func() {
			go writeToStdin(stdin, "1500000000.0\n")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())
			// known bug: u2date will insert a "\n" when the timestamp is the very last
			Ω(session.Wait().Out.Contents()).Should(Equal([]uint8("2017-07-13 19:40:00 -0700 PDT\n")))
		})
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})
})

func writeToStdin(stdin io.WriteCloser, stdinString string) {
	defer stdin.Close()
	_, err := io.WriteString(stdin, stdinString)
	if err != nil {
		log.Fatal(err)
	}
}
