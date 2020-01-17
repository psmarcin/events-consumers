package get_content

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Payload struct {
	URL string `json:"url"`
	Selector string `json:"selector"`
}

// New parse raw bytes to Payload struct
func New(raw []byte) (Payload, error) {
	p := Payload{}
	err := json.Unmarshal(raw, &p)
	if err != nil {
		return p, errors.Wrap(err, "parse failed")
	}
	return p, nil
}