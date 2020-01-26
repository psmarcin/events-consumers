package process_content

import "time"

type WebContent struct {
	Url       string
	Selector  string
	Value     string
	CreatedAt time.Time
}
