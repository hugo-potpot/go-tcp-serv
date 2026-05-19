package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"tcp-serv/internal/config"
)

type Server struct {
	host string
	port string
}

//type Protocol interface {
//	handleRequest(conn net.Conn)
//}

type Client struct {
	conn net.Conn
}

func New(config *config.Config) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn: conn,
		}
		go handleRequest(client.conn)
	}
}

func handleRequest(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			return
		}
		fmt.Printf("Message incoming: %s", string(message))
		conn.Write([]byte("Message received.\n"))
	}
}
