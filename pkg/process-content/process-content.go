package process_content

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
)

var (
	projectID = "events-consumer"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

var (
	sendMessageTopicID = os.Getenv("SEND_MESSAGE_TOPIC_ID")
)

func Process(ctx context.Context, m PubSubMessage) error {
	err := publish(ctx, sendMessageTopicID, string(m.Data))

	if err != nil {
		fmt.Printf("can't publish message to topis %s for Process", sendMessageTopicID)
		return err
	}

	return nil
}

func publish(ctx context.Context, topicID, message string) error {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	t := client.Topic(topicID)
	t.Publish(ctx, &pubsub.Message{
		Data: []byte(message),
	})
	return nil
}
