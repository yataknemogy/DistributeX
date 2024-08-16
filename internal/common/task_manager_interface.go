package common

import (
	"DistributeX/DistributeX/api/taskmanagerpb"
	"context"
	"time"
)

type TaskManagerInterface interface {
	RegisterNode(id, address string)
	UnregisterNode(id string)
	AddTask(payload interface{}, deadline time.Time) string
	ProcessTasks()
	NodeHeartbeat(ctx context.Context, req *taskmanagerpb.HeartbeatRequest) (*taskmanagerpb.HeartbeatResponse, error)
}
