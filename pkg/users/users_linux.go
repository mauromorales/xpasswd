package users

import (
	"bufio"
	"fmt"
	"os"
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

func (u LinuxUser) UID() string {
	return u.uid
}

func (u LinuxUser) GID() string {
	return u.gid
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

// List returns a list of users on a Linux system
func List() ([]User, error) {
	users := make([]User, 0)

	file, err := os.Open("/etc/passwd")
	if err != nil {
		fmt.Println("Error opening the file:", err)
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
	}

	// Check if there were errors during scanning
	if err := scanner.Err(); err != nil {
		return users, fmt.Errorf("error reading the file: %w", err)
	}

	return users, nil
}
