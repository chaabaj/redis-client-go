package command

import "testing"

func TestEcho(t *testing.T) {
	expectedRawCommand := "ECHO \"test\""
	echoCmd := Echo("test")
	if echoCmd.Encode() != expectedRawCommand {
		t.Errorf("Expected %s, got: %s", expectedRawCommand, echoCmd)
	}
}
