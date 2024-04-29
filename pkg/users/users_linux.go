package users

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type LinuxUser struct {
	login         string
	uid           string
	gid           string
	nameOrComment string
	home          string
	shell         string
	interpreter   string
}

func NewUserList() UserList {
	return &LinuxUserList{path: "/etc/passwd"}
}

// LinuxUserList is a list of Linux users
type LinuxUserList struct {
	CommonUserList
	path string
}

func (u LinuxUser) UID() (int, error) {
	return strconv.Atoi(u.uid)
}

func (u LinuxUser) GID() (int, error) {
	return strconv.Atoi(u.gid)
}

func (u LinuxUser) Username() string {
	return u.login
}

func (u LinuxUser) HomeDir() string {
	return u.home
}

func (u LinuxUser) Shell() string {
	return u.shell
}

func (u LinuxUser) RealName() string {
	return u.nameOrComment
}

func (l *LinuxUserList) SetPath(path string) {
	l.path = path
}

func (l *LinuxUserList) Load() error {
	_, err := l.GetAll()
	return err
}

// GetAll returns all users in the list
func (l *LinuxUserList) GetAll() ([]User, error) {
	users := make([]User, 0)

	file, err := os.Open(l.path)
	if err != nil {
		return users, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read each line
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into parts using ':' as the delimiter
		parts := strings.Split(line, ":")

		// Check if the line is correctly formatted with 7 fields
		if len(parts) != 7 {
			return users, fmt.Errorf("unexpected format: %s", line)
		}

		user := LinuxUser{
			login:         parts[0],
			uid:           parts[2],
			gid:           parts[3],
			nameOrComment: parts[4],
			home:          parts[5],
			shell:         parts[6],
			interpreter:   parts[6],
		}
		users = append(users, user)

		uid, err := user.UID()
		if err != nil {
			return users, fmt.Errorf("failed to convert UID to int: %w", err)
		}

		if uid > l.lastUID {
			l.lastUID = uid
		}
	}

	// Check if there were errors during scanning
	if err := scanner.Err(); err != nil {
		return users, fmt.Errorf("error reading the file: %w", err)
	}

	l.users = users

	return users, nil
}
