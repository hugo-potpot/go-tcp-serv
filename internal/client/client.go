package client

import "net"

type Client struct {
	Conn   net.Conn
	IsLive bool
}

func NewClient(conn net.Conn, isLive bool) *Client {
	return &Client{Conn: conn, IsLive: isLive}
}
