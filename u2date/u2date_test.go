package main_test

import (
	"io"
	"log"
	"os/exec"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("u2date", func() {
	var stdin io.WriteCloser
	var cmd *exec.Cmd
	var args []string

	if runtime.GOOS == "windows" {
		panic("Tests won't work under Windows; they depend on the TZ environment variable")
	}

	JustBeforeEach(func() {
		cmd = exec.Command(pathToU2datetCLI, args...)
		var err error
		stdin, err = cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
	})

	When("there are no timestamps to convert", func() {
		When("passed a zero-length file", func() {
			It("return a zero-length file", func() {
				go writeToStdin(stdin, "")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(string(session.Wait().Out.Contents())).Should(Equal(""))
			})
		})

		When("passed a file with no carriage return on the last line", func() {
			It("appends a carriage return (this is a mostly-harmless bug)", func() {
				go writeToStdin(stdin, "a\nb\nc")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(string(session.Wait().Out.Contents())).Should(Equal("a\nb\nc\n"))
			})
		})

	})

	When("it recognizes seconds since the epoch", func() {
		When("passed a file which ends with a convertible time with no carriage return", func() {
			When("the time has a decimal point", func() {
				It("converts it & appends a carriage return (this is a mostly-harmless bug)", func() {
					go writeToStdin(stdin, "1500000000.0")
					session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
					Ω(err).ShouldNot(HaveOccurred())
					// Dmitriy & I were coding on BOSH-on-IPv6 when 1.5 billion seconds rolled
					Ω(string(session.Wait().Out.Contents())).Should(Equal("2017-07-13 19:40:00 -0700 PDT\n"))
				})
			})
			When("the time has no decimal point", func() {
				It("converts it & appends a carriage return (this is a mostly-harmless bug)", func() {
					go writeToStdin(stdin, "1500000000")
					session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(string(session.Wait().Out.Contents())).Should(Equal("2017-07-13 19:40:00 -0700 PDT\n"))
				})
			})
		})
		When("there are several timestamps", func() {
			It("converts them all", func() {
				go writeToStdin(stdin, "💜1500000000❤️1500000001.🧡1500000002.0💛\n")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(string(session.Wait().Out.Contents())).Should(Equal("💜2017-07-13 19:40:00 -0700 PDT❤️2017-07-13 19:40:01 -0700 PDT.🧡2017-07-13 19:40:02 -0700 PDT💛\n"))
			})
		})
		When("passed a file containing a timestamp that's just shy of 2 billion", func() {
			It("converts it", func() {
				go writeToStdin(stdin, "1999999999.9")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(string(session.Wait().Out.Contents())).Should(Equal("2033-05-17 20:33:19.9 -0700 PDT\n"))
			})
		})
		When("passed a file containing a timestamp that's 1 billion", func() {
			It("converts it", func() {
				go writeToStdin(stdin, "💛1000000000.0💚1000000000.💙1000000000💜")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(string(session.Wait().Out.Contents())).Should(Equal("💛2001-09-08 18:46:40 -0700 PDT💚2001-09-08 18:46:40 -0700 PDT.💙2001-09-08 18:46:40 -0700 PDT💜\n"))
			})
		})
	})

	When("it shouldn't convert a number", func() {
		When("passed a file containing a timestamp that's 2 billion or greater", func() {
			It("doesn't convert it", func() {
				go writeToStdin(stdin, "2000000000.0 2000000000")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(string(session.Wait().Out.Contents())).Should(Equal("2000000000.0 2000000000\n"))
			})
		})
		When("passed a file containing a timestamp that's ten billion or greater", func() {
			It("doesn't convert it", func() {
				go writeToStdin(stdin, "11500000000.0 11500000000\n")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(string(session.Wait().Out.Contents())).Should(Equal("11500000000.0 11500000000\n"))
			})
		})

		When("passed a file containing a timestamp that's just shy of 1 billion", func() {
			It("doesn't convert it", func() {
				go writeToStdin(stdin, "999999999.9")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				// known bug: u2date will insert a "\n" when the timestamp is the very last
				Ω(string(session.Wait().Out.Contents())).Should(Equal("999999999.9\n"))
			})
		})
	})

	When("it recognizes nanoseconds since the epoch", func() {
		When("passed a file containing a recognizable timestamp", func() {
			It("converts it", func() {
				go writeToStdin(stdin, "1500000000000000000\n")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(string(session.Wait().Out.Contents())).Should(Equal("2017-07-13 19:40:00 -0700 PDT\n"))
			})
		})

		When("passed a file containing a timestamp that's 2 quadrillion or greater", func() {
			It("doesn't convert it", func() {
				go writeToStdin(stdin, "2000000000000000000")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(string(session.Wait().Out.Contents())).Should(Equal("2000000000000000000\n"))
			})
		})

		When("passed a file containing a timestamp that's just shy of 2 quadrillion", func() {
			It("converts it", func() {
				go writeToStdin(stdin, "1999999999999999999")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(string(session.Wait().Out.Contents())).Should(Equal("2033-05-17 20:33:19.999999999 -0700 PDT\n"))
			})
		})

		When("passed a file containing a timestamp that's 1 quadrillion", func() {
			It("doesn't convert it", func() {
				go writeToStdin(stdin, "1000000000000000000")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(string(session.Wait().Out.Contents())).Should(Equal("2001-09-08 18:46:40 -0700 PDT\n"))
			})
		})

		When("passed a file containing a timestamp that's just shy of 1 quadrillion", func() {
			It("doesn't convert it", func() {
				go writeToStdin(stdin, "999999999999999")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				// known bug: u2date will insert a "\n" when the timestamp is the very last
				Ω(string(session.Wait().Out.Contents())).Should(Equal("999999999999999\n"))
			})
		})

		When("passed a file containing a timestamp that's greater than 1 quadrillion and has decimal points", func() {
			It("converts it and leaves the decimal points unmolested", func() {
				// go writeToStdin(stdin, "150000000000000000.0 1500000000000000000. 1500000000000000000")
				go writeToStdin(stdin, "1500000000000000000.0 1500000000000000000. 1500000000000000000")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				// known bug: u2date will insert a "\n" when the timestamp is the very last
				Ω(string(session.Wait().Out.Contents())).Should(Equal("2017-07-13 19:40:00 -0700 PDT.0 2017-07-13 19:40:00 -0700 PDT. 2017-07-13 19:40:00 -0700 PDT\n"))
			})
		})
		When("using the '-wrap=\"' flag", func() {
			BeforeEach(func() {
				args = []string{"-wrap=\""}
			})
			It("blindly puts double-quotes around the date", func() {
				go writeToStdin(stdin, `1500000000000000000💜1500000000❤️1500000001.🧡1500000002.0💛`)
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(string(session.Wait().Out.Contents())).Should(Equal("\"2017-07-13 19:40:00 -0700 PDT\"💜\"2017-07-13 19:40:00 -0700 PDT\"❤️\"2017-07-13 19:40:01 -0700 PDT\".🧡\"2017-07-13 19:40:02 -0700 PDT\"💛\n"))
			})
		})
	})
})

func writeToStdin(stdin io.WriteCloser, stdinString string) {
	defer func() { _ = stdin.Close() }()
	_, err := io.WriteString(stdin, stdinString)
	if err != nil {
		log.Fatal(err)
	}
}
