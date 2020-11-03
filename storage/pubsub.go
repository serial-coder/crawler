// Google cloud Pub/Sub implementation of Storage
package storage

import (
	"cloud.google.com/go/pubsub"
	"context"
	"google.golang.org/api/option"
	"time"
)

type PubSub struct {
	ctx           context.Context
	client        *pubsub.Client
	topics        map[string]*pubsub.Topic        // name of topic => *pubsub.Topic mapping
	subscriptions map[string]*pubsub.Subscription // name of subscription => *pubsub.Subscription mapping
}

func NewPubSub(project string, opts ...option.ClientOption) (*PubSub, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, project, opts...)
	if err != nil {
		return nil, err
	}
	return &PubSub{client: client}, nil
}

func (p *PubSub) InitChannelsStorage(channels []string) error {
	var err error
	p.ctx = context.Background()

	for _, channel := range channels {
		var topic *pubsub.Topic
		topic = p.client.Topic(channel)
		if topic == nil {
			topic, err = p.client.CreateTopic(p.ctx, channel)
			if err != nil {
				return err
			}
		}

		topic.EnableMessageOrdering = true
		p.topics[channel] = topic

		var sub *pubsub.Subscription
		sub = p.client.Subscription(channel)
		if sub == nil {
			sub, err = p.client.CreateSubscription(p.ctx, channel, pubsub.SubscriptionConfig{
				Topic:                 topic,
				AckDeadline:           10 * time.Second,
				ExpirationPolicy:      25 * time.Hour,
				EnableMessageOrdering: true,
			})
			if err != nil {
				return nil
			}
		}
		sub.ReceiveSettings.Synchronous = true
		p.subscriptions[channel] = sub
	}
	return nil
}

// Put stores message to topic.
func (p *PubSub) Put(topic string, msg []byte) error {
	res := p.topics[topic].Publish(p.ctx, &pubsub.Message{
		ID:          topic,
		Data:        msg,
		OrderingKey: "0",
	})
	<-res.Ready()
	return nil
}

// Get reads message from topic.
func (p *PubSub) Get(topic string) ([]byte, error) {
	var data []byte
	err := p.subscriptions[topic].Receive(p.ctx, func(ctx context.Context, m *pubsub.Message) {
		copy(data, m.Data)
		m.Ack()
	})
	if err != context.Canceled {
		return nil, err
	}
	return data, nil
}

// Detele deletes topic and subscription specified by key.
func (p *PubSub) Delete(key string) error {
	err := p.topics[key].Delete(p.ctx)
	if err != nil {
		return err
	}
	err = p.subscriptions[key].Delete(p.ctx)
	return err
}

// Close stops all running goroutines related to topics.
func (p *PubSub) Close() error {
	for _, topic := range p.topics {
		topic.Stop()
	}
	return nil
}
