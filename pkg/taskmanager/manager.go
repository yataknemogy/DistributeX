package taskmanager

import (
	"context"
	"log"
	"sync"
	"time"

	"DistributeX/DistributeX/api/executionnodepb"
	"DistributeX/DistributeX/api/taskmanagerpb"
	"DistributeX/internal/common"
	"DistributeX/internal/storage"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// TaskManager implements the TaskManagerInterface, which is enforced by this line.
var _ common.TaskManagerInterface = (*TaskManager)(nil)

type TaskManager struct {
	nodes       map[string]*common.ExecutionNode
	tasksQueue  chan *common.Task
	results     map[string]interface{}
	mu          sync.Mutex
	distribute  func(*common.Task, []*common.ExecutionNode) *common.ExecutionNode
	nodeTimeout time.Duration
	cache       *storage.RedisCache
}

type ExecutionNodeServer struct {
	executionnodepb.UnimplementedExecutionNodeServer
}

func NewTaskManager(distribute func(*common.Task, []*common.ExecutionNode) *common.ExecutionNode, nodeTimeout time.Duration, cache *storage.RedisCache) *TaskManager {
	return &TaskManager{
		nodes:       make(map[string]*common.ExecutionNode),
		tasksQueue:  make(chan *common.Task, 100),
		results:     make(map[string]interface{}),
		distribute:  distribute,
		nodeTimeout: nodeTimeout,
		cache:       cache,
	}
}

func (tm *TaskManager) RegisterNode(id, address string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.nodes[id] = &common.ExecutionNode{
		ID:       id,
		Address:  address,
		LastSeen: time.Now(),
	}
	log.Printf("Node %s registered with address %s", id, address)
}

func (tm *TaskManager) UnregisterNode(id string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	delete(tm.nodes, id)
	log.Printf("Node %s unregistered", id)
}

func (tm *TaskManager) AddTask(payload interface{}, deadline time.Time) string {
	task := &common.Task{
		ID:       uuid.New().String(),
		Payload:  payload,
		Created:  time.Now(),
		Deadline: deadline,
	}

	if result, found := tm.GetCachedResult(task.ID); found {
		log.Printf("Task %s found in cache, returning cached result", task.ID)
		tm.mu.Lock()
		tm.results[task.ID] = result
		tm.mu.Unlock()
		return task.ID
	}

	tm.tasksQueue <- task
	log.Printf("Task %s added to queue", task.ID)

	go tm.processTask(task)

	return task.ID
}

func (tm *TaskManager) CacheResult(taskID string, result string) {
	err := tm.cache.Set(taskID, result, time.Hour)
	if err != nil {
		log.Printf("Failed to cache result for task %s: %v", taskID, err)
	}
}

func (tm *TaskManager) GetCachedResult(taskID string) (string, bool) {
	val, err := tm.cache.Get(taskID)
	if err == redis.Nil {
		return "", false
	} else if err != nil {
		log.Printf("Failed to get cached result for task %s: %v", taskID, err)
		return "", false
	}
	return val, true
}

func (tm *TaskManager) ProcessTasks() {
	go func() {
		for task := range tm.tasksQueue {
			tm.processTask(task)
		}
	}()
}

func (tm *TaskManager) processTask(task *common.Task) {
	tm.mu.Lock()
	nodes := tm.getActiveNodes()
	tm.mu.Unlock()

	if len(nodes) == 0 {
		log.Printf("No available nodes to process task %s", task.ID)
		return
	}

	node := tm.distribute(task, nodes)
	tm.sendTaskToNode(task, node)
}

func (tm *TaskManager) sendTaskToNode(task *common.Task, node *common.ExecutionNode) {
	log.Printf("Sending task %s to node %s", task.ID, node.ID)
	node.Load++

	conn, err := grpc.Dial(node.Address, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to node %s: %v", node.ID, err)
		return
	}
	defer conn.Close()

	client := executionnodepb.NewExecutionNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.ExecuteTask(ctx, &executionnodepb.TaskRequest{
		TaskId:  task.ID,
		Payload: task.Payload.(string),
	})

	if err != nil {
		log.Printf("Failed to execute task %s on node %s: %v", task.ID, node.ID, err)
		return
	}

	log.Printf("Task %s executed successfully on node %s", task.ID, node.ID)
	tm.CacheResult(task.ID, response.Status)
}

func (tm *TaskManager) getActiveNodes() []*common.ExecutionNode {
	var activeNodes []*common.ExecutionNode
	now := time.Now()

	for _, node := range tm.nodes {
		if now.Sub(node.LastSeen) < tm.nodeTimeout {
			activeNodes = append(activeNodes, node)
		}
	}
	return activeNodes
}

func (tm *TaskManager) NodeHeartbeat(ctx context.Context, req *taskmanagerpb.HeartbeatRequest) (*taskmanagerpb.HeartbeatResponse, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	node, exists := tm.nodes[req.NodeId]
	if exists {
		node.LastSeen = time.Now()
		log.Printf("Heartbeat received from node %s", req.NodeId)
		return &taskmanagerpb.HeartbeatResponse{Status: "active"}, nil
	}

	log.Printf("Unknown node %s sent heartbeat", req.NodeId)
	return &taskmanagerpb.HeartbeatResponse{Status: "unknown"}, nil
}
