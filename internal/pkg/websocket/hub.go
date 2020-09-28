package websocket

import (
	"math/rand"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type node struct {
	ch     chan Pong
	nodeID string
}

type Hub struct {
	pending      map[string]node
	receivers    map[string]node
	registerLock *sync.Mutex
	lastReceiver int
}

type BroadcastFn func(Pong) int

type RegisteredNodesFn func() []string

type RandomUnicastFn func(Pong) error

func NewHub() *Hub {
	return &Hub{
		receivers:    make(map[string]node),
		pending:      make(map[string]node),
		registerLock: &sync.Mutex{},
		lastReceiver: -1,
	}
}

func (h Hub) Add(ch chan Pong) string {
	id := uuid.New().String()
	h.pending[id] = node{ch: ch}
	return id
}

func (h Hub) Register(internalID, externalID string) {
	temp := h.pending[internalID]
	temp.nodeID = externalID
	h.receivers[internalID] = temp
	delete(h.pending, internalID)
}

func (h Hub) RegisterAtomically(internalID, externalID string) []string {
	h.registerLock.Lock()
	defer h.registerLock.Unlock()
	nodes := h.RegisteredNodes()
	h.Register(internalID, externalID)
	return nodes
}

func (h Hub) Unregister(internalID string) {
	h.registerLock.Lock()
	defer h.registerLock.Unlock()
	delete(h.receivers, internalID)
	delete(h.pending, internalID)
}

func (h Hub) Broadcast(message Pong) int {
	for _, node := range h.receivers {
		node.ch <- message
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
	for _, node := range h.receivers {
		if arrayContains(blacklist, node.nodeID) {
			continue
		}
		node.ch <- message
		sentCount++
		if sentCount == receiveCount {
			return sentCount
		}
	}
	return sentCount
}

func (h *Hub) RandomUnicast(message Pong) error {
	h.registerLock.Lock()
	defer h.registerLock.Unlock()
	receiverNum := rand.Intn(len(h.receivers))
	for receiverNum == h.lastReceiver {
		receiverNum = rand.Intn(len(h.receivers))
	}
	h.lastReceiver = receiverNum
	num := 0
	for _, receiver := range h.receivers {
		if num == receiverNum {
			receiver.ch <- message
			return nil
		}
		num++
	}
	return errors.Errorf("Receiver number (%d) is greater than the number of registered receivers (%d)", receiverNum, len(h.receivers))
}

func (h Hub) RegisteredNodes() (nodes []string) {
	for _, node := range h.receivers {
		nodes = append(nodes, node.nodeID)
	}
	return
}
