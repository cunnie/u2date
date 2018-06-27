package u2date_test

import (
	. "github.com/cunnie/u2date/u2date"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("U2date", func() {
	stdinString := "Hello!"
	stdin := strings.NewReader(stdinString)
	It("should return 'aha'", func() {
		Expect((U2date(stdin))).To(Equal(stdinString))
	})
})
