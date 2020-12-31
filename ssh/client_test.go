package ssh

import (
	"flag"
	"testing"
)

func TestClient_MultiRun(t *testing.T) {
	flag.Parse()
	_ = Init()
	client, _ := NewSshClient(Conf)
	client.Connect()
	defer client.Close()

	cmd := "pwd"

	out := make(chan *EchoMsg, 2000)
	client.MultiRun(cmd, out)
	for v := range out {
		t.Log(v.Addr, v.Content)
	}
}
