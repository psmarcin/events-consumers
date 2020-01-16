package message

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	telegramAPIKey = os.Getenv("TELEGRAM_API_KEY")
	telegramChannelID = os.Getenv("TELEGRAM_CHANNEL_ID")
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Send receive PubSubMessage and process it
func Send(ctx context.Context, m PubSubMessage) error {
	message := string(m.Data)
	err := SendMessageToChannel(telegramAPIKey, message, telegramChannelID)
	if err != nil {
		return err
	}

	fmt.Printf("Message: %s, sent to %s", message, telegramChannelID)
	return nil

}

func SendMessageToChannel(APIKey, message, channelId string) error {
	if APIKey == "" {
		return errors.New("APIKey can't be empty")
	}
	if message == "" {
		return errors.New("message can't be empty")
	}
	if channelId == "" {
		return errors.New("channelId can't be empty")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", APIKey, channelId, message)

	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(res)
	fmt.Println(string(body))

	return nil
}