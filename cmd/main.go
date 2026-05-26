package main

import (
	"tcp-serv/internal/config"
	"tcp-serv/internal/server"
)

func main() {
	serv := server.New(
		&config.Config{
			Host: "localhost",
			Port: "8080",
		})
	serv.Run()
}
