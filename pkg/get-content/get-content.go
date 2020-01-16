package get_content

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
	processContentTopicID = os.Getenv("PROCESS_CONTENT_TOPIC_ID")
)

// Get make external request to get get-content
func Get(ctx context.Context, m PubSubMessage) error {
	message := fmt.Sprintf("message from get_content '%s'", m.Data)
	fmt.Printf("%s", message)
	err := publish(ctx, processContentTopicID, message)

	if err != nil {
		fmt.Printf("can't publish message to topis %s for Get", processContentTopicID)
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
