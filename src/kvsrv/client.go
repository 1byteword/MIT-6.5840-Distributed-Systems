package kvsrv

import (
	"crypto/rand"
	"math/big"
	"sync"
	"time"

	"6.5840/labrpc"
)

type Clerk struct {
	server *labrpc.ClientEnd
	mu     sync.Mutex
	id     int64
	seqNum int
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(server *labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.server = server
	ck.id = nrand()
	ck.seqNum = 0
	return ck
}

func (ck *Clerk) Get(key string) string {
	ck.mu.Lock()
	seqNum := ck.seqNum
	ck.seqNum++
	ck.mu.Unlock()

	args := GetArgs{
		Key:    key,
		Id:     ck.id,
		SeqNum: seqNum,
	}
	for {
		var reply GetReply
		ok := ck.server.Call("KVServer.Get", &args, &reply)
		if ok && !reply.WrongLeader {
			return reply.Value
		}
		time.Sleep(time.Duration(50+nrand()%50) * time.Millisecond)
	}
}

func (ck *Clerk) PutAppend(key string, value string, op string) string {
	ck.mu.Lock()
	seqNum := ck.seqNum
	ck.seqNum++
	ck.mu.Unlock()

	args := PutAppendArgs{
		Key:    key,
		Value:  value,
		Op:     op,
		Id:     ck.id,
		SeqNum: seqNum,
	}
	for {
		var reply PutAppendReply
		ok := ck.server.Call("KVServer.PutAppend", &args, &reply)
		if ok && !reply.WrongLeader {
			return reply.Value
		}
		time.Sleep(time.Duration(50+nrand()%50) * time.Millisecond)
	}
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}

func (ck *Clerk) Append(key string, value string) string {
	return ck.PutAppend(key, value, "Append")
}
