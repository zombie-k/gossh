package main

import (
	"flag"
	"fmt"
	gossh "gossh/ssh"
)

func main() {
	flag.Parse()
	if err := gossh.Init(); err != nil {
		panic(err)
		return
	}
	client, _ := gossh.NewSshClient(gossh.Conf)
	client.Connect()
	defer client.Close()

	out := make(chan *gossh.EchoMsg, 2000)
	client.MultiRun(gossh.Conf.Cmd, out)
	for v := range out {
		if gossh.Conf.EchoIp {
			fmt.Printf("%s=:>%s", v.Addr, v.Content)
		} else {
			fmt.Printf("%s", v.Content)
		}
	}
}
