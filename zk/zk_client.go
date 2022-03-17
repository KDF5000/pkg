package zk

import (
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"
)

type optionFunc func(opt *zkOption)

type zkOption struct {
	SessionTimeout time.Duration
	Auth           []byte
}

type ZKClient struct {
	addr   []string
	option zkOption
	// Export public method of *zk.Conn
	*zk.Conn
}

func WithSessionTimeout(timeout time.Duration) optionFunc {
	return func(opt *zkOption) {
		opt.SessionTimeout = timeout
	}
}

func WithAuth(key []byte) optionFunc {
	return func(opt *zkOption) {
		// should copy?
		opt.Auth = key
	}
}

// addrs can be domain name or specified ip address
func NewZKClient(addr []string, options ...optionFunc) (*ZKClient, error) {
	cli := &ZKClient{
		addr: addr,
	}
	for _, option := range options {
		option(&cli.option)
	}

	if cli.option.SessionTimeout <= 0 {
		cli.option.SessionTimeout = 10 * time.Second
	}

	// TODO(kongdefei): use self-defined addr resolver to cache
	// dns result and refresh dns address periodically
	err := cli.ReConnect()
	if err != nil {
		return nil, err
	}

	return cli, nil
}

// ReConnet will re-connect zk servers and get an new zk.Conn
// Call this after Close if you want to use the same ZKClient
func (cli *ZKClient) ReConnect() (err error) {
	if cli.Conn != nil && cli.State() != zk.StateDisconnected {
		return fmt.Errorf("the connection is normal")
	}

	cli.Conn, _, err = zk.Connect(cli.addr, cli.option.SessionTimeout, zk.WithHostProvider(&DNSHostProvider{}))
	if err != nil {
		return err
	}

	if cli.option.Auth != nil && len(cli.option.Auth) > 0 {
		err = cli.AddAuth("digest", []byte(cli.option.Auth))
	}

	return err
}

// func ConnectWithDialer(servers []string, sessionTimeout time.Duration, dialer Dialer) (*Conn, <-chan Event, error)
// func (c *Conn) AddAuth(scheme string, auth []byte) error
// func (c *Conn) Children(path string) ([]string, *Stat, error)
// func (c *Conn) ChildrenW(path string) ([]string, *Stat, <-chan Event, error)
// func (c *Conn) Close()
// func (c *Conn) Create(path string, data []byte, flags int32, acl []ACL) (string, error)
// func (c *Conn) CreateContainer(path string, data []byte, flags int32, acl []ACL) (string, error)
// func (c *Conn) CreateProtectedEphemeralSequential(path string, data []byte, acl []ACL) (string, error)
// func (c *Conn) CreateTTL(path string, data []byte, flags int32, acl []ACL, ttl time.Duration) (string, error)
// func (c *Conn) Delete(path string, version int32) error
// func (c *Conn) Exists(path string) (bool, *Stat, error)
// func (c *Conn) ExistsW(path string) (bool, *Stat, <-chan Event, error)
// func (c *Conn) Get(path string) ([]byte, *Stat, error)
// func (c *Conn) GetACL(path string) ([]ACL, *Stat, error)
// func (c *Conn) GetW(path string) ([]byte, *Stat, <-chan Event, error)
// func (c *Conn) IncrementalReconfig(joining, leaving []string, version int64) (*Stat, error)
// func (c *Conn) Multi(ops ...interface{}) ([]MultiResponse, error)
// func (c *Conn) Reconfig(members []string, version int64) (*Stat, error)
// func (c *Conn) Server() string
// func (c *Conn) SessionID() int64
// func (c *Conn) Set(path string, data []byte, version int32) (*Stat, error)
// func (c *Conn) SetACL(path string, acl []ACL, version int32) (*Stat, error)
// func (c *Conn) SetLogger(l Logger)
// func (c *Conn) State() State
// func (c *Conn) Sync(path string) (string, error)
