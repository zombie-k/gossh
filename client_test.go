package gossh

import (
	"flag"
	"testing"
)

func TestClient_MultiRun(t *testing.T) {
	flag.Parse()
	if err := Init(); err != nil {
		panic(err)
		return
	}
	client, _ := NewSshClient(Conf)
	client.Connect()
	defer client.Close()

	cmd := "sudo docker ps"

	out := make(chan *EchoMsg, 2000)
	client.MultiRun(cmd, out)
	for v := range out {
		t.Log(v.Addr, v.Content)
	}
}
