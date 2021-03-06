package gossh

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	c       *Config
	SshConf *ssh.ClientConfig
	Infos   []*Info
	Modes   *ssh.TerminalModes
	Ch      chan *EchoMsg
	wg      *sync.WaitGroup
}

type Info struct {
	Addr    string
	Port    int
	Client  *ssh.Client
	Session *ssh.Session
	Cmd     string
}

type EchoMsg struct {
	Addr    string
	Content string
	err     error
}

func NewSshClient(conf *Config) (*Client, error) {
	client := &Client{
		c:     conf,
		Infos: make([]*Info, 0),
		Ch:    make(chan *EchoMsg, 1000),
		wg:    new(sync.WaitGroup),
	}
	sshConfig := &ssh.ClientConfig{
		Config:          ssh.Config{},
		User:            conf.User,
		Auth:            make([]ssh.AuthMethod, 0),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(conf.Timeout) * time.Second,
	}
	sshConfig.SetDefaults()
	if conf.PrivateKey == "" && conf.IdRsa == "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.Password(conf.Password))
	} else {
		var pemBytes []byte
		var err error
		if conf.IdRsa != "" {
			pemBytes = []byte(conf.IdRsa)
		} else {
			pemBytes, err = ioutil.ReadFile(conf.PrivateKey)
			if err != nil {
				return nil, err
			}
		}
		var signer ssh.Signer
		if conf.Password == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
			if err != nil {
				return nil, err
			}
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(conf.Password))
			if err != nil {
				return nil, err
			}
		}
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
	}
	client.SshConf = sshConfig

	if conf.RemoteHosts != "" {
		remoteSlice := strings.Split(conf.RemoteHosts, ",")
		for _, v := range remoteSlice {
			if v != "" {
				addrPort := strings.Split(strings.TrimSpace(v), ":")
				addr := addrPort[0]
				port := 22
				if len(addrPort) == 2 {
					port, _ = strconv.Atoi(addrPort[1])
				}
				client.Infos = append(client.Infos, &Info{
					Addr: addr,
					Port: port,
				})
			}
		}
	}

	client.Modes = &ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	return client, nil
}

func (cli *Client) RegInfoHostsStr(cmd string, remoteHosts string) {
	for _, v := range strings.Split(remoteHosts, ",") {
		addrPort := strings.Split(strings.TrimSpace(v), ":")
		addr := addrPort[0]
		port := 22
		if len(addrPort) == 2 {
			port, _ = strconv.Atoi(addrPort[1])
		}
		cli.Infos = append(cli.Infos, &Info{
			Addr: addr,
			Port: port,
			Cmd:  cmd,
		})
	}
	return
}

// call RegInfoHostsSlice,RegInfoHostsStr before Connect
func (cli *Client) RegInfoHostsSlice(cmd string, remoteHosts ...string) {
	for _, v := range remoteHosts {
		addrPort := strings.Split(strings.TrimSpace(v), ":")
		addr := addrPort[0]
		port := 22
		if len(addrPort) == 2 {
			port, _ = strconv.Atoi(addrPort[1])
		}
		cli.Infos = append(cli.Infos, &Info{
			Addr: addr,
			Port: port,
			Cmd:  cmd,
		})
	}
	return
}

func (cli *Client) Connect() {
	for _, v := range cli.Infos {
		remotes := fmt.Sprintf("%s:%d", v.Addr, v.Port)
		client, err := ssh.Dial("tcp", remotes, cli.SshConf)
		if err != nil {
			fmt.Printf("[%s] dial error. %+v", remotes, err)
			continue
		}
		v.Client = client
		session, err := client.NewSession()
		if err != nil {
			fmt.Printf("[%s] new session error. %+v", remotes, err)
			continue
		}
		if err := session.RequestPty("xterm", 80, 40, *cli.Modes); err != nil {
			fmt.Printf("[%s] request pty error. %+v", remotes, err)
			continue
		}
		v.Session = session
	}
}

func (cli *Client) Run(cmd string) {
	for _, v := range cli.Infos {
		var stdoutBuf bytes.Buffer
		v.Session.Stdout = &stdoutBuf
		if v.Cmd == "" {
			v.Cmd = cmd
		}
		err := v.Session.Run(cmd)
		if err != nil {
			fmt.Printf("run error.%v", err)
		}
		fmt.Println(v.Addr, v.Session.Stdout)
	}
}

func (cli *Client) MultiRun(cmd string, out chan<- *EchoMsg) {
	defer close(out)
	wg := new(sync.WaitGroup)
	for _, sess := range cli.Infos {
		v := sess
		if v.Session == nil {
			continue
		}
		wg.Add(1)
		var stdoutBuf bytes.Buffer
		v.Session.Stdout = &stdoutBuf
		go func() {
			defer wg.Done()
			if v.Cmd == "" {
				v.Cmd = cmd
			}
			if cli.c.ReplacePattern != "" {
				v.Cmd = strings.ReplaceAll(v.Cmd, cli.c.ReplacePattern, v.Addr)
			}
			err := v.Session.Run(v.Cmd)
			msg := &EchoMsg{
				Addr:    v.Addr,
				Content: fmt.Sprint(v.Session.Stdout),
				err:     err,
			}
			out <- msg
		}()
	}
	wg.Wait()
}

func (cli *Client) Close() {
	for _, v := range cli.Infos {
		if v.Client != nil {
			if err := v.Client.Close(); err != nil {
				fmt.Printf("close error. %v", err)
			}
		}
	}
}
