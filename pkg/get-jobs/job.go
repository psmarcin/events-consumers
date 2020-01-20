package get_jobs

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Job struct {
	Url      string `json:"url"`
	Selector string `json:"selector"`
}


func (p *Job) Serialize() ([]byte, error) {
	serialized, err := json.Marshal(&p)
	if err != nil {
		return nil, errors.Wrap(err, "can't serialize Job")
	}
	return serialized, err
}
