package store

import "time"

var (
	CacheExpiryTime   = 5 * time.Minute
	TasksCacheKeyBase = "tasks"
)
