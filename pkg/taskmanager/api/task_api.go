package api

import (
	"DistributeX/internal/common"
	"context"
	"log"
	"time"

	"DistributeX/DistributeX/api/taskmanagerpb"
)

type TaskManagerAPI struct {
	taskmanagerpb.UnimplementedTaskManagerServer
	manager common.TaskManagerInterface
}

func NewTaskManagerAPI(manager common.TaskManagerInterface) *TaskManagerAPI {
	return &TaskManagerAPI{manager: manager}
}

func (api *TaskManagerAPI) SendTask(ctx context.Context, req *taskmanagerpb.TaskRequest) (*taskmanagerpb.TaskResponse, error) {
	taskID := api.manager.AddTask(req.Payload, time.Now().Add(10*time.Second))
	log.Printf("Task %s added via API", taskID)
	return &taskmanagerpb.TaskResponse{Status: "Task added successfully"}, nil
}

func (api *TaskManagerAPI) NodeHeartbeat(ctx context.Context, req *taskmanagerpb.HeartbeatRequest) (*taskmanagerpb.HeartbeatResponse, error) {
	response, err := api.manager.NodeHeartbeat(ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}
