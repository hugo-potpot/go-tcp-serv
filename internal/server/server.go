package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"tcp-serv/internal/config"
	"tcp-serv/internal/ratelimiter"
)

type Server struct {
	host          string
	port          string
	ipRateLimiter *ratelimiter.IPRateLimiter
}

//type Protocol interface {
//	handleRequest(conn net.Conn)
//}

type Client struct {
	conn net.Conn
}

func New(config *config.Config) *Server {
	return &Server{
		host:          config.Host,
		port:          config.Port,
		ipRateLimiter: ratelimiter.NewIPRateLimiter(),
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
		conn.Write([]byte("220 " + server.host + "\r\n"))
		log.Printf("New Client: %s", conn.RemoteAddr())
		// add conn to rate limiter
		if !server.ipRateLimiter.Allow(conn.RemoteAddr().String()) {
			conn.Write([]byte("429 Too Many Requests \n"))
			log.Printf("[%s] 429 Too Many Requests", conn.RemoteAddr())
			conn.Close()
			continue
		}
		client := &Client{
			conn: conn,
		}
		go server.handleRequest(client.conn)
	}
}

func (server *Server) handleRequest(conn net.Conn) {
	defer conn.Close()
	defer server.ipRateLimiter.Release(conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		log.Printf("Message incoming: %s\n", message)
		conn.Write([]byte("Message received.\n"))
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
