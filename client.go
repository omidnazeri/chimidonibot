// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore
package main

import (
	"./receiver"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var addr = flag.String("addr", "79.175.163.71:3000", "http service address")

func sendAutoChat(c *websocket.Conn, interrupt *chan os.Signal, done *chan struct{}) {
	//ticker := time.NewTicker(time.Second)
	//defer ticker.Stop()

	for {
		select {
		case <-*done:
			return
		//case <-ticker.C:
		//	err := c.WriteMessage(websocket.TextMessage, []byte("{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NTAzOTg3LCJpYXQiOjE1NTY3NDAwNjUsImV4cCI6MTYxOTgxMjA2NX0.zBOOrurTVRgJNyzBsUwzU7U7vQatj4p5N8U_wOgJXxM\",\"action\":\"CHAT\",\"chattext\":\"salam\"}"))
		//	if err != nil {
		//		log.Println("write:", err)
		//		return
		//	}
		case <-*interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-*done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func messageReceiver(c *websocket.Conn, done *chan struct{}) {
	defer close(*done)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		//n := bytes.Index(message, []byte{0})
		//messageString := ""
		//for _, e := range message {
		//	messageString += string(e)
		//}
		//for _, c := range message {
		//	fmt.Print(string(c))
		//}
		messageString := string(message[:])

		var receiverObj receiver.MessageReceiver

		data := receiver.GeneralMessageReceiver{}
		data.FillData(message)

		switch data.Command {
		case "PLAYERCOUNT":
			// do nothing
		case "SHOWCHAT":
			// do nothing
		case "QUESTION":
			receiverObj = &receiver.QuestionMessageReceiver{}
		default:
			receiverObj = &receiver.GeneralMessageReceiver{}
		}

		if receiverObj != nil {
			fmt.Println(messageString)
			receiverObj.FillData(message)
			receiverObj.OnReceive(c, done)
		}
	}
}

func sendHello(c *websocket.Conn) {
	err := c.WriteMessage(websocket.TextMessage, []byte("{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NTAzOTg3LCJpYXQiOjE1NTI4MzQ2MzcsImV4cCI6MTYxNTkwNjYzN30.00hCKmHkkVkZJyCXP0SUiEXGT-EbHqQJTwDH6DBVTRI\",\"action\":\"HELLO\",\"firstname\":\"\",\"lastname\":\"\",\"debugstring\":\"\",\"haslife\":0,\"balance\":0,\"coins\":0,\"buildNumber\":0,\"AppVersion\":\"Bot\",\"AndroidVersion\":\"Bot\",\"DeviceName\":\"Bot\"}"))
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go messageReceiver(c, &done)
	sendHello(c)

	sendAutoChat(c, &interrupt, &done)
}
