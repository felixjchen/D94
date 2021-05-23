package raft

import (
	"log"
	"math/rand"
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

func (rf *Raft) becomeFollower(term int) {
	rf.currentTerm = term
	rf.state = Follower
	rf.votedFor = NoVote
}

func (rf *Raft) getElectionTimeout() time.Duration {
	min := 200
	max := 500
	random_election_timeout := rand.Intn(max-min) + min
	return time.Millisecond * time.Duration(random_election_timeout)
}

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type LogEntry struct {
	Command interface{}
	Term    int
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
	Entries      []LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term    int
	Success bool
}
