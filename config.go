package main

import (
	"time"
)

const (
	RedisAddr         = "localhost:6379"
	RetryInterval     = 5 * time.Second
	DefaultMaxRequests = 10
	VIPMaxRequests     = 20
	AdminMaxRequests   = 50
	SlidingWindowSize  = time.Minute
	LogEnabled         = true
	NotifyService      = "Telegram"
)
