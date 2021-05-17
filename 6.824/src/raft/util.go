package raft

import (
	"log"
	"time"
)

// Debugging
const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

func (rf *Raft) readLastLeaderRPC() time.Time {
	// Read lastLeaderRPC property, locking appropriately
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.lastLeaderRPC
}

func (rf *Raft) heartbeat() {
	// appendEntries empty all peers
	rf.mu.Lock()
	defer rf.mu.Unlock()

	for i := range rf.peers {
		if i != rf.me {
			args := AppendEntriesArgs{}
			reply := AppendEntriesReply{}
			go rf.sendAppendEntries(i, &args, &reply, []string{})
		}
	}

	go func() {
		// queue another heartbeat if im leader
		time.Sleep(20000)
		_, isleader := rf.GetState()
		if isleader {
			rf.heartbeat()
		}
	}()
}

type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

type RequestVoteReply struct {
	// Your data here (2A).
	Term        int
	VoteGranted bool
}

type AppendEntriesArgs struct {
	Term         int
	LeaderId     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []string
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term    int
	Success bool
}
