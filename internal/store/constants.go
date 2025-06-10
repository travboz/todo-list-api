package store

import "time"

var (
	ExpiryTime        = 5 * time.Minute
	TasksCacheKeyBase = "tasks"
)
