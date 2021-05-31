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
	"bytes"

	"sync"
	"sync/atomic"
	"time"

	"6.824/labgob"
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
	// My Volatile
	apply     chan ApplyMsg
	heartbeat chan bool
	state     string

	// Persistent state on all servers
	currentTerm int
	votedFor    int
	log         []LogEntry

	// Volatile state on all servers
	commitIndex int
	lastApplied int

	// Volatile state on leaders
	nextIndex  []int
	matchIndex []int
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	// Your code here (2A).
	rf.mu.Lock()
	defer rf.mu.Unlock()

	term := rf.currentTerm
	isleader := rf.state == Leader

	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	w := new(bytes.Buffer)
	e := labgob.NewEncoder(w)
	// We serialize only three important states...
	e.Encode(rf.currentTerm)
	e.Encode(rf.votedFor)
	e.Encode(rf.log)
	data := w.Bytes()
	rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	r := bytes.NewBuffer(data)
	d := labgob.NewDecoder(r)
	// Deserialize 3 importatnt states...
	var currentTerm int
	var votedFor int
	var log []LogEntry
	if d.Decode(&currentTerm) != nil ||
		d.Decode(&votedFor) != nil ||
		d.Decode(&log) != nil {
		// error...
	} else {
		rf.currentTerm = currentTerm
		rf.votedFor = votedFor
		rf.log = log
	}
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

	// Rules for Servers (§5.1)
	if args.Term > rf.currentTerm {
		rf.becomeFollower(args.Term)
	}

	// 1. Reply false if term < currentTerm (§5.1)
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return
	}

	// 2. If votedFor is null or candidateId, and candidate’s log is at least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)

	// "first come first serve"

	// up-to-date =
	// a) If the logs have last entries with different terms, then the log with the later term is more up-to-date.
	// b) If the logs end with the same term, then whichever log is longer is more up-to-date.
	atLeastUpToDate := args.LastLogTerm > rf.log[len(rf.log)-1].Term || (args.LastLogTerm == rf.log[len(rf.log)-1].Term && args.LastLogIndex >= len(rf.log)-1)

	if (rf.votedFor == NoVote || rf.votedFor == args.CandidateId) && atLeastUpToDate {
		rf.votedFor = args.CandidateId
		rf.persist()
		reply.Term = rf.currentTerm
		reply.VoteGranted = true
	} else {
		// no vote
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
	}
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// Rules for Servers (§5.1)
	if args.Term > rf.currentTerm {
		rf.becomeFollower(args.Term)
	}

	// 1. Reply false if term < currentTerm (§5.1)
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.Success = false
		return
	}

	// heartbeat
	rf.heartbeat <- true

	// 2. Reply false if log doesn’t contain an entry at prevLogIndex whose term matches prevLogTerm (§5.3)
	if !(args.PrevLogIndex < len(rf.log)) || rf.log[args.PrevLogIndex].Term != args.PrevLogTerm {
		reply.Term = rf.currentTerm
		reply.Success = false
		// we find the next index to try by walking back terms
		start := min(args.PrevLogIndex, len(rf.log)-1)
		conflictTerm := rf.log[start].Term
		for i := start; i > 0; i-- {
			if rf.log[i].Term != conflictTerm {
				reply.NextIndex = i
				break
			}
		}
		return
	}

	// 3. If an existing entry conflicts with a new one (same index but different terms), delete the existing entry and all that follow it (§5.3)
	// 4. Append any new entries not already in the log
	rf.log = append(rf.log[:args.PrevLogIndex+1], args.Entries...)
	rf.persist()

	// 5. If leaderCommit > commitIndex, set commitIndex = min(leaderCommit, index of last new entry)
	if args.LeaderCommit > rf.commitIndex {
		rf.commitIndex = min(args.LeaderCommit, len(rf.log)-1)
	}

	reply.Term = rf.currentTerm
	reply.Success = true
	return
}

func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)

	rf.mu.Lock()
	defer rf.mu.Unlock()

	// Rules for Servers (§5.1)
	if reply.Term > rf.currentTerm {
		rf.becomeFollower(reply.Term)
	}

	return ok
}

func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)

	rf.mu.Lock()
	defer rf.mu.Unlock()

	// Rules for Servers (§5.1)
	if reply.Term > rf.currentTerm {
		rf.becomeFollower(reply.Term)
	}
	return ok
}

func (rf *Raft) sendHeartbeat() {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// replicate logs onto followers
	// once majority replicated, commit
	for i := 0; i < len(rf.peers); i++ {
		if i != rf.me {
			go func(peer int) {
				rf.mu.Lock()
				reply := &AppendEntriesReply{}
				args := &AppendEntriesArgs{
					Term:         rf.currentTerm,
					LeaderId:     rf.me,
					PrevLogIndex: rf.nextIndex[peer] - 1,
					PrevLogTerm:  rf.log[rf.nextIndex[peer]-1].Term,
					Entries:      append([]LogEntry{}, rf.log[rf.nextIndex[peer]:]...),
					LeaderCommit: rf.commitIndex,
				}

				// if len(rf.log) -1 >= rf.nextIndex[peer] {
				// 	args.PrevLogIndex = rf.nextIndex[peer]
				// }
				rf.mu.Unlock()

				rf.sendAppendEntries(peer, args, reply)

				rf.mu.Lock()
				defer rf.mu.Unlock()
				if reply.Success {
					// Peer has been caught up!
					// rf.nextIndex[peer] = len(rf.log)
					// rf.matchIndex[peer] = len(rf.log) - 1
					rf.nextIndex[peer] = rf.nextIndex[peer] + len(args.Entries)
					rf.matchIndex[peer] = rf.nextIndex[peer] + len(args.Entries) - 1
				} else {
					// backoff
					rf.nextIndex[peer] = max(reply.NextIndex, 1)
				}

				// If there exists an N such that N > commitIndex, a majority of matchIndex[i] ≥ N, and log[N].term == currentTerm: set commitIndex = N (§5.3, §5.4)
				for n := len(rf.log) - 1; n > rf.commitIndex && rf.log[n].Term == rf.currentTerm; n-- {
					replicas := 1
					for i := 0; i < len(rf.peers); i++ {
						if i != rf.me {
							if rf.matchIndex[i] >= n {
								replicas++
							}
						}
					}
					if replicas > len(rf.peers)/2 {
						rf.commitIndex = n
						break
					}
				}
			}(i)
		}
	}
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
	rf.mu.Lock()
	defer rf.mu.Unlock()
	index := len(rf.log)
	term := rf.currentTerm
	isLeader := rf.state == Leader

	if isLeader {
		// If command received from client: append entry to local log, respond after entry applied to state machine (§5.3)
		newEntry := LogEntry{
			Command: command,
			Term:    term,
		}
		rf.log = append(rf.log, newEntry)
		rf.persist()
	}

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

// If commitIndex > lastApplied: increment lastApplied, apply log[lastApplied] to state machine (§5.3)
func (rf *Raft) applyCheck() {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// Rules for Servers (§5.3)
	for rf.commitIndex > rf.lastApplied {
		rf.lastApplied++
		message := ApplyMsg{
			CommandValid: true,
			CommandIndex: rf.lastApplied,
			Command:      rf.log[rf.lastApplied].Command,
		}
		rf.apply <- message
	}
}

// The ticker go routine starts a new election if this peer hasn't received
// heartsbeats recently.
func (rf *Raft) ticker() {
	for !rf.killed() {
		rf.mu.Lock()
		state := rf.state
		rf.mu.Unlock()

		// Apply commited logs
		go rf.applyCheck()

		switch state {

		case Follower:
			select {
			// viable leader exists
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
			rf.persist()
			rf.mu.Unlock()

			electionWon := make(chan bool)

			// begin election
			// get votes, in parallel goroutines
			go func() {
				rf.mu.Lock()
				defer rf.mu.Unlock()
				args := &RequestVoteArgs{
					Term:         rf.currentTerm,
					CandidateId:  rf.me,
					LastLogIndex: len(rf.log) - 1,
					LastLogTerm:  rf.log[len(rf.log)-1].Term,
				}

				votes := 1
				votes_mu := sync.Mutex{}
				// broadcast, get votes from everyone!
				for i := 0; i < len(rf.peers); i++ {
					if i != rf.me {
						go func(peer int) {
							reply := &RequestVoteReply{}
							rf.sendRequestVote(peer, args, reply)

							rf.mu.Lock()
							defer rf.mu.Unlock()
							if rf.state == Candidate && rf.currentTerm == args.Term && reply.VoteGranted {
								votes_mu.Lock()
								defer votes_mu.Unlock()
								votes++
								if votes > len(rf.peers)/2 {
									// Volatile state on leaders:
									rf.nextIndex = []int{}
									rf.matchIndex = []int{}
									for i := 0; i < len(rf.peers); i++ {
										rf.nextIndex = append(rf.nextIndex, len(rf.log))
										rf.matchIndex = append(rf.matchIndex, 0)
									}
									rf.state = Leader

									electionWon <- true
								}
							}
						}(i)
					}
				}
			}()

			select {
			// majority votes recieved
			case <-electionWon:
			// viable leader exists
			case <-rf.heartbeat:
			// election timeout
			case <-time.After(rf.getElectionTimeout()):
			}

		case Leader:

			go rf.sendHeartbeat()
			time.Sleep(time.Millisecond * 60)
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
	// Volatile - mine
	rf.state = Follower
	rf.apply = applyCh
	rf.heartbeat = make(chan bool, 3)

	// Persistent state on all servers
	rf.currentTerm = 1
	rf.votedFor = NoVote
	emptyEntry := LogEntry{
		Term: 0,
	}
	rf.log = []LogEntry{emptyEntry}

	// Volatile state on all servers:
	rf.commitIndex = 0
	rf.lastApplied = 0

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// start ticker goroutine to start elections
	go rf.ticker()

	return rf
}
