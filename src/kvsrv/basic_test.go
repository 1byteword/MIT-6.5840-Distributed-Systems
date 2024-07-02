package kvsrv

import (
	"testing"
)

func TestKVServerBasicOperations(t *testing.T) {
	// Create a server
	server := StartKVServer()

	// Test Put
	t.Run("TestPut", func(t *testing.T) {
		putArgs := &PutArgs{Key: "key1", Value: "value1"}
		putReply := &PutReply{}
		server.Put(putArgs, putReply)

		getArgs := &GetArgs{Key: "key1"}
		getReply := &GetReply{}
		server.Get(getArgs, getReply)
		if getReply.Value != "value1" {
			t.Errorf("Expected value1, got %s", getReply.Value)
		}
	})

	// Test Get (non-existent key)
	t.Run("TestGetNonExistent", func(t *testing.T) {
		getArgs := &GetArgs{Key: "non-existent"}
		getReply := &GetReply{}
		server.Get(getArgs, getReply)
		if getReply.Value != "" {
			t.Errorf("Expected empty string, got %s", getReply.Value)
		}
	})

	// Test Append
	t.Run("TestAppend", func(t *testing.T) {
		appendArgs1 := &AppendArgs{Key: "key2", Value: "Hello"}
		appendReply1 := &AppendReply{}
		server.Append(appendArgs1, appendReply1)
		if appendReply1.Value != "" {
			t.Errorf("Expected empty string, got %s", appendReply1.Value)
		}

		appendArgs2 := &AppendArgs{Key: "key2", Value: " World"}
		appendReply2 := &AppendReply{}
		server.Append(appendArgs2, appendReply2)
		if appendReply2.Value != "Hello" {
			t.Errorf("Expected 'Hello', got %s", appendReply2.Value)
		}

		getArgs := &GetArgs{Key: "key2"}
		getReply := &GetReply{}
		server.Get(getArgs, getReply)
		if getReply.Value != "Hello World" {
			t.Errorf("Expected 'Hello World', got %s", getReply.Value)
		}
	})

	// Test Put after Append
	t.Run("TestPutAfterAppend", func(t *testing.T) {
		putArgs := &PutArgs{Key: "key2", Value: "New Value"}
		putReply := &PutReply{}
		server.Put(putArgs, putReply)

		getArgs := &GetArgs{Key: "key2"}
		getReply := &GetReply{}
		server.Get(getArgs, getReply)
		if getReply.Value != "New Value" {
			t.Errorf("Expected 'New Value', got %s", getReply.Value)
		}
	})

	// Test multiple operations
	t.Run("TestMultipleOperations", func(t *testing.T) {
		putArgs := &PutArgs{Key: "key3", Value: "First"}
		putReply := &PutReply{}
		server.Put(putArgs, putReply)

		appendArgs1 := &AppendArgs{Key: "key3", Value: " Second"}
		appendReply1 := &AppendReply{}
		server.Append(appendArgs1, appendReply1)

		appendArgs2 := &AppendArgs{Key: "key3", Value: " Third"}
		appendReply2 := &AppendReply{}
		server.Append(appendArgs2, appendReply2)

		getArgs := &GetArgs{Key: "key3"}
		getReply := &GetReply{}
		server.Get(getArgs, getReply)
		if getReply.Value != "First Second Third" {
			t.Errorf("Expected 'First Second Third', got %s", getReply.Value)
		}
	})
}
