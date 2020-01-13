package content

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

var (
	projectID = "events-consumer"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Get make external request to get content
func Get(ctx context.Context, m PubSubMessage) error {
	topicID := "process_content"
	message := fmt.Sprintf("message from get_content '%s'", m.Data)
	fmt.Printf("%s", message)
	err := publish(ctx, topicID, message)

	if err != nil {
		fmt.Printf("can't publish message to topis %s for Get", topicID)
		return err
	}
	return nil
}

// Process receives content and check if changed
func Process(ctx context.Context, m PubSubMessage) error {
	topicID := "send_message"
	message := fmt.Sprintf("message from process_content '%s'", m.Data)
	fmt.Printf("%s", message)
	err := publish(ctx, topicID, message)

	if err != nil {
		fmt.Printf("can't publish message to topis %s for Process", topicID)
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
