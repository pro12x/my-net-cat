package main

import (
	"fmt"
	"log"
	"net-cat/internal/pkg"
	"os"
)

func main() {
	str := "Janel-9"
	fmt.Println(pkg.ValidName(str))
	log.Fatalln(str)
	os.Exit(0)
}
