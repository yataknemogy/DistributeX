package scheduler

import "DistributeX/internal/common"

func WeightedRoundRobin(nodes []*common.ExecutionNode) func() *common.ExecutionNode {
	var currentIndex int
	var currentWeight int

	gcd := func(a, b int) int {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}

	gcdOfWeights := nodes[0].Weight
	for _, node := range nodes {
		gcdOfWeights = gcd(gcdOfWeights, node.Weight)
	}

	maxWeight := 0
	for _, node := range nodes {
		if node.Weight > maxWeight {
			maxWeight = node.Weight
		}
	}

	return func() *common.ExecutionNode {
		for {
			currentIndex = (currentIndex + 1) % len(nodes)
			if currentIndex == 0 {
				currentWeight = currentWeight - gcdOfWeights
				if currentWeight <= 0 {
					currentWeight = maxWeight
				}
			}

			if nodes[currentIndex].Weight >= currentWeight {
				return nodes[currentIndex]
			}
		}
	}
}

func LeastConnections(nodes []*common.ExecutionNode) *common.ExecutionNode {
	var selectedNode *common.ExecutionNode
	minConnections := int(^uint(0) >> 1)

	for _, node := range nodes {
		if node.Load < minConnections {
			selectedNode = node
			minConnections = node.Load
		}
	}

	return selectedNode
}
