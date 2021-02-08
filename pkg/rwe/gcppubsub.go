package rwe

import (
	"sync"
	"context"

	"github.com/sirupsen/logrus"

	"cloud.google.com/go/pubsub"
)

//initPubSubClient creates a client for a given projectID
//Env variable GOOGLE_APPLICATION_CREDENTIALS="[PATH]" must be set
//https://cloud.google.com/docs/authentication/production 
func initPubSubClient(projectID string) context.Context, pubsub.Client {
	ctx := context.Background()

    client, err := pubsub.NewClient(ctx, projectID)
    if err != nil {
            logrus.Fatalf("Failed to create client: %v", err)
	}
	return ctx, client
}

//InitTopic returns context, client, and topic for consumtion by InitSubscription or Publish
func InitTopic(projectID , topicName string) context.Context, pubsub.Client, pubsub.Topic {
	ctx, client := initPubSubClient(projectID)
	topic, err := client.CreateTopic(ctx, topicName)
	if err != nil {
		logrus.Fatal("Failed topic init")
	}
	return ctx, client, topic
}

//InitSubscription
func InitSubscription(ctx context.Context, client pubsub.Client, topic pubsub.Topic, subName string) context.Context, pubsub.Subscription {
		// Create a new subscription to the previously created topic
	// with the given name.
	sub, err := client.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{
		Topic:            topic,
		AckDeadline:      10 * time.Second,
		ExpirationPolicy: 25 * time.Hour,
	})
	if err != nil {
		logrus.Fatal("Failed subscription init")
	}
	return ctx, sub
}

func Receive(ctx context.Context, sub pubsub.Subscription, done chan struct{}, dataOut chan json.RawMessage, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case done:
			return
		default:
			err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
				dataOut <- m
				m.Ack()
			})
			if err != context.Canceled {
				logrus.Fatal("Failed receive")
			}
		}
	}
}

func Publish(ctx context.Context, topic pubsub.Topic, done chan struct{}, dataIn chan json.RawMessage, wg *sync.WaitGroup) {
	defer wg.Done()
	defer topic.Stop()
	for {
		select {
		case done:
			return
		case dataPub := <-dataIn:
			res := topic.Publish(ctx, &pubsub.Message{Data: dataPub})
		}
	}
}