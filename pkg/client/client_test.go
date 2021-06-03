package client

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis-client/pkg/command"
)

func TestConnect(t *testing.T) {
	client, err := Connect("localhost", 6379)
	if err != nil {
		t.Error("Cannot connect to redis client")
	}
	defer client.Close()

}

func TestSendCommand(t *testing.T) {
	redisClient, err := Connect("localhost", 6379)
	if err != nil {
		t.Error("Cannot connect to redis client")
	}
	defer redisClient.Close()
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	response, err := redisClient.SendCommand(command.Echo("test"), ctx)
	if err != nil {
		t.Error(err.Error())
	}
	println(response)
}
