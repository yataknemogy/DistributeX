package common

import (
	"time"
)

type Task struct {
	ID       string
	Payload  interface{}
	Created  time.Time
	Deadline time.Time
}

type ExecutionNode struct {
	ID       string
	Address  string
	Weight   int
	Load     int
	LastSeen time.Time
}

type Worker struct {
	ID string
}

type Config struct {
	RedisAddress   string
	GRPCPort       string
	NodeTimeout    time.Duration
	KafkaBroker    string
	KafkaTaskTopic string
	KafkaLogTopic  string
}
