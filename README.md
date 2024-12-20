# xpasswd

Like [passwd](https://github.com/willdonnelly/passwd) but cross-platform and with a license.

You can list all users in a system, which would be the equivalent of reading the /etc/passwd file on Linux.

```Go
import (
	"fmt"

	"github.com/mauromorales/xpasswd/pkg/users"
)

func main() {
	// print all the users in the system
	list := users.NewUserList()
	users, err := list.GetAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, user := range users {
		fmt.Println(user.Username())
	}

}
```

Each user, Implements the following interface

```Go
type User interface {
	// UID returns the user's unique ID
	UID() string
	// GID returns the user's group ID
	GID() string
	// Username returns the user's username
	Username() string
	// HomeDir returns the user's home directory
	HomeDir() string
	// Shell returns the user's shell
	Shell() string
	// RealName returns the user's real name
	RealName() string
}
```

## TODO

- [x] Linux support
- [x] MacOS support
- [ ] BSD support 
- [ ] Windows support (I don't have a windows machine atm)
