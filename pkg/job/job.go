package job

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

var (
	projectID    = "events-consumer"
	collectionID = "jobs"
)

type Job struct {
	Command  string `json:"command"`
	Selector string `json:"selector"`
	Name     string `json:"name"`
}

func (p *Job) Serialize() ([]byte, error) {
	serialized, err := json.Marshal(&p)
	if err != nil {
		return nil, errors.Wrap(err, "can't serialize Job")
	}
	return serialized, err
}

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

var (
	getContentTopicID = os.Getenv("GET_CONTENT_TOPIC_ID")
)

func GetJobs(ctx context.Context, m PubSubMessage) error {
	// Get a Firestore client.
	bgCtx := context.Background()
	client, err := firestore.NewClient(bgCtx, projectID)
	if err != nil {
		return errors.Wrap(err, "can't creat client for firebase")
	}

	// Close client when done.
	defer client.Close()

	jobs, err := getDocuments(client, bgCtx, collectionID)
	if err != nil {
		return errors.Wrap(err, "getting document failed")
	}

	for _, job := range jobs {
		message, err := job.Serialize()
		if err != nil {
			fmt.Printf("can't get message %+v", job)
			continue
		}

		err = publish(ctx, getContentTopicID, message)
		fmt.Printf("published %s\n", message)

		if err != nil {
			fmt.Printf("can't publish message to topis %s for GetJobs", getContentTopicID)
			return err
		}
	}

	fmt.Printf("started %d jobs", len(jobs))

	return nil
}

func getDocuments(
	client *firestore.Client,
	ctx context.Context,
	collectionID string,
) ([]Job, error) {
	var jobs []Job

	query := client.
		Collection(collectionID)

	iter := query.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return jobs, errors.Wrap(err, "can't iterate over next item")
		}
		data := doc.Data()
		jobs = append(jobs, Job{
			Command:  fmt.Sprintf("%s", data["command"]),
			Selector: fmt.Sprintf("%s", data["selector"]),
			Name:     fmt.Sprintf("%s", data["name"]),
		})
	}

	return jobs, nil
}

func publish(ctx context.Context, topicID string, message []byte) error {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	t := client.Topic(topicID)
	t.Publish(ctx, &pubsub.Message{
		Data: message,
	})
	return nil
}
