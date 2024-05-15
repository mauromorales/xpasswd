package users

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListDarwing", func() {
	It("parse record", func() {
		rootRecord := `NFSHomeDirectory: /var/root
		Password: *
		PrimaryGroupID: 0
		RealName:
		 System Administrator
		RecordName:
		 root
		 BUILTIN\Local System
		UniqueID: 0
		UserShell: /bin/sh`

		got := parseRecord(rootRecord)
		Expect(got.UID()).To(Equal(0))
		Expect(got.GID()).To(Equal(0))
		Expect(got.HomeDir()).To(Equal("/var/root"))
		Expect(got.Shell()).To(Equal("/bin/sh"))
		Expect(got.Username()).To(Equal("root"))
		Expect(got.RealName()).To(Equal("System Administrator"))
		Expect(got.Password()).To(Equal("*"))
	})
})

var _ = Describe("DarwinUser", func() {
	Describe("Get", func() {
		var list DarwinUserList
		rootUser := DarwinUser{uniqueID: "0", primaryGroupID: "0", recordName: "root", password: "*", nFSHomeDirectory: "/root", userShell: "/bin/bash", realName: "root"}
		barbazUser := DarwinUser{uniqueID: "1000", primaryGroupID: "1000", recordName: "barbaz", password: "*", nFSHomeDirectory: "/home/barbaz", userShell: "/bin/bash", realName: "Bar Baz"}
		users := []User{rootUser, barbazUser}

		Context("when the user is not present in the list", func() {
			JustBeforeEach(func() {
				list = DarwinUserList{}
			})

			It("returns nil", func() {
				got := list.Get("foobar")
				Expect(got).To(BeNil())
			})
		})

		Context("when the user is present", func() {
			JustBeforeEach(func() {
				list = DarwinUserList{
					CommonUserList{users: users, lastUID: 1000},
				}
			})

			It("returns the user", func() {
				got := list.Get("root")
				Expect(got).To(Equal(rootUser))
				Expect(list.GenerateUID()).To(Equal(1001))
			})
		})
	})
})
