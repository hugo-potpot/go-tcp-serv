package protocol

import (
	"errors"
	"fmt"
	"tcp-serv/internal/client"
)

type EhloCommand struct{}

func (c *EhloCommand) Handle(client *client.Client, args string) error {
	if len(args) < 2 {
		return errors.New("500 You need to entry your name (EHLO <name>)\n")
	}
	client.IsLive = true
	client.Conn.Write([]byte(fmt.Sprintf("250 Pleased to meet you %s \r\n", args)))
	return nil
}
