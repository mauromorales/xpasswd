# xpasswd

Like [passwd](https://github.com/willdonnelly/passwd) but with a license.

Reads and enumerates system users — on Linux by parsing `/etc/passwd` directly, on macOS via `dscl`. It does not write, lock, or otherwise modify anything.

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

`NewUserList()` returns a `UserList`, which also supports pointing at a specific file (Linux only — see Supported platforms below), looking a single user up by name, and picking a free UID:

```Go
type UserList interface {
	// Get returns a user from the list by username
	Get(username string) User
	// GetAll returns every user it could successfully parse, together with
	// an error describing every line it could not.
	GetAll() ([]User, error)
	// GenerateUID returns lastUID + 1, or 0 if the list is empty
	GenerateUID() int
	// GenerateUIDInRange finds the lowest unused UID in [minimum, maximum]
	GenerateUIDInRange(minimum, maximum int) (int, error)
	LastUID() int
	// SetPath points the list at a specific file. Linux only; a no-op on
	// macOS, which always reads from the live directory service.
	SetPath(path string)
	// Load populates the list; equivalent to calling GetAll and discarding
	// the returned slice.
	Load() error
}
```

Each user implements:

```Go
type User interface {
	// UID returns the user's unique ID
	UID() (int, error)
	// GID returns the user's group ID
	GID() (int, error)
	// Username returns the user's username
	Username() string
	// Password returns the user's password (this is usually not used but we include it for completeness)
	Password() string
	// HomeDir returns the user's home directory
	HomeDir() string
	// Shell returns the user's shell
	Shell() string
	// RealName returns the user's real name
	RealName() string
}
```

## Supported platforms

Linux and macOS only. The package compiles on every `GOOS` (its interfaces carry no build constraint), but `NewUserList()` is only defined for `linux` and `darwin` — building a consumer for any other target fails with `undefined: users.NewUserList`.

BSD (`master.passwd`, a 10-field format, not the 7 fields `parseRecord` expects) would be a reasonable contribution if anyone needs it. Windows isn't a good fit for this library's model — there's no path-based passwd-equivalent, and the `User` interface's shape (integer UID, Unix home dir, shell) doesn't map onto SIDs — so it's not planned.

## Adopters

Known direct users of this library, in case you're wondering whether a change here is safe:

- [mudler/yip](https://github.com/mudler/yip) — cloud-init style provisioning tool
- [mudler/entities](https://github.com/mudler/entities) — builds passwd/group/shadow writing on top of this package's enumeration and free-UID allocation
- [harvester/yip](https://github.com/harvester/yip) — Harvester's fork of yip

Transitively, that puts this package underneath [Kairos](https://github.com/kairos-io), [Harvester](https://github.com/harvester/harvester), and [Rancher's elemental-toolkit](https://github.com/rancher/elemental-toolkit) — none of which depend on it directly, but all of which depend on `yip` or `entities`, which do.

If you're using this library somewhere not listed here, please open an issue — it helps prioritize what's safe to change.
