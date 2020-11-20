/*
Copyright LLC Newity. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Google cloud Pub/Sub implementation of Storage
package storage

import (
	"errors"
	stan "github.com/nats-io/stan.go"
	"sync"
)

type Nats struct {
	Connection    stan.Conn
	channels      []string
	subscriptions []stan.Subscription
}

func NewNats(clusterID, clientID, natsURL string, maxPubAcksInflight int) (*Nats, error) {
	if maxPubAcksInflight == 0 {
		maxPubAcksInflight = stan.DefaultMaxPubAcksInflight // if arg is null, set to default (16384)
	}
	conn, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL), stan.MaxPubAcksInflight(maxPubAcksInflight))
	if err != nil {
		return nil, err
	}

	return &Nats{
		Connection: conn,
	}, nil
}

func (n *Nats) InitChannelsStorage(channels []string) error {
	n.channels = channels
	return nil
}

// Put stores message to topic.
func (n *Nats) Put(topic string, msg []byte) error {
	n.Connection.Publish(topic, msg) // sync call, wait for ACK from NATS Streaming
	return nil
}

// Get reads one message from the topic and closes channel.
func (n *Nats) Get(topic string) ([]byte, error) {
	var data []byte
	sub, err := n.Connection.QueueSubscribe(topic, topic, func(m *stan.Msg) {
		data = m.Data
		m.Ack()
	})
	n.subscriptions = append(n.subscriptions, sub)
	return data, err
}

// GetStream reads a stream of messages from topic and writes them to the channel.
func (n *Nats) GetStream(topic string) (<-chan []byte, <-chan error) {
	ch, errch := make(chan []byte), make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		sub, err := n.Connection.QueueSubscribe(topic, topic, func(m *stan.Msg) {
			ch <- m.Data
			m.Ack()
		})
		n.subscriptions = append(n.subscriptions, sub)
		if err != nil {
			errch <- err
		}
	}()
	wg.Wait()
	return ch, errch
}

// Detele does not work for Nats.
func (n *Nats) Delete(key string) error {
	return errors.New("Not implemented in Nats")
}

// Close stops all running goroutines related to topics.
func (n *Nats) Close() error {
	for _, sub := range n.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			return err
		}
		if err := sub.Close(); err != nil {
			return err
		}
	}
	return n.Connection.Close()
}
