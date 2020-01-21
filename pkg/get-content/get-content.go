package get_content

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
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
	// transform payload to struc
	payload, err := NewIncomingPayload(m.Data)
	if  err != nil {
		return err
	}

	// make http request
	body, err := getRequest(payload.URL)
	if err != nil{
		return errors.Wrap(err, "can't get page " + payload.URL)
	}

	// select content from body
	value, err := getContent(body, payload.Selector)
	if err != nil{
		return errors.Wrap(err, "can't get value for selector " + payload.Selector)
	}

	// log
	message := fmt.Sprintf("Page %s has changed value to %s", payload.URL, value)
	fmt.Printf("%s", message)

	outgoingPayload := OutgoingPayload{
		URL:      payload.URL,
		Selector: payload.Selector,
		Content:  value,
	}

	// serialize payload
	serialized, err := outgoingPayload.Serialize()
	if err !=nil{
		return err
	}

	// publish event
	err = publish(ctx, processContentTopicID, serialized)
	if err != nil {
		return errors.Wrap(err, "can't publish message")
	}

	return nil
}

func getRequest(url string) (*html.Node, error){
	client := resty.New()

	resp, err := client.R().
		SetDoNotParseResponse(true).
		Get(url)

	if err != nil {
		return nil, errors.Wrap(err, "can't get " + url)
	}

	if resp.IsError() {
		return nil, errors.New(fmt.Sprintf("request failed with status code %s", resp.Status()))
	}

	nodes, err := html.Parse(resp.RawBody())
	if err != nil {
		return nil, errors.Wrap(err, "can't parse html")
	}

	return nodes, nil
}

func getContent(body *html.Node, selector string) (string, error){
	doc := goquery.NewDocumentFromNode(body)

	value := doc.Find(selector).First().Text()
	cleanedValue := strings.TrimSpace(value)
	return cleanedValue, nil
}

func publish(ctx context.Context, topicID string, message []byte) error {
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


