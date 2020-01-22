package process_content

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

var (
	projectID = "events-consumer"
	collectionID = "web_contents"
	errDocumentNotFound = errors.New("Document not found")
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

var (
	sendMessageTopicID = os.Getenv("SEND_MESSAGE_TOPIC_ID")
)

func Process(ctx context.Context, m PubSubMessage) error {
	payload, err := NewIncomingPayload(m.Data)
	if  err != nil {
		return err
	}

	// Get a Firestore client.
	bgCtx := context.Background()
	client, err := firestore.NewClient(bgCtx, projectID)
	if err != nil {
		return errors.Wrap(err, "can't creat client for firebase")
	}

	// Close client when done.
	defer client.Close()

	wc, err := getDocument(client, bgCtx, collectionID, payload)
	if err == errDocumentNotFound {
		err = addDocument(client, bgCtx, collectionID, payload, payload.Content)
		if err != nil {
			return errors.Wrap(err, "adding document failed")
		}
	}
	if err != nil && err != errDocumentNotFound {
		return errors.Wrap(err, "adding document failed")
	}

	contentChanged := hasContentChanged(wc.Value, payload.Content)

	if contentChanged == false {
		fmt.Printf("content hasn't change, still %s", payload.Content)
		return nil
	}

	err = addDocument(client, bgCtx, collectionID, payload, payload.Content)
	if err != nil {
		return errors.Wrap(err, "adding document failed")
	}

	message:= fmt.Sprintf("Content changes on page %s, was: %s, now: %s", payload.Command, wc.Value ,payload.Content)
	err = publish(ctx, sendMessageTopicID, message)

	if err != nil {
		fmt.Printf("can't publish message to topis %s for Process", sendMessageTopicID)
		return err
	}

	return nil
}

func addDocument(
	client *firestore.Client,
	ctx context.Context,
	collectionID string,
	payload IncomingPayload,
	value string,
	) error {
	_, _, err := client.Collection(collectionID).Add(ctx, map[string]interface{}{
		"command": payload.Command,
		"selector":  payload.Selector,
		"value":  value,
		"createdAt": time.Now(),
	})
	if err != nil {
		return errors.Wrap(err, "can't add new document")
	}

	return nil
}

func getDocument(
	client *firestore.Client,
	ctx context.Context,
	collectionID string,
	payload IncomingPayload) (WebContent, error) {
	wc := WebContent{}

	query := client.
		Collection(collectionID).
		Where("command", "==", payload.Command).
		Where("selector", "==", payload.Selector).
		OrderBy("createdAt", firestore.Desc).
		Limit(1)

	iter := query.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return wc, errors.Wrap(err, "can't iterate over next item")
		}
		data := doc.Data()
		wc.Selector = fmt.Sprintf("%s", data["selector"])
		wc.Url = fmt.Sprintf("%s", data["url"])
		wc.Value = fmt.Sprintf("%s", data["value"])
		return wc, nil
	}

	return wc, errDocumentNotFound
}

func hasContentChanged(contentA, contentB string) bool {
	return contentA != contentB
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
