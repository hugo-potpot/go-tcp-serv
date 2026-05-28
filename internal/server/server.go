package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"tcp-serv/internal/client"
	"tcp-serv/internal/config"
	"tcp-serv/internal/protocol"
	"tcp-serv/internal/ratelimiter"
	"time"
)

const TIMEOUT = 10 * time.Second

type Server struct {
	host          string
	port          string
	ipRateLimiter *ratelimiter.IPRateLimiter
	mutex         sync.Mutex
	commands      map[string]Command
}

type Command interface {
	Handle(client *client.Client, args string) error
}

func New(config *config.Config) *Server {
	return &Server{
		host:          config.Host,
		port:          config.Port,
		ipRateLimiter: ratelimiter.NewIPRateLimiter(),
		commands: map[string]Command{
			"EHLO": &protocol.EhloCommand{},
			"DATE": &protocol.DateCommand{},
			"BYE":  &protocol.ByeCommand{},
		},
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

		newClient := client.NewClient(
			conn,
			false,
		)
		go server.handleRequest(newClient)
	}
}

func (server *Server) handleRequest(client *client.Client) {
	defer client.Conn.Close()
	defer server.ipRateLimiter.Release(client.Conn.RemoteAddr().String())

	for {
		client.Conn.SetReadDeadline(time.Now().Add(TIMEOUT))
		message, err := bufio.NewReader(client.Conn).ReadString('\n')
		client.Conn.SetReadDeadline(time.Now().Add(TIMEOUT))
		if err != nil {
			log.Println(err)
			return
		}

		if len(message) == 0 {
			continue
		}

		parts := strings.Fields(message)
		verb := strings.ToUpper(parts[0])
		args := strings.Join(parts[1:], " ")

		if cmd, ok := server.commands[verb]; ok {
			err := cmd.Handle(client, args)
			if err != nil {
				log.Printf("Error handling command %s: %s", parts[0], err)
				client.Conn.Write([]byte(err.Error() + "\r\n"))
			}
		} else {
			client.Conn.Write([]byte("500 Unknown command\r\n"))
		}
	}
}
