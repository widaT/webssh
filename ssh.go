package webssh

import (
	"io/ioutil"
	"time"

	"golang.org/x/crypto/ssh"
)

type AuthModel int8

const (
	PASSWORD AuthModel = iota + 1
	PUBLICKEY
)

type SSHClientConfig struct {
	AuthModel AuthModel
	HostAddr  string
	User      string
	Password  string
	KeyPath   string
	Timeout   time.Duration
}

func SSHClientConfigPassword(hostAddr, user, Password string) *SSHClientConfig {
	return &SSHClientConfig{
		Timeout:   time.Second * 5,
		AuthModel: PASSWORD,
		HostAddr:  hostAddr,
		User:      user,
		Password:  Password,
	}
}

func SSHClientConfigPulicKey(hostAddr, user, keyPath string) *SSHClientConfig {
	return &SSHClientConfig{
		Timeout:   time.Second * 5,
		AuthModel: PUBLICKEY,
		HostAddr:  hostAddr,
		User:      user,
		KeyPath:   keyPath,
	}
}

func NewSSHClient(conf *SSHClientConfig) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		Timeout:         conf.Timeout,
		User:            conf.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //忽略know_hosts检查
	}
	switch conf.AuthModel {
	case PASSWORD:
		config.Auth = []ssh.AuthMethod{ssh.Password(conf.Password)}
	case PUBLICKEY:
		signer, err := getKey(conf.KeyPath)
		if err != nil {
			return nil, err
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	}
	c, err := ssh.Dial("tcp", conf.HostAddr, config)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func getKey(keyPath string) (ssh.Signer, error) {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(key)
}
