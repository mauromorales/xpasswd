// Package users provides primitives for user information on a system.
package users

// User is an interface that represents a user on a system
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
