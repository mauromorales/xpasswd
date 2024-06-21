package users

import (
	"log"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LinuxUserList", func() {
	Describe("GenerateUIDInRange", func() {
		var file *os.File
		var err error
		var list LinuxUserList

		BeforeEach(func() {
			file, err = os.CreateTemp("", "passwd")
			Expect(err).ToNot(HaveOccurred())
			DeferCleanup(func() {
				defer os.Remove(file.Name())
			})

			_, err = file.WriteString("root:x:0:0:root:/root:/bin/bash\n")
			Expect(err).ToNot(HaveOccurred())
			_, err = file.WriteString("foo:x:1000:1000:foo:/home/foo:/bin/bash\n")
			Expect(err).ToNot(HaveOccurred())
			_, err = file.WriteString("foo:x:1001:1000:foo:/home/foo:/bin/bash\n")
			_, err = file.WriteString("foo:x:1001:1000:foo:/home/foo:/bin/bash\n")
			Expect(err).ToNot(HaveOccurred())

			list = LinuxUserList{}
			list.SetPath(file.Name())
			Expect(list.Load()).ToNot(HaveOccurred())
		})

		When("a uid is available in the range", func() {
			var minimum, maximum int
			BeforeEach(func() {
				minimum = 1000
				maximum = 2000
			})

			It("returns the minimum available uid", func() {
				r, err := list.GenerateUIDInRange(minimum, maximum)
				Expect(err).ToNot(HaveOccurred())
				Expect(r).To(Equal(1002))
			})
		})

		When("there is no available uid", func() {
			var minimum, maximum int
			BeforeEach(func() {
				minimum = 1000
				maximum = 1001
			})

			It("returns an error", func() {
				_, err := list.GenerateUIDInRange(minimum, maximum)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("no available uid in range"))
			})
		})
	})

	It("parses a record", func() {
		rootRecord := `root:x:0:0:root:/root:/bin/bash`

		got, err := parseRecord(rootRecord)
		Expect(err).To(BeNil())

		Expect(got.UID()).To(Equal(0))
		Expect(got.GID()).To(Equal(0))
		Expect(got.HomeDir()).To(Equal("/root"))
		Expect(got.Shell()).To(Equal("/bin/bash"))
		Expect(got.Username()).To(Equal("root"))
		Expect(got.RealName()).To(Equal("root"))
		Expect(got.Password()).To(Equal("x"))
	})

	It("Gets all users", func() {
		file, err := os.CreateTemp("", "passwd")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(file.Name())

		_, err = file.WriteString("root:x:0:0:root:/root:/bin/bash\n")
		Expect(err).ToNot(HaveOccurred())
		_, err = file.WriteString("foo:x:1000:1000:foo:/home/foo:/bin/bash\n")
		Expect(err).ToNot(HaveOccurred())

		list := LinuxUserList{}
		list.SetPath(file.Name())
		err = list.Load()
		Expect(err).ToNot(HaveOccurred())

		user := list.Get("root")
		Expect(user).ToNot(BeNil())
		user = list.Get("foo")
		Expect(user).ToNot(BeNil())
		user = list.Get("bar")
		Expect(user).To(BeNil())
		Expect(list.GenerateUID()).To(Equal(1001))
	})

	It("Returns users even when there might be errors reading", func() {
		file, err := os.CreateTemp("", "passwd")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(file.Name())

		_, err = file.WriteString("root:x:0:0:root:/root:/bin/bash\n")
		Expect(err).ToNot(HaveOccurred())
		_, err = file.WriteString("foo:x:1000::foo:/home/foo:/bin/bash\n")
		Expect(err).ToNot(HaveOccurred())
		_, err = file.WriteString("baz:x:1000:1000:foo:/home/foo:/bin/bash:asdf\n")
		Expect(err).ToNot(HaveOccurred())

		list := LinuxUserList{}
		list.SetPath(file.Name())
		users, err := list.GetAll()
		Expect(err).To(HaveOccurred())
		Expect(users).To(HaveLen(2))
		foo := list.Get("foo")
		uid, err := foo.GID()
		Expect(uid).To(Equal(0))
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("DarwinUser", func() {
	Describe("Get", func() {
		var list LinuxUserList
		rootUser := LinuxUser{uid: "0", gid: "0", login: "root", password: "*", userHomeDir: "/root", usercommandInterpreter: "/bin/bash", userNameOrComment: "root"}
		barbazUser := LinuxUser{uid: "1000", gid: "1000", login: "barbaz", password: "*", userHomeDir: "/home/barbaz", usercommandInterpreter: "/bin/bash", userNameOrComment: "Bar Baz"}
		users := []User{rootUser, barbazUser}

		Context("when the user is not present in the list", func() {
			JustBeforeEach(func() {
				list = LinuxUserList{}
			})

			It("returns nil", func() {
				got := list.Get("foobar")
				Expect(got).To(BeNil())
			})
		})

		Context("when the user is present", func() {
			JustBeforeEach(func() {
				list = LinuxUserList{
					CommonUserList: CommonUserList{users: users, lastUID: 1000},
					path:           "/etc/passwd",
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
