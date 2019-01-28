package gassh

import (
	"fmt"
	"log"
)

func DocLsFile() {
	// 返回一个sshconfig结构体
	config := Password("root", "123456")
	// 秘钥认证
	//PrivateKey("root", "/root/.ssh/id_rsa")
	//连接某台机器,返回一个连接
	conn, err := config.Connect("192.168.1.10:22")
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	out, err := conn.ExecShell("ls")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)
}
