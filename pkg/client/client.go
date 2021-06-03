package client

import (
	"bufio"
	"context"
	"fmt"
	"net"

	"github.com/go-redis-client/pkg/command"
)

type RedisClient struct {
	host string
	port int16
	conn net.Conn
}

func Connect(host string, port int16) (*RedisClient, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	redisClient := RedisClient{host: host, port: port, conn: conn}
	return &redisClient, nil
}

type ConnectionNotInitialized struct{}

func (e *ConnectionNotInitialized) Error() string {
	return "redis connection is not initialized please call connect at first"
}

func (client *RedisClient) SendCommand(command command.Command, ctx context.Context) (string, error) {
	if client.conn == nil {
		return "", &ConnectionNotInitialized{}
	}
	deadline, ok := ctx.Deadline()
	if ok {
		client.conn.SetDeadline(deadline)
	}
	n, err := fmt.Fprint(client.conn, command.Encode())
	println(n)
	if err != nil {
		return "", err
	}
	buffer := bufio.NewReader(client.conn)
	bytes, err := buffer.ReadBytes('\n')
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (client *RedisClient) Close() error {
	return client.conn.Close()
}
