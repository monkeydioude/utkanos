package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const (
	server  = "192.168.1.109:2224"
	channel = "#general"
	pass    = "1234"
)

var nickname = "utkanos"

type Channels []string

var channels Channels = []string{"#general", "#animals", "#u-t-k-a", "#nooOOOOooos", "#seg", "#f_a_u_l_t", "#bitcoin"}

func (c Channels) JoinAll(send func(string)) {
	send(fmt.Sprintf("JOIN %s", strings.Join(channels, ",")))
}

func (c Channels) SendAll(send func(string)) {
	for _, channel := range c {
		send(fmt.Sprintf("PRIVMSG %s :Mrrrrrrrrr", channel))
	}
}

func (c Channels) SendRandom(send func(string)) {
	send(fmt.Sprintf("PRIVMSG %s :Mrrrrrrrrr", c[rand.Intn(len(channels))]))
}

func sendFactory(conn net.Conn) func(string) {
	return func(msg string) {
		fmt.Fprintf(conn, "%s\r\n", msg)
		fmt.Println(msg)
	}
}

func spam(sendFn func(string)) {
	for {
		channels.SendRandom(sendFn)
		time.Sleep(20 * time.Millisecond)
	}
}

func waitMsg(conn net.Conn, sendFn func(string)) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("<< %s\n", line)

		// Respond to PING messages to stay connected
		if len(line) > 4 && line[:4] == "PING" {
			sendFn("PONG " + line[5:])
		}
	}
}

func permaReadAndSend(sendFn func(string)) {
	input := bufio.NewScanner(os.Stdin)

	time.Sleep(20 * time.Millisecond)
	for fmt.Print(">> "); input.Scan(); {
		sendFn(input.Text())
		time.Sleep(100 * time.Millisecond)
		fmt.Print(">> ")
	}
}

func main() {
	var conn net.Conn
	var err error

	if len(os.Args) > 1 {
		nickname = os.Args[1]
	}
	for conn == nil {
		conn, err = net.Dial("tcp", server)
		if err != nil {
			fmt.Println("Error connecting, retry in 5s", "err", err.Error())
			time.Sleep(5 * time.Second)
		}
	}
	defer conn.Close()

	send := sendFactory(conn)
	fmt.Printf("Connected to server. addr=%s", server)
	// Authenticate with the server
	send("PASS " + pass)
	send("NICK " + nickname)
	send("USER " + nickname + " 0 * :" + nickname)
	channels.JoinAll(send)
	// time.Sleep(200 * time.Millisecond)
	send("WHOIS " + nickname)

	// go spam(send)
	go permaReadAndSend(send)
	waitMsg(conn, send)
}
