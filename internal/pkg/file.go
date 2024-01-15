package pkg

import (
	"log"
	"net"
	"os"
)

func WriteToFile(text string) {
	file, err := os.OpenFile("logs/log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = file.WriteString(text)
	if err != nil {
		return
	}
}

func ReadFromFile(connection net.Conn) {
	file, err := os.OpenFile("logs/log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	buffer := make([]byte, 4096)
	for {
		_, err := file.Read(buffer)
		if err != nil {
			break
		}
		_, err = connection.Write(buffer)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if err := file.Close(); err != nil {
		log.Fatalln(err)
	}
}
