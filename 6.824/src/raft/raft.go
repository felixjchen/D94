package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	//	"bytes"

	"sync"
	"sync/atomic"
	"time"

	"6.824/labrpc"
)

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in part 2D you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh, but set CommandValid to false for these
// other uses.
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int

	// For 2D:
	SnapshotValid bool
	Snapshot      []byte
	SnapshotTerm  int
	SnapshotIndex int
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	// Persistent state on all servers
	currentTerm int
	votedFor    int
	log         []LogEntry

	// Volatile state on all servers
	// commitIndex int
	// lastApplied int

	heartbeat   chan bool
	electionWin chan bool
	votes       int
	state       string

	// Volatile state on leaders
	// nextIndex  []int
	// matchIndex []int
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (2A).
	rf.mu.Lock()
	defer rf.mu.Unlock()

	term = rf.currentTerm
	isleader = rf.state == Leader

	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

//
// A service wants to switch to snapshot.  Only do so if Raft hasn't
// have more recent info since it communicate the snapshot on applyCh.
//
func (rf *Raft) CondInstallSnapshot(lastIncludedTerm int, lastIncludedIndex int, snapshot []byte) bool {

	// Your code here (2D).

	return true
}

// the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (2D).

}

func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// 1) Term < currentTerm
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return
	}

	// Rules for Servers (§5.1)
	if args.Term > rf.currentTerm {
		rf.becomeFollower(args.Term)
	}

	// 2) Conditions for giving vote
	// If I haven't voted and candidate is at least up to date... I will vote for them
	// up-to-date = candidate is my term or higher, candidate has my logs or more
	if rf.votedFor == NoVote || rf.votedFor == args.CandidateId {
		rf.votedFor = args.CandidateId
		reply.Term = rf.currentTerm
		reply.VoteGranted = true
	} else {
		// no vote
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
	}
}

func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {

	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	// 5.1
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.state != Candidate || rf.currentTerm != args.Term {
		return ok
	}

	// Rules for Servers (§5.1)
	if reply.Term > rf.currentTerm {
		rf.becomeFollower(reply.Term)
	}

	if reply.VoteGranted {
		rf.votes++

		if rf.votes > len(rf.peers)/2 {
			rf.state = Leader
			rf.electionWin <- true
		}
	}
	return ok

}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// 1) term < currentTerm reply false
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.Success = false
		return
	}

	// Rules for Servers (§5.1)
	if args.Term > rf.currentTerm {
		rf.becomeFollower(args.Term)
	}

	// 2) Handle AppendEntries
	reply.Term = rf.currentTerm
	reply.Success = true

	// heartbeat
	rf.heartbeat <- true
}

func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	// 5.1
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if !ok || rf.state != Leader || args.Term != rf.currentTerm {
		return ok
	}
	if reply.Term > rf.currentTerm {
		rf.state = Follower
		rf.currentTerm = reply.Term
		rf.votedFor = NoVote
	}
	return ok
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (2B).

	return index, term, isLeader
}

//
// the tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. your code can use killed() to
// check whether Kill() has been called. the use of atomic avoids the
// need for a lock.
//
// the issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. any goroutine with a long-running loop
// should call killed() to check whether it should stop.
//
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

// The ticker go routine starts a new election if this peer hasn't received
// heartsbeats recently.
func (rf *Raft) ticker() {
	// election timer
	for {
		rf.mu.Lock()
		state := rf.state
		rf.mu.Unlock()

		switch state {
		case Follower:
			select {
			case <-rf.heartbeat:
			// convert to candidate
			case <-time.After(rf.getElectionTimeout()):
				rf.mu.Lock()
				rf.state = Candidate
				rf.mu.Unlock()
			}

		case Candidate:
			rf.mu.Lock()
			rf.currentTerm++
			rf.votedFor = rf.me
			rf.votes = 1
			rf.mu.Unlock()

			// begin election
			// get votes, in parallel goroutines
			go func() {
				rf.mu.Lock()
				defer rf.mu.Unlock()
				args := RequestVoteArgs{
					Term:        rf.currentTerm,
					CandidateId: rf.me,
				}
				for i := 0; i < len(rf.peers); i++ {
					if i != rf.me {
						reply := RequestVoteReply{}
						go rf.sendRequestVote(i, &args, &reply)
					}
				}
			}()

			select {
			case <-rf.heartbeat:
				rf.mu.Lock()
				rf.state = Follower
				rf.mu.Unlock()
			case <-rf.electionWin:
				rf.mu.Lock()
				rf.state = Leader
				rf.mu.Unlock()
			case <-time.After(rf.getElectionTimeout()):
			}

		case Leader:

			go func() {
				rf.mu.Lock()
				defer rf.mu.Unlock()

				args := AppendEntriesArgs{
					Term:     rf.currentTerm,
					LeaderId: rf.me,
					Entries:  []LogEntry{},
				}
				for i := 0; i < len(rf.peers); i++ {
					if i != rf.me {
						reply := AppendEntriesReply{}
						go rf.sendAppendEntries(i, &args, &reply)
					}
				}
			}()

			time.Sleep(time.Millisecond * 150)
		}

	}
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}

	rf.mu.Lock()
	defer rf.mu.Unlock()

	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).
	rf.state = Follower
	rf.votes = 0
	rf.heartbeat = make(chan bool, 100)
	rf.electionWin = make(chan bool, 100)

	rf.currentTerm = 1
	rf.votedFor = NoVote
	rf.log = []LogEntry{}

	// rf.commitIndex = 0
	// rf.lastApplied = 0

	// rf.nextIndex = []int{}
	// rf.matchIndex = []int{}
	// for i := 0; i < len(rf.peers); i++ {
	// 	rf.nextIndex = append(rf.nextIndex, 1)
	// 	rf.matchIndex = append(rf.matchIndex, 0)
	// }

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// start ticker goroutine to start elections
	go rf.ticker()

	return rf
}
