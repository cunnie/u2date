package u2date_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/cunnie/u2date/u2date"
)

var _ = Describe("U2date", func() {
	It("should return 'aha'", func() {
		Expect(U2date()).To(Equal("aha"))
	})
})
