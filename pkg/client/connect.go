package client

import (
	"fmt"
	"net"
)

type RedisClient struct {
	host	string
	port	int16
	conn    *net.Conn
}

func Connect(host string, port int16) (*RedisClient, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	redisClient := RedisClient{host: host, port: port, conn: &conn}
	return &redisClient, nil
}
