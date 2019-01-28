package gassh

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
)

type SshConfig struct {
	config *ssh.ClientConfig
	addr   string
}

func Password(username string, password string) *SshConfig {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return &SshConfig{config: config}
}

func PrivateKey(username string, privitekeypath string) *SshConfig {
	key, err := ioutil.ReadFile(privitekeypath)
	fmt.Println(privitekeypath)
	if err != nil {
		panic("key not found")
	}
	signer, err := ssh.ParsePrivateKey(key)
	fmt.Println()
	if err != nil {
		log.Fatal(err)
	}
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return &SshConfig{config: config}
}

type SshConn struct {
	conn *ssh.Client
}

//type Session struct {
//	session *ssh.Session
//}

func (sc *SshConfig) Connect(addr string) (*SshConn, error) {

	client, err := ssh.Dial("tcp", addr, sc.config)
	if err != nil {
		return nil, err
	}
	return &SshConn{client}, nil
}

func (sc *SshConn) Close() {
	sc.conn.Close()
}

func (sc *SshConn) ExecShell(cmd string) ([]byte, error) {

	session, err := sc.conn.NewSession()
	defer session.Close()
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	session.Stdout = &b
	session.Run(cmd)

	return b.Bytes(), nil
}

func (sc *SshConn) ExecShellToString(cmd string) (string, error) {

	session, err := sc.conn.NewSession()
	defer session.Close()
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	session.Stdout = &b
	session.Run(cmd)

	return b.String(), nil
}

func (sc *SshConn) NewSession() {
	session, err := sc.conn.NewSession()
	defer session.Close()
	if err != nil {
		log.Fatal("unable to create session: ", err)
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}
	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}
	session.Wait()
}
