package pkg

import (
	"bufio"
	"fmt"
	"net"
)

func Broadcasting(connection net.Conn, username string) {
	// Fermer la connexion à la fin de l'exécution de la fonction
	// defer connection.Close()
	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil {
			fmt.Println("Error closing connection")
		}
	}(connection)

	// Lecture des messages entrants sans interruption
	for {
		format1 := fmt.Sprintf("[%v][%v]: ", GetTime(), username)
		// connection.Write([]byte(format))
		_, err2 := connection.Write([]byte(format1))
		if err2 != nil {
			return
		}

		// Lecture des messages entrants de l'utilisateur
		buf := bufio.NewReader(connection)
		content, err := buf.ReadString(10)
		if err != nil {
			fmt.Println(username, "has left our chat...")
			break
		}

		if !ValidMessage(content) {
			continue
		}

		content = content[:len(content)-1]
		format2 := fmt.Sprintf("\n[%v][%v]: [%v]", GetTime(), username, content)
		format3 := fmt.Sprintf("[%v][%v]: [%v]", GetTime(), username, content)
		fmt.Println(format3)

		// TODO: Save here
		// SaveMessage(format3)

		for user := range Users {
			if connection != user {
				// user.Write([]byte(format))
				_, err3 := user.Write([]byte(format2))
				if err3 != nil {
					return
				}
				format4 := fmt.Sprintf("[%v][%v]: ", GetTime(), Users[user])
				// user.Write([]byte(format))
				_, err4 := user.Write([]byte(format4))
				if err4 != nil {
					return
				}
			}
		}
	}
	RemoveConnections <- connection
}
