package raft

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	Debug     = false
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

// Log helpers
func (rf *Raft) lastLogEntryTerm() int {
	lastLogEntryTerm := rf.lastIncludedTerm
	if len(rf.log) != 0 {
		lastLogEntryTerm = rf.lastLogEntry().Term
	}
	return lastLogEntryTerm
}

func (rf *Raft) lastLogEntryIndex() int {
	lastLogEntryIndex := rf.lastIncludedIndex
	if len(rf.log) != 0 {
		lastLogEntryIndex = rf.lastLogEntry().Index
	}
	return lastLogEntryIndex
}
func (rf *Raft) lastLogEntry() LogEntry {
	return rf.log[len(rf.log)-1]
}
func (rf *Raft) firstLogEntry() LogEntry {
	return rf.log[0]
}
func (rf *Raft) getAdjustedIndex(i int) int {
	offset := rf.log[0].Index
	adjustedIndex := i - offset
	return adjustedIndex
}
func (rf *Raft) logEntry(i int) LogEntry {
	adjustedIndex := rf.getAdjustedIndex(i)
	return rf.log[adjustedIndex]
}

func (rf *Raft) getElectionTimeout() time.Duration {
	min := 200
	max := 500
	random_election_timeout := rand.Intn(max-min) + min
	return time.Millisecond * time.Duration(random_election_timeout)
}

// Utility
func min(x int, y int) int {
	return int(math.Min(float64(x), float64(y)))
}
func max(x int, y int) int {
	return int(math.Max(float64(x), float64(y)))
}

func DP(a ...interface{}) {
	r := ""
	for _, i := range a {
		r += fmt.Sprintf(" %+v ", i)
	}
	DPrintf("%s \n \n", r)
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
	Index   int
}

// func encodeArray(log []LogEntry) []byte {
// 	w := new(bytes.Buffer)
// 	e := labgob.NewEncoder(w)
// 	for _, entry := range log {
// 		iCommand := entry.Command.(int)
// 		e.Encode(iCommand)
// 	}
// 	return w.Bytes()
// }
