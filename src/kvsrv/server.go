package kvsrv

import (
	"sync"
	"time"
)

const (
	maxLastAppliedSize = 1000
	cleanupInterval    = 1 * time.Minute
)

type OperationResult struct {
	SeqNum    int
	Value     string
	Timestamp time.Time
}

type KVServer struct {
	mu          sync.Mutex
	data        map[string]string
	lastApplied map[int64]OperationResult
	cleanupDone chan bool
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if lastOp, exists := kv.lastApplied[args.Id]; exists && args.SeqNum <= lastOp.SeqNum {
		reply.Value = lastOp.Value
		reply.WrongLeader = false
		return
	}

	value, exists := kv.data[args.Key]
	if exists {
		reply.Value = value
	} else {
		reply.Value = ""
	}

	kv.lastApplied[args.Id] = OperationResult{SeqNum: args.SeqNum, Value: reply.Value, Timestamp: time.Now()}
	reply.WrongLeader = false
	kv.cleanupIfNeeded()
}

func (kv *KVServer) PutAppend(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if lastOp, exists := kv.lastApplied[args.Id]; exists && args.SeqNum <= lastOp.SeqNum {
		reply.Value = lastOp.Value
		reply.WrongLeader = false
		return
	}

	oldValue := kv.data[args.Key]

	if args.Op == "Put" {
		kv.data[args.Key] = args.Value
	} else if args.Op == "Append" {
		kv.data[args.Key] += args.Value
	}

	kv.lastApplied[args.Id] = OperationResult{SeqNum: args.SeqNum, Value: oldValue, Timestamp: time.Now()}
	reply.WrongLeader = false
	reply.Value = oldValue
	kv.cleanupIfNeeded()
}

func (kv *KVServer) cleanupIfNeeded() {
	if len(kv.lastApplied) > maxLastAppliedSize {
		kv.cleanup()
	}
}

func (kv *KVServer) cleanup() {
	now := time.Now()
	for clientId, op := range kv.lastApplied {
		if now.Sub(op.Timestamp) > cleanupInterval {
			delete(kv.lastApplied, clientId)
		}
	}
}

func (kv *KVServer) periodicCleanup() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			kv.mu.Lock()
			kv.cleanup()
			kv.mu.Unlock()
		case <-kv.cleanupDone:
			return
		}
	}
}

func StartKVServer() *KVServer {
	kv := new(KVServer)
	kv.data = make(map[string]string)
	kv.lastApplied = make(map[int64]OperationResult)
	kv.cleanupDone = make(chan bool)
	go kv.periodicCleanup()
	return kv
}

func (kv *KVServer) Kill() {
	close(kv.cleanupDone)
}
