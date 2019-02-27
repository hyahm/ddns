package main

import (
	"flag"
	"fmt"
	"galog"
	"gamail"
	"gassh"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var IP string
var RemoteIp string
var RemoteFile string
var AfterExec string

func main() {
	//gaconfig.InitConf("ddns")

	galog.DefInitLogger("log",0,false)

	afterExec := flag.String("a", "nginx -s reload", "reload http server command")
	interval := flag.Int("t", 60, "interval time, default 60 Second")
	remoteIp := flag.String("r", "", "server ip")
	remoteFile := flag.String("f", "", "remote server nginx config file path")
	port := flag.Int("P", 22, "ssh port")
	username := flag.String("u", "root", "ssh username")
	password := flag.String("p", "", "ssh password")
	keyfile := flag.String("k", "~/.ssh/id_rsa", "ssh private key")
	defaultip := flag.String("d", "", "default ip")

	emailuser := flag.String("eu","","email username")
	emailpwd := flag.String("ep","","email password")
	emailport := flag.Int("eo",0,"email username")
	emailto := flag.String("et","","email to who")
	flag.Parse()
	fmt.Println(flag.Lookup("f"))
	fmt.Println(*remoteFile)
	//func main() {
	//
	//	verbose := flag.String("verbose", "on", "Turn verbose mode on or off.")
	//	flag.Parse()
	//
	//	fmt.Println(flag.Lookup("verbose")) // print the Flag struct
	//
	//	fmt.Println(*verbose)
	//
	//}
	os.Exit(0)
	if *remoteIp == "" {
		log.Fatal("server ip must be need, use -r serverip")
	}
	if *remoteFile == "" {
		log.Fatal("remotefile ip must be need, ps: -f /usr/local/nginx/conf/vhost/aa.conf")
	}
	if *defaultip == "" {
		log.Fatalf("defaultip ip must be need in %s",*remoteFile)
	}
	IP = *defaultip
	RemoteIp = *remoteIp
	AfterExec = *afterExec
	RemoteFile = *remoteFile
	for {
		resp, err := http.Get("http://ip.hyahm.com")
		if err != nil {
			galog.Error(err.Error())
			time.Sleep(time.Second * 3)
			continue
		}
		s, _ := ioutil.ReadAll(resp.Body)
		s = s[:len(s)-1]

		if string(s) != IP {
			if *emailuser != "" && *emailpwd != "" && *emailport != 0 && *emailto != "" {
				email(*emailuser,*emailpwd ,*emailport,string(s),*emailto)
			}

			if *password != "" {
				err = execShellWithPassword(*username, *password, *port, string(s))
				if err != nil {
					email(*emailuser,*emailpwd ,*emailport,"update fail",*emailto)
				}
			} else {
				err = execShellWithKey(*username, *keyfile, *port, string(s))
				if err != nil {
					email(*emailuser,*emailpwd ,*emailport,"update fail",*emailto)
				}
			}
			email(*emailuser,*emailpwd ,*emailport,"update success",*emailto)
			IP = string(s)
		}
		fmt.Println(string(s))
		time.Sleep(time.Second * time.Duration(*interval))
	}

}

func execShellWithPassword(username string, password string, port int, newip string) error {
	sconf := gassh.Password(username, password)
	sshconn, err := sconf.Connect(fmt.Sprintf("%s:%d", RemoteIp, port))
	if err != nil {
		galog.Error(err.Error())
		return err
	}
	defer sshconn.Close()
	exec(sshconn, newip)
	return nil
}

func execShellWithKey(username string, key string, port int, newip string) error {
	sconf := gassh.PrivateKey(username, key)
	// 连接远程机器
	sshconn, err := sconf.Connect(fmt.Sprintf("%s:%d", RemoteIp, port))
	if err != nil {
		galog.Error(err.Error())
		return err
	}
	fmt.Println("远程连接成功")
	defer sshconn.Close()
	err = exec(sshconn, newip)
	if err != nil {
		return err
	}
	return nil
}

func exec(conn *gassh.SshConn, newip string) error {
	rfs := strings.Split(RemoteFile,";")
	for _,v := range rfs {
		_, err := conn.ExecShell(fmt.Sprintf(" sed -i 's/%s/%s/g' %s", IP, newip,v))
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	fmt.Println("change success")
	_, err = conn.ExecShell(AfterExec)
	if err != nil {
		return err
	}
	fmt.Println("reload service success")
	return nil
}

func email(username string,password string,port int,content string, to string) {
	econf := gamail.Newmailconfig()
	econf.Username = username
	econf.Password = password
	econf.Subject = "ip changed"
	econf.Content = content
	econf.Port = port
	econf.Tolist = []string{to}
	econf.SendMail()
}