package types

import "time"

// Job represent a transcoding job to be done by a worker
type Job struct {
	ID     string    `json:"id"`
	Args   []string  `json:"args"`
	Expiry time.Time `json:"expiry"`
}
