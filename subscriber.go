package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SubscriberClient struct {
	UserName  string
	Password  string
	ClientID  string
	BrokerURL string
	Topics    map[string]byte
}

func subscribe(client mqtt.Client, subClient SubscriberClient) {
	token := client.SubscribeMultiple(subClient.Topics, func(c mqtt.Client, m mqtt.Message) {
		log.Printf("Received message: %s from topic %s\n", m.Payload(), m.Topic())
	})
	token.Wait()
	log.Printf("subscribed to topics: %v", subClient.Topics)
}

func sub() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	subscriberClient := SubscriberClient{
		BrokerURL: "",           // replace with your broker url and port
		UserName:  "sub1",       // replace with your username
		Password:  "secret",     // replace with your password
		ClientID:  "sub_client", // replace with your client ID
		Topics: map[string]byte{ // replace with your topics
			"topic/device/temperature": 0,
			"topic/device/moisture":    0,
			"topic/device/light":       0,
		},
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(subscriberClient.BrokerURL)
	opts.SetClientID(subscriberClient.ClientID)
	opts.SetUsername(subscriberClient.UserName)
	opts.SetPassword(subscriberClient.Password)

	tlsConfig := newTLSConfig()
	opts.SetTLSConfig(tlsConfig)

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("subscriber connected")
		subscribe(c, subscriberClient)
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("subscriber lost connection: %v\n", err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("subscriber failed to connect: %v", token.Error())
	}

	<-c
	client.Disconnect(2000) // gracefully close connection after 2 seconds
}
