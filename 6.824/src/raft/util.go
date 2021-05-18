package raft

import (
	"log"
	"time"
)

// Debugging
const (
	Debug     = true
	Follower  = "follower"
	Candidate = "candidate"
	Leader    = "leader"
	NoVote    = -1
)

func getEpoch() int64 {
	return time.Now().UnixNano()
}

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type LogEntry struct {
	Command string
	Term    int
}

type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term        int
	CandidateId int
	// LastLogIndex int
	// LastLogTerm  int
}

type RequestVoteReply struct {
	// Your data here (2A).
	Term        int
	VoteGranted bool
}

type AppendEntriesArgs struct {
	Term     int
	LeaderId int
	// PrevLogIndex int
	// PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term    int
	Success bool
}
