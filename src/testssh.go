package main

import (
	"fmt"
	"gassh"
	"log"
)

func main() {
	gs:=gassh.Password("root","1")
	conn,err := gs.Connect("192.168.1.200:22")
	if err != nil {
		log.Fatal(err)

	}
	defer conn.Close()
	ls,err := conn.ExecShell("ls")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(ls))
	df,err := conn.ExecShell("df")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(df))
}
