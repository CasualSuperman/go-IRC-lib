package main

import (
	irc "./_obj/irc"
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
)

func main() {
	user := irc.NewUser("CasualSuperman", "rwertman", false, "Robert Wertman", "0176092")
	conn, err := irc.Connect("irc.foonetic.net:6667", user)
	if err != nil {
		panic(err)
	}
	i := 0
	done := false
	quit := make(chan bool)
	go func(quit chan bool) {
		for !done {
			data := <-signal.Incoming
			if data.String() == "SIGINT: interrupt" {
				quit <- true
			} else {
				fmt.Println(data.String())
			}
		}
	}(quit)
	input := make(chan string)
	for !done {
		i++
		select {
			case data, ok := <-conn.Recv():
				fmt.Printf(data.Tmpl(), data.Data()...)
				if !ok {
					done = true
				}
			case send := <-input:
				if string(send[0]) != "/"{
					conn.Send(irc.NewPrivateMessage("#ufeff", send))
				} else {
					if pos := strings.Index(strings.ToLower(send), "whois"); pos == 1 {
						conn.Write([]byte(": WHOIS" + send[pos+5:] + "\n"))
					} else {
						conn.Write([]byte(": " + send[1:] + "\n"))
					}
				}
			case done = <-quit:
		}
		if i == 20 {
			conn.Send(irc.NewJoinMessage("#ufeff"))
			go func() {
				line := ""
				reader := bufio.NewReader(os.Stdin)
				for {
					var buf []byte
					isPrefix := true
					for isPrefix {
						buf, isPrefix, err = reader.ReadLine()
						if err != nil {
							if err.Error() == "EOF" {
								done = true
							} else {
								fmt.Printf("ERROR: %s", err.Error())
							}
						} else {
							line += string(buf)
						}
					}
					if !done {
						input <- line
						line = ""
					}
				}
			}()
		}
	}
	conn.Send(irc.NewQuitMessage("Closed"))
}
