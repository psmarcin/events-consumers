package content

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

var (
	collectionID        = "web_contents"
	errDocumentNotFound = errors.New("Document not found")
	sendMessageTopicID  = os.Getenv("SEND_MESSAGE_TOPIC_ID")
)

func Process(ctx context.Context, m PubSubMessage) error {
	payload, err := NewProcessPayload(m.Data)
	if err != nil {
		return err
	}

	bgCtx := context.Background()
	firestoreClient, err := firestore.NewClient(bgCtx, projectID)
	if err != nil {
		return fmt.Errorf("%s", errors.Wrap(err, "can't creat firestoreClient for firebase"))
	}

	// Close firestoreClient when done.
	defer firestoreClient.Close()

	// get latest value
	wc, err := getDocument(firestoreClient, bgCtx, collectionID, payload)
	if err == errDocumentNotFound {
		err = addDocument(firestoreClient, bgCtx, collectionID, payload, payload.Content)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[%s] adding document failed", payload.Name))
		}
	}
	if err != nil && err != errDocumentNotFound {
		return errors.Wrap(err, fmt.Sprintf("[%s] get document failed", payload.Name))
	}

	// check if content changed
	if contentChanged := hasContentChanged(wc.Value, payload.Content); contentChanged == false {
		fmt.Printf("[%s] content hasn't change, still %s", payload.Name, payload.Content)
		return nil
	}

	// content has changed so we have to add it to database
	err = addDocument(firestoreClient, bgCtx, collectionID, payload, payload.Content)
	if err != nil {
		return errors.Wrap(err, "adding document failed")
	}

	// create pubsub client
	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	message := fmt.Sprintf("[%s] %s → %s", payload.Name, wc.Value, payload.Content)
	err = publish(pubsubClient, ctx, sendMessageTopicID, message)

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
	payload ProcessPayload,
	value string,
) error {
	_, _, err := client.Collection(collectionID).Add(ctx, map[string]interface{}{
		"command":   payload.Command,
		"selector":  payload.Selector,
		"value":     value,
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
	payload ProcessPayload,
) (WebContent, error) {
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
