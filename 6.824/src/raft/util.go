package raft

import (
	"fmt"
	"log"
	"math"
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
	rf.persist()
}

func min(x int, y int) int {
	return int(math.Min(float64(x), float64(y)))
}

func max(x int, y int) int {
	return int(math.Max(float64(x), float64(y)))
}

func (rf *Raft) getRandomElectionTimeout() time.Duration {
	min := 200
	max := 500
	random_election_timeout := rand.Intn(max-min) + min
	return time.Millisecond * time.Duration(random_election_timeout)
}

func DP(a ...interface{}) {
	r := ""
	for _, i := range a {
		r += fmt.Sprintf(" %+v ", i)
	}
	DPrintf("%s", r)
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
	Term      int
	NextIndex int
	Success   bool
}
