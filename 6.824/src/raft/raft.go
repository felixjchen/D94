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

	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	//	"6.824/labgob"
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
	state       string
	currentTerm int
	votedFor    int
	log         []LogEntry

	commitIndex int
	lastApplied int
	lastAppend  time.Time

	nextIndex  []int
	matchIndex []int
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (2A).
	rf.mu.Lock()
	term = rf.currentTerm
	isleader = (rf.state == Leader)
	rf.mu.Unlock()

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

	// 5.1
	if args.Term > rf.currentTerm {
		rf.state = Follower
		rf.currentTerm = args.Term
	}

	// 5.1
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return
	}

	// If I haven't voted and candidate is at least up to date... I will vote for them
	// up-to-date = candidate is my term or higher, candidate has my logs or more
	if rf.votedFor == NoVote || rf.votedFor == args.CandidateId {

		rf.votedFor = args.CandidateId
		reply.Term = rf.currentTerm
		reply.VoteGranted = true

		fmt.Println("RequestedVote:", args.CandidateId, "from", rf.me)
		return
	}

	// no vote
	reply.Term = rf.currentTerm
	reply.VoteGranted = false
	return
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// 5.1
	if args.Term > rf.currentTerm {
		rf.state = Follower
		rf.currentTerm = args.Term
	}

	// 5.1
	if args.Term < rf.currentTerm {
		rf.state = Follower
		rf.currentTerm = args.Term

		reply.Term = args.Term
		reply.Success = false
		return
	}

	// Handle AppendEntries
	rf.lastAppend = time.Now()

	reply.Term = rf.currentTerm
	reply.Success = true
	return

}

func (rf *Raft) sendRequestVote(server int, reply *RequestVoteReply) bool {
	args := RequestVoteArgs{}
	rf.mu.Lock()
	args.Term = rf.currentTerm
	args.CandidateId = rf.me
	rf.mu.Unlock()

	ok := rf.peers[server].Call("Raft.RequestVote", &args, reply)
	// 5.1
	rf.mu.Lock()
	if reply.Term > rf.currentTerm {
		rf.state = Follower
		rf.currentTerm = args.Term
	}
	rf.mu.Unlock()

	return ok
}
func (rf *Raft) sendAppendEntries(server int, reply *AppendEntriesReply, entries []LogEntry) bool {
	args := AppendEntriesArgs{}
	rf.mu.Lock()
	args.Term = rf.currentTerm
	args.LeaderId = rf.me
	args.Entries = entries
	args.LeaderCommit = rf.commitIndex
	rf.mu.Unlock()

	ok := rf.peers[server].Call("Raft.AppendEntries", &args, reply)
	// 5.1
	rf.mu.Lock()
	if reply.Term > rf.currentTerm {
		rf.state = Follower
		rf.currentTerm = args.Term
	}
	rf.mu.Unlock()
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
func (rf *Raft) election_ticker() {
	// election timer
	for rf.killed() == false {

		// Your code here to check if a leader election should
		// be started and to randomize sleeping time using
		rf.mu.Lock()
		prev_append := rf.lastAppend
		rf.mu.Unlock()

		// Random sleep
		min := 300
		max := 500
		random_election_timeout := rand.Intn(max-min) + min
		time.Sleep(time.Millisecond * time.Duration(random_election_timeout))

		rf.mu.Lock()
		thisAppend := rf.lastAppend
		isfollower := rf.state == Follower
		rf.mu.Unlock()

		// trigger election
		if isfollower && prev_append == thisAppend {
			rf.election()
		}
	}
}

func (rf *Raft) election() {
	// convert to candidate
	rf.mu.Lock()
	rf.currentTerm++
	rf.votedFor = rf.me
	rf.state = Candidate

	me := rf.me
	peer_count := len(rf.peers)
	prev_append := rf.lastAppend
	fmt.Println("Term:", rf.currentTerm, "Peer:", me, "Begins Election ")
	rf.mu.Unlock()

	// begin election
	// get votes, in parallel goroutines
	votes := 1
	votes_mu := &sync.Mutex{}
	votes_wg := sync.WaitGroup{}
	for i := 0; i < peer_count; i++ {
		if i != me {
			votes_wg.Add(1)
			go func(peer int) {
				reply := RequestVoteReply{}
				ok := rf.sendRequestVote(peer, &reply)

				if ok && reply.VoteGranted {
					votes_mu.Lock()
					votes++
					votes_mu.Unlock()
				}
				votes_wg.Done()
			}(i)
		}
	}

	votes_wg.Wait()

	rf.mu.Lock()
	fmt.Println("Term:", rf.currentTerm, "Peer:", me, "Election Votes: ", votes)
	this_append := rf.lastAppend
	rf.mu.Unlock()

	// become leader
	if votes > peer_count/2 {
		fmt.Println("Term:", rf.currentTerm, "Peer:", me, "Becomes Leader")
		rf.mu.Lock()
		rf.state = Leader
		rf.mu.Unlock()

		// send empty append entries heartbeat job
		go func() {
			rf.mu.Lock()
			state := rf.state
			rf.mu.Unlock()
			for state == Leader {
				for i := 0; i < peer_count; i++ {
					if i != me {
						go func(peer int) {
							reply := AppendEntriesReply{}
							rf.sendAppendEntries(peer, &reply, []LogEntry{})

							rf.mu.Lock()
							if reply.Term > rf.currentTerm {
								rf.state = Follower
								rf.currentTerm = reply.Term
							}
							rf.mu.Unlock()
						}(i)
					}
				}

				time.Sleep(200 * time.Millisecond)
				rf.mu.Lock()
				state = rf.state
				rf.mu.Unlock()
			}
		}()

		return
	}

	fmt.Println(me, "FAILS ELECTION")

	// become follower, if a valid append came through during election
	if prev_append != this_append {
		rf.mu.Lock()
		rf.state = Follower
		rf.votedFor = NoVote
		rf.mu.Unlock()
		return
	}
	rf.mu.Lock()
	rf.state = Follower
	rf.votedFor = NoVote
	rf.mu.Unlock()
	return

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
	rf.lastAppend = time.Now()

	rf.currentTerm = 1
	rf.votedFor = NoVote
	rf.log = []LogEntry{}

	rf.commitIndex = 0
	rf.lastApplied = 0

	rf.nextIndex = []int{}
	rf.matchIndex = []int{}
	for i := 0; i < len(rf.peers); i++ {
		rf.nextIndex = append(rf.nextIndex, 1)
		rf.matchIndex = append(rf.matchIndex, 0)
	}

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// start ticker goroutine to start elections
	go rf.election_ticker()

	return rf
}
