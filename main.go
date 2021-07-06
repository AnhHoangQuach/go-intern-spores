package main

import "github.com/AnhHoangQuach/go-intern-spores/controllers"

var server = controllers.Server{}

func main() {
	// Connect DB
	server.Initialize()

	server.Run(":8080")
}