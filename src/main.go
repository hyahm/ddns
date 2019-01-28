package main

import (
	"flag"
	"fmt"
	"gassh"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var IP string
var RemoteIp string
var RemoteFile string
var AfterExec string

func main() {
	afterExec := flag.String("a", "nginx -s reload", "remote ip")
	interval := flag.Int("t", 1, "interval time, minute")
	remoteIp := flag.String("r", "", "remote ip")
	remoteFile := flag.String("f", "", "remote server file path")
	port := flag.Int("P", 22, "ssh port")
	username := flag.String("u", "root", "ssh username")
	password := flag.String("p", "", "ssh password")
	keyfile := flag.String("k", "~/.ssh/id_rsa", "ssh private key")
	flag.Parse()
	if *remoteIp == "" {
		log.Fatal("remote ip must be need")
	}
	if *remoteFile == "" {
		log.Fatal("remotefile ip must be need")
	}
	RemoteIp = *remoteIp
	AfterExec = *afterExec
	RemoteFile = *remoteFile
	for {
		resp, err := http.Get("http://ip.hyahm.com")
		if err != nil {
			time.Sleep(time.Second * 3)
			continue
		}
		s, _ := ioutil.ReadAll(resp.Body)
		s = s[:len(s)-1]

		if string(s) != IP {
			if *password != "" {
				execShellWithPassword(*username, *password, *port, string(s))
			} else {
				execShellWithKey(*username, *keyfile, *port, string(s))
			}
			IP = string(s)
		}
		fmt.Println(string(s))
		time.Sleep(time.Minute * time.Duration(*interval))
	}

}

func execShellWithPassword(username string, password string, port int, newip string) error {
	sconf := gassh.Password(username, password)
	sshconn, err := sconf.Connect(fmt.Sprintf("%s:%d", RemoteIp, port))
	if err != nil {
		return err
	}
	defer sshconn.Close()
	exec(sshconn, newip)
	return nil
}

func execShellWithKey(username string, key string, port int, newip string) error {
	sconf := gassh.PrivateKey(username, key)
	sshconn, err := sconf.Connect(fmt.Sprintf("%s:%d", RemoteIp, port))
	if err != nil {
		return err
	}
	defer sshconn.Close()
	err = exec(sshconn, newip)
	if err != nil {
		return err
	}
	return nil
}

func exec(conn *gassh.SshConn, newip string) error {

	_, err := conn.ExecShell(fmt.Sprintf(` sed -i "s/%s/%s/g`, IP, newip))
	if err != nil {
		return err
	}
	fmt.Println("change success")
	_, err = conn.ExecShell(AfterExec)
	if err != nil {
		return err
	}
	fmt.Println("reload service success")
	return nil
}
