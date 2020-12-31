package ssh

import (
	"flag"
	"github.com/BurntSushi/toml"
	"time"
)

type Config struct {
	User        string
	Password    string
	PrivateKey  string
	Timeout     time.Duration
	Ciphers     []string
	RemoteHosts []RemoteHost
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
	_timeout     time.Duration
	_ciphers     string
	_remoteHosts RemoteHostFlag
)

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
	addFlags(flag.CommandLine)
}

func addFlags(fs *flag.FlagSet) {
	_remoteHosts.Set("10.85.132.235,10.85.132.217")

	fs.StringVar(&_user, "user", "", "user")
	fs.StringVar(&_password, "password", "", "password")
	fs.StringVar(&_privateKey, "key", "", "private key path")
	fs.DurationVar(&_timeout, "timeout", 10*time.Second, "timeout")
	fs.Var(&_remoteHosts, "ips", "remote ssh ip and port. eg:'127.0.0.1:22,127.0.0.2:22'")
}

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
