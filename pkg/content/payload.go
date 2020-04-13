package content

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type IncomingPayload struct {
	Command  string `json:"command"`
	Selector string `json:"selector"`
	Name     string `json:"name"`
}

// NewIncomingPayload parse raw bytes to IncomingPayload struct
func NewIncomingPayload(raw []byte) (IncomingPayload, error) {
	p := IncomingPayload{}
	err := json.Unmarshal(raw, &p)
	if err != nil {
		return p, errors.Wrap(err, "parse failed")
	}
	return p, nil
}

type OutgoingPayload struct {
	Command  string `json:"command"`
	Selector string `json:"selector"`
	Content  string `json:"content"`
	Name     string `json:"name"`
}

func (p *OutgoingPayload) Serialize() (string, error) {
	serialized, err := json.Marshal(&p)
	if err != nil {
		return "", errors.Wrap(err, "can't serialize OutgoingPayload")
	}
	return string(serialized), err
}


type ProcessPayload struct {
	Command  string `json:"command"`
	Selector string `json:"selector"`
	Content  string `json:"content"`
	Name     string `json:"name"`
}

// NewIncomingPayload parse raw bytes to IncomingPayload struct
func NewProcessPayload(raw []byte) (ProcessPayload, error) {
	p := ProcessPayload{}
	err := json.Unmarshal(raw, &p)
	if err != nil {
		return p, errors.Wrap(err, "parse failed")
	}
	return p, nil
}


type WebContent struct {
	Url       string
	Selector  string
	Value     string
	CreatedAt time.Time
}

