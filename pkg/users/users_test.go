package users

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// stubUser is a minimal User implementation for exercising CommonUserList
// in isolation. There's no exported "common" user type to build test
// fixtures from - Linux and Darwin each ship their own, and both are
// platform-gated, so a plain package-external interface implementation is
// what's left for a test file that has to compile on every platform.
type stubUser struct {
	uid, gid           int
	username, password string
	homeDir, shell     string
	realName           string
}

func (u stubUser) UID() (int, error) { return u.uid, nil }
func (u stubUser) GID() (int, error) { return u.gid, nil }
func (u stubUser) Username() string  { return u.username }
func (u stubUser) Password() string  { return u.password }
func (u stubUser) HomeDir() string   { return u.homeDir }
func (u stubUser) Shell() string     { return u.shell }
func (u stubUser) RealName() string  { return u.realName }

var _ = Describe("users", func() {
	Describe("Get", func() {
		var list CommonUserList
		rootUser := stubUser{uid: 0, gid: 0, username: "root", homeDir: "/root", shell: "/bin/bash", realName: "root"}
		barbazUser := stubUser{uid: 1000, gid: 1000, username: "barbaz", homeDir: "/home/barbaz", shell: "/bin/bash", realName: "Bar Baz"}
		users := []User{rootUser, barbazUser}

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
				Expect(got).To(Equal(rootUser))
			})
		})
	})

	Describe("GenerateUID", func() {
		var list CommonUserList
		user := stubUser{uid: 0, gid: 0, username: "root", homeDir: "/root", shell: "/bin/bash", realName: "root"}
		foobar := stubUser{uid: 1000, gid: 1000, username: "foobar", homeDir: "/home/foobar", shell: "/bin/bash", realName: "foo bar"}
		users := []User{user, foobar}

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
				list = CommonUserList{users: users, lastUID: 1000}
			})

			It("returns the next available UID", func() {
				got := list.GenerateUID()
				Expect(got).To(Equal(1001))
			})
		})
	})
})
