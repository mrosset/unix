package main

import (
	"code.google.com/p/go.crypto/ssh"
	"code.google.com/p/go.crypto/ssh/terminal"
	"fmt"
	"io/ioutil"
	"log"
)

var (
	config = &ssh.ServerConfig{
		PasswordCallback: authUserPass,
	}
)

func init() {
	log.SetPrefix("sshs: ")
	log.SetFlags(log.Lshortfile)
}

func main() {
	pem, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		log.Fatal(err)
	}
	err = config.SetRSAPrivateKey(pem)
	if err != nil {
		log.Fatal(err)
	}
	lis, err := ssh.Listen("tcp", "0.0.0.0:2022", config)
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()
	conn, err := lis.Accept()
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Handshake()
	if err != nil {
		log.Fatal(err)
	}

	for {
		channel, err := conn.Accept()
		if err != nil {
			log.Fatal(err)
		}
		channel.Accept()

		term := terminal.NewTerminal(channel, "> ")
		serverTerm := &ssh.ServerTerminal{
			Term:    term,
			Channel: channel,
		}
		go func() {
			defer channel.Close()
			for {
				line, err := serverTerm.ReadLine()
				if err != nil {
					break
				}
				fmt.Println(line)
			}
		}()
	}
}

func authUserPass(c *ssh.ServerConn, user, pass string) bool {
	return true
}
