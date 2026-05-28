package protocol

import (
	"tcp-serv/internal/client"
)

type ByeCommand struct {
	OnBye func(c *client.Client)
}

func (c *ByeCommand) Handle(client *client.Client, args string) error {
	client.Conn.Write([]byte("221 Bye"))
	if c.OnBye != nil {
		c.OnBye(client)
	}
	client.Conn.Close()
	return nil
}
