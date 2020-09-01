package client

import "testing"

func TestConnect(t *testing.T) {
	_, err := Connect("localhost", 6379)
	if err != nil {
		t.Error("Cannot connect to redis client")
	}
}
