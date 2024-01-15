package pkg

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

func IsPrintable(s string) bool {
	reg := regexp.MustCompile(`^[ -~]+$`)
	return reg.MatchString(s)
}

func ValidMessage(message string) bool {
	message = strings.TrimSuffix(message, "\n")

	if len(message) == 0 || len(strings.TrimSpace(message)) == 0 {
		fmt.Println("Message cannot be empty!, write something")
		return false
	}
	if !IsPrintable(message) {
		fmt.Println("Your message is not valid")
		return false
	}
	return true
}

func ValidName(connection net.Conn, username string) bool {
	username = strings.TrimSuffix(username, "\n")

	// Vérification du nom entré par l'utilisateur
	if len(username) == 0 || len(strings.TrimSpace(username)) == 0 {
		fmt.Println("This name is not valid")
		// connection.Write([]byte("The username cannot be empty!\n"))
		_, err := connection.Write([]byte("The username cannot be empty!\n"))
		if err != nil {
			return false
		}
		return false
	}

	// Vérification du nom dans la liste des utilisateurs existants
	for _, user := range Users {
		if user == username {
			fmt.Println("User Already Exists")
			//connection.Write([]byte("This username already exists, try another.\n"))
			_, err := connection.Write([]byte("This username already exists, try another.\n"))
			if err != nil {
				return false
			}
			return false
		}
	}
	return true
}
