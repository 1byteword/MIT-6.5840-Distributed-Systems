package kvsrv

type PutArgs struct {
	Key   string
	Value string
}

type PutReply struct {
	Value string
}

type AppendArgs struct {
	Key   string
	Value string
}

type AppendReply struct {
	Value string
}

type GetArgs struct {
	Key string
}

type GetReply struct {
	Value string
}
