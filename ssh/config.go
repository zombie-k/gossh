package ssh

import (
	"flag"
	"github.com/BurntSushi/toml"
)

type Config struct {
	User        string
	Password    string
	PrivateKey  string
	Timeout     int64
	Ciphers     []string
	RemoteHosts string
	Cmd         string
	EchoIp      bool
	IdRsa       string
}

type RemoteHost struct {
	Addr string
	Port int
}

var (
	confPath     string
	Conf         = &Config{}
	_user        string
	_password    string
	_privateKey  string
	_timeout     int64
	_ciphers     string
	_remoteHosts string
)

func init() {
	flag.StringVar(&confPath, "conf", "ssh/ssh.toml", "default config path")
	flag.StringVar(&_user, "user", "", "user")
	flag.StringVar(&_password, "password", "", "password")
	flag.StringVar(&_privateKey, "key", "", "private key path")
	flag.Int64Var(&_timeout, "timeout", 10, "timeout")
	flag.StringVar(&_remoteHosts, "host", "10.85.132.235,10.85.132.217", "remote ssh ip and port. eg:'127.0.0.1:22,127.0.0.2:22'")
}

/*
func addFlags(fs *flag.FlagSet) {
	_remoteHosts.Set("10.85.132.235,10.85.132.217")

	fs.StringVar(&_user, "user", "", "user")
	fs.StringVar(&_password, "password", "", "password")
	fs.StringVar(&_privateKey, "key", "", "private key path")
	fs.DurationVar(&_timeout, "timeout", 10*time.Second, "timeout")
	fs.Var(&_remoteHosts, "ips", "remote ssh ip and port. eg:'127.0.0.1:22,127.0.0.2:22'")
}
*/

func Init() (err error) {
	if confPath == "" {
		Conf = &Config{
			User:        _user,
			Password:    _password,
			PrivateKey:  _privateKey,
			Timeout:     _timeout,
			RemoteHosts: _remoteHosts,
		}
	} else {
		_, err = toml.DecodeFile(confPath, &Conf)
	}
	return
}
