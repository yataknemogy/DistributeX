package taskmanager

import (
	"log"
	"net"

	"DistributeX/DistributeX/api/executionnodepb"
	"DistributeX/DistributeX/api/taskmanagerpb"
	"DistributeX/internal/common"
	"DistributeX/internal/config"
	"DistributeX/internal/monitoring"
	"DistributeX/internal/storage"
	grpcUtils "DistributeX/pkg/common/grpc"
	"DistributeX/pkg/common/logger"
	"DistributeX/pkg/executionnode/metrics"
	taskManagerAPI "DistributeX/pkg/taskmanager/api"
	"DistributeX/pkg/taskmanager/scheduler"
)

func StartTaskManager() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logger.NewLogger("app.log")
	logger.Info("Logger initialized")

	redisCache := storage.NewRedisCache(cfg.RedisAddress, "", 0)

	distributeFunc := func(task *common.Task, nodes []*common.ExecutionNode) *common.ExecutionNode {
		return scheduler.LeastConnections(nodes)
	}

	taskManager := NewTaskManager(distributeFunc, cfg.NodeTimeout, redisCache)

	taskManagerAPIInstance := taskManagerAPI.NewTaskManagerAPI(taskManager)

	go metrics.CollectMetrics()
	go monitoring.MonitorSystemMetrics()
	go taskManager.ProcessTasks()

	taskManager.RegisterNode("node1", "localhost:50051")
	taskManager.RegisterNode("node2", "localhost:50052")

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpcUtils.CreateGRPCServer()
	executionnodepb.RegisterExecutionNodeServer(grpcServer, &ExecutionNodeServer{})
	taskmanagerpb.RegisterTaskManagerServer(grpcServer, taskManagerAPIInstance)

	// Запуск gRPC сервера
	logger.Info("Starting gRPC server...")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
