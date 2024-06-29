package kvsrv

import (
	"log"
	"sync"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type KVServer struct {
	mu sync.Mutex

	data map[string]string
	// Your definitions here.
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()

	value, exists := kv.data[args.Key]
	if exists {
		reply.Value = value
	} else {
		reply.Value = ""
	}
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()

	reply.Value = kv.data[args.Key]
	kv.data[args.Key] = args.Value
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()

	oldValue, exists := kv.data[args.Key]
	if !exists {
		oldValue = ""
	}

	kv.data[args.Key] = oldValue + args.Value
	reply.Value = oldValue
}

func StartKVServer() *KVServer {
	kv := new(KVServer)
	kv.data = make(map[string]string)

	return kv
}
