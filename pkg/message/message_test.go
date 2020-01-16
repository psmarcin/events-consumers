package message

import (
	"errors"
	"testing"
)

func TestSendMessageToChannel(t *testing.T) {
	err := SendMessageToChannel("","", "")
	if err == nil {
		t.Errorf("Should return error: 'APIKey can't be empty', got: %s, want: %s", err, errors.New("APIKey can't be empty"))
	}

	err = SendMessageToChannel("api","", "")
	if err == nil {
		t.Errorf("Should return error: 'message can't be empty', got: %s, want: %s", err, errors.New("message can't be empty"))
	}

	err = SendMessageToChannel("api","message", "")
	if err == nil {
		t.Errorf("Should return error: 'channelId can't be empty', got: %s, want: %s", err, errors.New("channelId can't be empty"))
	}
}