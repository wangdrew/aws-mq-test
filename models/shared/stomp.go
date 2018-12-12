package shared

import (
	"crypto/tls"
	"fmt"
	"github.com/go-stomp/stomp"
	"log"
)

type StompClient struct {
	conn *stomp.Conn
}

func NewStompClient(host, username, password string) *StompClient {
	tlsConfig := &tls.Config{}
	tlsConn, err := tls.Dial("tcp", host+":61614", tlsConfig)
	if err != nil {
		log.Fatal("Dialing stomp server:", err)
	}
	conn, err := stomp.Connect(tlsConn, stomp.ConnOpt.Login(username, password))
	if err != nil {
		log.Fatal("Connecting to stomp server:", err)
	}
	return &StompClient{conn}
}

func (c *StompClient) Send(destination string, message string) {
	err := c.conn.Send(destination, "text/plain", []byte(message))
	if err != nil {
		log.Fatal("Sending message to stomp server:", err)
	}
	fmt.Println("Sent - destination:", destination, "message:", message)
}

type StompMessageHandler func(destination string, message string)

func (c *StompClient) Subscribe(destination string, handler StompMessageHandler) {
	go func() {
		sub, err := c.conn.Subscribe(destination, stomp.AckClient)
		if err != nil {
			log.Fatal("Subscribing to stomp server:", err)
		}
		for {
			msg := <-sub.C
			if msg.Err != nil {
				log.Fatal("Receving messages from stomp server:", err)
			}
			handler(destination, string(msg.Body))
			err = c.conn.Ack(msg)
			if err != nil {
				log.Fatal("Acking message to stomp server:", err)
			}
		}
	}()
}
