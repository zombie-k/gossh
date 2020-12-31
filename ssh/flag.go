package ssh

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type CommaStrSlice []string

func (c *CommaStrSlice) String() string {
	return fmt.Sprint(*c)
}

func (c *CommaStrSlice) Set(value string) error {
	for _, i := range strings.Split(value, ",") {
		*c = append(*c, strings.TrimSpace(i))
	}
	return nil
}

type RemoteHostFlag []RemoteHost

func (rhf *RemoteHostFlag) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	l := len(*rhf)
	for k, v := range *rhf {
		buf.WriteString("{")
		buf.WriteString(v.Addr)
		buf.WriteString(" ")
		buf.WriteString(strconv.Itoa(v.Port))
		buf.WriteString("}")
		if k < l-1 {
			buf.WriteString(" ")
		}
	}
	buf.WriteString("]")
	return buf.String()
}

func (rhf *RemoteHostFlag) Set(value string) error {
	for _, v := range strings.Split(value, ",") {
		port := 22
		kv := strings.Split(v, ":")
		ip := kv[0]
		if len(kv) == 2 {
			port, _ = strconv.Atoi(kv[1])
		}
		if len(ip) > 0 {
			*rhf = append(*rhf, RemoteHost{
				Addr: ip,
				Port: port,
			})
		}
	}
	return nil
}
