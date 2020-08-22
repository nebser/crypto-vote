package websocket

import "sync"

type Hub struct {
	receivers map[string]chan Pong
	owner     string
	lock      *sync.WaitGroup
}

func NewHub(owner string) Hub {
	lock := sync.WaitGroup{}
	lock.Add(1)
	return Hub{
		receivers: make(map[string]chan Pong),
		owner:     owner,
		lock:      &lock,
	}
}

func (h Hub) Register(nodeID string, ch chan Pong) {
	h.receivers[nodeID] = ch
}

func (h Hub) RegisterAtomically(nodeID string, ch chan Pong) []string {
	h.lock.Wait()
	defer h.lock.Done()
	nodes := h.RegisteredNodes()
	h.Register(nodeID, ch)
	return nodes
}

func (h Hub) Unregister(nodeID string) {
	delete(h.receivers, nodeID)
}

func (h Hub) Broadcast(message Pong) int {
	for _, ch := range h.receivers {
		ch <- message
	}
	return len(h.receivers)
}

func arrayContains(array []string, target string) bool {
	for _, elem := range array {
		if elem == target {
			return true
		}
	}
	return false
}

func (h Hub) Multicast(message Pong, receiveCount int, blacklist []string) int {
	sentCount := 0
	for node, ch := range h.receivers {
		if arrayContains(blacklist, node) {
			continue
		}
		ch <- message
		sentCount++
		if sentCount == receiveCount {
			return sentCount
		}
	}
	return sentCount
}

func (h Hub) RegisteredNodes() (nodes []string) {
	for node := range h.receivers {
		nodes = append(nodes, node)
	}
	return
}
