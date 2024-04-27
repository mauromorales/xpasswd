package users

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListDarwing", func() {
	It("parse record", func() {
		rootRecord := `NFSHomeDirectory: /var/root
		PrimaryGroupID: 0
		RealName:
		 System Administrator
		RecordName:
		 root
		 BUILTIN\Local System
		UniqueID: 0
		UserShell: /bin/sh`

		got := parseRecord(rootRecord)
		Expect(got.UID()).To(Equal("0"))
		Expect(got.GID()).To(Equal("0"))
		Expect(got.HomeDir()).To(Equal("/var/root"))
		Expect(got.Shell()).To(Equal("/bin/sh"))
		Expect(got.Username()).To(Equal("root"))
		Expect(got.RealName()).To(Equal("System Administrator"))
	})
})
