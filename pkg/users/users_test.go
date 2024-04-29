package users

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("users", func() {
	Describe("Get", func() {
		var list CommonUserList
		user := CommonUser{uid: "0", gid: "0", username: "root", homeDir: "/root", shell: "/bin/bash", realName: "root"}
		users := []User{user}

		Context("when the user is not present in the list", func() {
			JustBeforeEach(func() {
				list = CommonUserList{}
			})

			It("returns nil", func() {
				got := list.Get("foobar")
				Expect(got).To(BeNil())
			})
		})

		Context("when the user is present", func() {
			JustBeforeEach(func() {
				list = CommonUserList{users: users}
			})

			It("returns the user", func() {
				got := list.Get("root")
				Expect(got).To(Equal(user))
			})
		})
	})

	Describe("GenerateUID", func() {
		var list CommonUserList
		user := CommonUser{uid: "0", gid: "0", username: "root", homeDir: "/root", shell: "/bin/bash", realName: "root"}
		users := []User{user}

		Context("when the list is empty", func() {
			JustBeforeEach(func() {
				list = CommonUserList{}
			})

			It("returns 0", func() {
				got := list.GenerateUID()
				Expect(got).To(Equal(0))
			})
		})

		Context("when the list is not empty", func() {
			JustBeforeEach(func() {
				list = CommonUserList{users: users}
			})

			It("returns the next available UID", func() {
				got := list.GenerateUID()
				Expect(got).To(Equal(1))
			})
		})
	})
})
