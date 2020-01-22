package process_content


import (
"encoding/json"

"github.com/pkg/errors"
)

type IncomingPayload struct {
	Command string `json:"command"`
	Selector string `json:"selector"`
	Content string `json:"content"`
	Name string `json:"name"`
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

