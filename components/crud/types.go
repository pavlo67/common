package crud

import "time"

type Status struct {
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type Part struct {
	Key   []string
	Count uint64
}
