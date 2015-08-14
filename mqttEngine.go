package main

import (
	"encoding/json"
	"net/url"
	"sync"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

type mqttEngine struct {
	murl      *url.URL
	cli       *client.Client
	cliLock   sync.Mutex
	connected bool

	//publisher  *publisher
	pubTicker  *time.Ticker
	pollTicker *time.Ticker
}

func newMqttEngine() (*mqttEngine, error) {

	murl, err := url.Parse(*mqttURL)

	if err != nil {
		return nil, err
	}

	mq := &mqttEngine{}

	// Create an MQTT Client.
	cli := client.New(&client.Options{
		ErrorHandler: mq.handleClientError,
	})

	mq.murl = murl
	mq.cli = cli

	mq.attemptConnect()

	//mq.publisher = publisher
	mq.pollTicker = time.NewTicker(time.Second * 1)
	mq.pubTicker = time.NewTicker(time.Second * 15)

	go poll(mq)
	go publish(mq,"Ready to publish")

	return mq, nil
}

func (mq *mqttEngine) attemptConnect() bool {

	log.Debugf("Attempt Connect")

	mq.cliLock.Lock()

	defer mq.cliLock.Unlock()

	if mq.connected {
		log.Debugf("already connected")
		return mq.connected
	}

	log.Debugf("connecting to %s", mq.murl.Host)

	co := &client.ConnectOptions{
		Network:  mq.murl.Scheme,
		Address:  mq.murl.Host,
		ClientID: []byte("mqtt-webbrick"),
	}

	if mq.murl.User != nil {
		co.UserName = []byte(mq.murl.User.Username())
		if pass, ok := mq.murl.User.Password(); ok {
			co.Password = []byte(pass)
		}
	}

	// Connect to the MQTT Server.
	if err := mq.cli.Connect(co); err != nil {
		log.Errorf("failed to connect: %s", err)
		mq.connected = false
	} else {
		mq.connected = true
	}

	return mq.connected
}

func (mq *mqttEngine) handleClientError(err error) {

	log.Errorf("client error: %s", err)

	mq.cliLock.Lock()
	defer mq.cliLock.Unlock()

	mq.connected = false

	go func() {
		if err := mq.cli.Disconnect(); err != nil {
			log.Errorf("client disconnect error: %s", err)
		}
		log.Debugf("client disconnected")
	}()
}

func (mq *mqttEngine) disconnect() error {
	if mq.connected {
		return mq.cli.Disconnect()
	}
	return nil
}

func publish(mq *mqttEngine, message string) {

	for t := range mq.pubTicker.C {

		//metrics := mq.publisher.export()

		if !mq.attemptConnect() {
			log.Warningf("publish failed: not connected")
		}

		payload, _ := json.Marshal(struct {
			Time    int64                  `json:"ts"`
			Payload map[string]interface{} `json:"payload"`
		}{
			Time:    t.Unix(),
			Payload: []byte(message),
		})

		log.Debugf("publishing to %s length %d", wbTopic, len(payload))

		err := mq.cli.Publish(&client.PublishOptions{
			QoS:       mqtt.QoS0,
			TopicName: []byte(wbTopic),
			Message:   payload,
		})

		if err != nil {
			log.Errorf("error publishing: %s", err)
		}
	}

}



func poll(mq *mqttEngine) {

	for _ = range mq.pollTicker.C {
		//log.Debugf("Flush at %v", t)
		mq.publisher.flush()
	}

}
