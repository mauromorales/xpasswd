// Package users provides primitives for user information on a system.
package users

import (
	"errors"
	"fmt"
	"strconv"
)

// User is an interface that represents a user on a system
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

type CommonUser struct {
	uid      string
	gid      string
	username string
	password string
	homeDir  string
	shell    string
	realName string
}

func (u CommonUser) UID() (int, error) {
	return strconv.Atoi(u.uid)
}

func (u CommonUser) GID() (int, error) {
	return strconv.Atoi(u.gid)
}

func (u CommonUser) Username() string {
	return u.username
}

func (u CommonUser) Password() string {
	return u.password
}

func (u CommonUser) HomeDir() string {
	return u.homeDir
}

func (u CommonUser) Shell() string {
	return u.shell
}

func (u CommonUser) RealName() string {
	return u.realName
}

// UserList is an interface that represents a list of users
type UserList interface {
	// Get returns a user from the list by username
	Get(username string) User
	// GetAll returns every user it could successfully parse, together with
	// a joined error (see errors.Join) describing every line it could not.
	// The returned slice can be non-empty even when the error is non-nil:
	// callers that want to treat any malformed line as fatal should check
	// the error, and callers that want best-effort results should use the
	// slice and log or ignore the error.
	GetAll() ([]User, error)
	// GenerateUID returns lastUID + 1, or -1 if the list is empty.
	GenerateUID() int
	GenerateUIDInRange(int, int) (int, error)
	LastUID() int
	SetPath(path string)
	Load() error
}

// CommonUserList is a common implementation of UserList
type CommonUserList struct {
	users   []User
	lastUID int
}

// Get checks if a user with the given username exists
func (list CommonUserList) Get(username string) User {
	for _, user := range list.users {
		if user.Username() == username {
			return user
		}
	}
	return nil
}

func (list CommonUserList) LastUID() int {
	return list.lastUID
}

// GenerateUID returns lastUID + 1, or -1 if the list is empty (or hasn't
// been loaded). -1 is never a valid UID, unlike 0, which is root's.
func (list CommonUserList) GenerateUID() int {
	if len(list.users) == 0 {
		return -1
	}
	return list.lastUID + 1
}

// Finds the lowest available uid in the specified range.
// Returns an error if there is no available uid in that range.
func (list CommonUserList) GenerateUIDInRange(minimum, maximum int) (int, error) {
	userSet := make(map[int]struct{})
	for _, user := range list.users {
		uid, err := user.UID()
		if err != nil {
			return -1, fmt.Errorf("getting user's uid: %w", err)
		}
		userSet[uid] = struct{}{}
	}

	result := -1
	for i := minimum; i <= maximum; i++ {
		if _, found := userSet[i]; found {
			continue // uid in use, skip it
		}
		result = i // found a free one, stop here
		break
	}

	if result == -1 {
		return result, errors.New("no available uid in range")
	}

	return result, nil
}
