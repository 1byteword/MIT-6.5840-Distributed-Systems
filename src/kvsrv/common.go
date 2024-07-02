package kvsrv

type PutAppendArgs struct {
	Key    string
	Value  string
	Op     string // "Put" or "Append"
	Id     int64  // Unique client ID
	SeqNum int    // Request sequence number
}

type PutAppendReply struct {
	WrongLeader bool
	Value       string
}

type GetArgs struct {
	Key    string
	Id     int64 // Unique client ID
	SeqNum int   // Request sequence number
}

type GetReply struct {
	WrongLeader bool
	Value       string
}
