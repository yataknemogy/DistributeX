package executionnode

import (
	"context"
	"log"
	"time"

	"DistributeX/DistributeX/api/executionnodepb"
)

type ExecutionNode struct {
	ID       string
	Address  string
	Load     int
	LastSeen time.Time
}

func (node *ExecutionNode) ExecuteTask(ctx context.Context, req *executionnodepb.TaskRequest) (*executionnodepb.TaskResponse, error) {
	log.Printf("Executing task %s with payload: %s", req.TaskId, req.Payload)
	time.Sleep(2 * time.Second)
	log.Printf("Task %s completed", req.TaskId)
	return &executionnodepb.TaskResponse{Status: "completed"}, nil
}
