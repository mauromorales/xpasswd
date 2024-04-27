package xpasswd_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestXpasswd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Xpasswd Suite")
}
