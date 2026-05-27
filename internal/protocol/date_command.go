package protocol

import (
	"errors"
	"fmt"
	"tcp-serv/internal/client"
	"time"
)

type DateCommand struct{}

func (c *DateCommand) Handle(client *client.Client, args string) error {
	if !client.IsLive {
		return errors.New("550 Bad state \r\n")
	}

	date := time.Now()
	client.Conn.Write([]byte(fmt.Sprintf("250 %s \r\n", date.String())))
	return nil
}
