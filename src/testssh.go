package main

import (
	"io/ioutil"
	"log"
)

func main() {
	err := ioutil.WriteFile("a.txt",[]byte("hello"),0644)
	if err != nil {
		log.Fatal(err)
	}
}
