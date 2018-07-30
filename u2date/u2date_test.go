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

	Context("universal checks", func() {
		Describe("When passed a zero-length file", func() {
			It("return a zero-length file", func() {
				go writeToStdin(stdin, "")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(session.Wait().Out.Contents()).Should(Equal([]uint8("")))
			})
		})

		Describe("When passed a file with no carriage return on the last line", func() {
			It("appends a carriage return (this is a mostly-harmless bug)", func() {
				go writeToStdin(stdin, "a\nb\nc")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(session.Wait().Out.Contents()).Should(Equal([]uint8("a\nb\nc\n")))
			})
		})

		Describe("When passed a file which ends with a convertible time with no carriage return", func() {
			It("converts it & appends a carriage return (this is a mostly-harmless bug)", func() {
				go writeToStdin(stdin, "1500000000.0")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				// Dmitriy & I were coding on BOSH-on-IPv6 when 1.5 billion seconds rolled
				Ω(session.Wait().Out.Contents()).Should(Equal([]uint8("2017-07-13 19:40:00 -0700 PDT\n")))
			})
		})
	})

	Context("When it recognizes seconds since the epoch", func() {
		Describe("When passed a file containing a recognizable timestamp", func() {
			It("converts it", func() {
				go writeToStdin(stdin, "1500000000.0\n")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(session.Wait().Out.Contents()).Should(Equal([]uint8("2017-07-13 19:40:00 -0700 PDT\n")))
			})
		})

		Describe("When passed a file containing a timestamp that's 2 billion or greater", func() {
			It("doesn't convert it", func() {
				go writeToStdin(stdin, "2000000000.0")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(session.Wait().Out.Contents()).Should(Equal([]uint8("2000000000.0\n")))
			})
		})

		Describe("When passed a file containing a timestamp that's ten billion or greater", func() {
			It("doesn't converts it", func() {
				go writeToStdin(stdin, "11500000000.0\n")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(session.Wait().Out.Contents()).Should(Equal([]uint8("11500000000.0\n")))
			})
		})
		Describe("When passed a file containing a timestamp that's just shy of 2 billion", func() {
			It("converts it", func() {
				go writeToStdin(stdin, "1999999999.9")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(session.Wait().Out.Contents()).Should(Equal([]uint8("2033-05-17 20:33:19.9 -0700 PDT\n")))
			})
		})

		Describe("When passed a file containing a timestamp that's 1 billion", func() {
			It("doesn't convert it", func() {
				go writeToStdin(stdin, "1000000000.0")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(session.Wait().Out.Contents()).Should(Equal([]uint8("2001-09-08 18:46:40 -0700 PDT\n")))
			})
		})

		Describe("When passed a file containing a timestamp that's just shy of 1 billion", func() {
			It("doesn't convert it", func() {
				go writeToStdin(stdin, "999999999.9")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				// known bug: u2date will insert a "\n" when the timestamp is the very last
				Ω(session.Wait().Out.Contents()).Should(Equal([]uint8("999999999.9\n")))
			})
		})

		Describe("When passed a file containing a timestamp that doesn't have a decimal point", func() {
			It("doesn't convert it", func() {
				go writeToStdin(stdin, "1500000000.0 1500000000. 1500000000")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				// known bug: u2date will insert a "\n" when the timestamp is the very last
				Ω(session.Wait().Out.Contents()).Should(Equal([]uint8("2017-07-13 19:40:00 -0700 PDT 1500000000. 1500000000\n")))
			})
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
