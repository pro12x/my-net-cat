package pkg

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	PORT              = "8989"
	Connections       = make(chan net.Conn)
	Mutex             sync.Mutex
	Room              = 0
	Users             = make(map[net.Conn]string)
	Network           = make(chan net.Conn)
	RemoveConnections = make(chan net.Conn)
	Joins             = make(chan net.Conn)
	Messages          []string
)

func Run(tab []string) {
	// Contrôle des entrées
	if len(tab) < 2 {
		if len(tab) != 0 {
			PORT = tab[0]
		}

		if p, err := strconv.Atoi(PORT); err != nil || p < 0 {
			fmt.Println(PORT, "cannot be used as a port")
			os.Exit(0)
		}
		PORT = ":" + PORT

		// Le serveur écoute toutes les connexions entrantes sur ce PORT
		listener, err1 := net.Listen("tcp", PORT)
		if err1 != nil {
			fmt.Println("This port is not available")
			os.Exit(0)
		}
		defer func(listener net.Listener) {
			err2 := listener.Close()
			if err2 != nil {
				fmt.Println("Cannot close listening")
				os.Exit(0)
			}
		}(listener)

		fmt.Println("Server is listening on port", PORT)

		welcome, err3 := os.ReadFile("./assets/home.txt")
		if err3 != nil {
			fmt.Println("Cannot read home.txt file")
		}
		welcome = append(welcome, '\n')

		// Accepter constamment les nouvelles connexions
		// go acceptConnections(listener, Connections)
		go func() {
			for {
				// Accepte les nouvelles connexion entrantes
				connection, err4 := listener.Accept()
				if err4 != nil {
					fmt.Println("Cannot accept this connection!")
					err5 := connection.Close()
					if err5 != nil {
						return
					}
				}
				// Envoie la connexion acceptée à la chaine de connexions
				Connections <- connection
			}
		}()

		for {
			select {
			// Traitement des connexions entrantes
			case connection := <-Connections:
				Mutex.Lock()
				if Room >= 10 {
					fmt.Println("This Room Is Full!")
					connection.Write([]byte("This room is full, try again later!"))
					connection.Close()
					Mutex.Unlock()
					continue
				}
				Room++
				Mutex.Unlock()

				go func(connection net.Conn) {
					// Le message de bienvenue
					connection.Write(welcome)

					// Lecture des entrées de l'utilisateur
					buf := bufio.NewReader(connection)
					var username string

					for {
						connection.Write([]byte("[ENTER YOUR NAME]: "))
						// Lecture du nom de l'utilisateur
						userName, err6 := buf.ReadString(10)
						if err6 != nil {
							fmt.Println("Cannot read the userName")
							Mutex.Lock()
							Room--
							Mutex.Unlock()
							return
						}

						// Vérification de la validité et l'unicité du nom de l'utilisateur
						if !ValidName(connection, userName) {
							continue
						}

						// Prendre le nom sans le dernier caractère (\n)
						userName = userName[:len(userName)-1]
						// Verrouiller les autres goroutines
						Mutex.Lock()
						// Ajouter le nom de l'utilisateur dans le chaine de connexion
						Users[connection] = string(userName)
						// Déverrouiller
						Mutex.Unlock()
						// Arrêter la boucle après un nom d'utilisateur valide
						break
					}
					// TODO: read here

					// Ajoute la connexion de l'utilisateur à la chaine de réseau
					Network <- connection

					// Gérer les messages de l'utilisateur avec une nouvelle goroutine
					go Broadcasting(connection, username)
				}(connection)
			case removeConnection := <-RemoveConnections:
				// Notifier qu'un utilisateur a quitté le réseau
				notif := fmt.Sprintln(Users[removeConnection], "has left our chat...")

				for user := range Users {
					if removeConnection != user {
						// Notifier les autres utilisateurs du départ de celui-ci
						user.Write([]byte(notif))

						format := fmt.Sprintf("[%v][%v]: ", GetTime(), Users[user])
						user.Write([]byte(format))
					}
				}

				Mutex.Lock()
				Room -= 1
				delete(Users, removeConnection)
				Mutex.Unlock()

			case joiner := <-Joins:
				notif := fmt.Sprintln(Users[joiner], "has joined our chat...")
				log.Println(notif)
				for user := range Users {
					if joiner != user {
						user.Write([]byte("\n" + notif))
						msg := fmt.Sprintf("[%v][%v]: ", GetTime(), Users[user])
						user.Write([]byte(msg))
					}
				}
			}
		}
		return
	}
	fmt.Println("[USAGE]: ./TCPChat $port")
}

func GetTime() string {
	dateTime := time.Now()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", dateTime.Year(), dateTime.Month(), dateTime.Day(), dateTime.Hour(), dateTime.Minute(), dateTime.Second())
}

func acceptConnections(listener net.Listener, connections chan<- net.Conn) {
	func() {
		for {
			// Accepte les nouvelles connexion entrantes
			connection, err := listener.Accept()
			if err != nil {
				fmt.Println("Cannot accept this connection!")
				err := connection.Close()
				if err != nil {
					return
				}
			}
			// Envoie la connexion acceptée à la chaine de connexions
			connections <- connection
		}
	}()
}
