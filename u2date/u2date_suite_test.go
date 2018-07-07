package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestU2date(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "U2date Suite")
}
