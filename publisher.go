package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Publisher struct {
	ID        string
	UserName  string
	Password  string
	ClientID  string
	BrokerURL string
	Topic     string
}

func NewPublisher(id, username, topic, clientID string) Publisher {
	return Publisher{
		ID:        id,
		BrokerURL: "", // replace with your broker url and port
		UserName:  username,
		Password:  "secret", // replace with your password
		ClientID:  clientID,
		Topic:     topic,
	}
}

func newPublisherClient(publisher Publisher) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(publisher.BrokerURL)
	opts.SetClientID(publisher.ClientID)
	opts.SetUsername(publisher.UserName)
	opts.SetPassword(publisher.Password)

	tlsConfig := newTLSConfig()
	opts.SetTLSConfig(tlsConfig)

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("%s lost connection: %v\n", publisher.ID, err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("%s publisher failed to connect: %v", publisher.ID, token.Error())
	}

	return client
}

func (p Publisher) publishMessage(client mqtt.Client, duration time.Duration) {
	i := 0
	for range time.Tick(duration * time.Second) {
		text := fmt.Sprintf("%s is currently: %d", p.ID, i)
		token := client.Publish(p.Topic, 0, false, text)
		token.Wait()
		i++
	}
}

func newTLSConfig() *tls.Config {
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile("ssl_cert.crt")
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	certPool.AppendCertsFromPEM(ca)
	return &tls.Config{
		RootCAs: certPool,
	}
}

func pub() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	temperaturePublisher := NewPublisher("Temperature", "pub1", "topic/device/temperature", "temp_pub_client")
	moisturePublisher := NewPublisher("Moisture", "pub2", "topic/device/moisture", "moisture_pub_client")
	lightPublisher := NewPublisher("Light", "pub3", "topic/device/light", "light_pub_client")

	temperatureClient := newPublisherClient(temperaturePublisher)
	moistureClient := newPublisherClient(moisturePublisher)
	lightClient := newPublisherClient(lightPublisher)

	go func() {
		temperaturePublisher.publishMessage(temperatureClient, 3)
	}()

	go func() {
		moisturePublisher.publishMessage(moistureClient, 5)
	}()

	go func() {
		lightPublisher.publishMessage(lightClient, 7)
	}()

	<-c
}
