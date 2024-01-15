package main

import (
	"net-cat/internal/pkg"
	"os"
)

// Main function to start the chat server
func main() {
	pkg.Run(os.Args[1:])
}
