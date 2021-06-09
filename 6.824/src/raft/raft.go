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
	// My Volatile states
	apply     chan ApplyMsg
	heartbeat chan bool
	state     string

	// Persistent state on all servers
	currentTerm       int
	votedFor          int
	snapshot          []byte
	lastIncludedIndex int
	lastIncludedTerm  int
	log               []LogEntry

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
	e.Encode(rf.currentTerm)
	e.Encode(rf.votedFor)
	e.Encode(rf.lastIncludedIndex)
	e.Encode(rf.lastIncludedTerm)
	e.Encode(rf.snapshot)
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
	var snapshot []byte
	var lastIncludedIndex int
	var lastIncludedTerm int
	var log []LogEntry
	if d.Decode(&currentTerm) != nil ||
		d.Decode(&votedFor) != nil ||
		d.Decode(&lastIncludedIndex) != nil ||
		d.Decode(&lastIncludedTerm) != nil ||
		d.Decode(&snapshot) != nil ||
		d.Decode(&log) != nil {
		// error...
	} else {
		rf.currentTerm = currentTerm
		rf.votedFor = votedFor
		rf.lastIncludedIndex = lastIncludedIndex
		rf.lastIncludedTerm = lastIncludedTerm
		rf.snapshot = snapshot
		rf.log = log
	}
}

//
// A service wants to switch to snapshot.  Only do so if Raft hasn't
// have more recent info since it communicate the snapshot on applyCh.
//
func (rf *Raft) CondInstallSnapshot(lastIncludedTerm int, lastIncludedIndex int, snapshot []byte) bool {
	// rf gets garbage collected on occasion ....
	if rf == nil {
		return false
	}

	rf.mu.Lock()
	defer rf.mu.Unlock()
	DP("CondInstallSnapshot", rf.me, rf.snapshot, rf.log, rf.commitIndex, lastIncludedIndex, snapshot)

	// If snapshot breaks into this servers commited logs... return false
	// We don't need this snapshot anymore, our logs are caught up
	if rf.commitIndex > lastIncludedIndex {
		return false
	}

	// 5. Save snapshot file, discard any existing or partial snapshot with a smaller index
	rf.snapshot = snapshot
	rf.lastIncludedIndex = lastIncludedIndex
	rf.lastIncludedTerm = lastIncludedTerm

	// 6. If existing log entry has same index and term as snapshot’s last included entry, retain log entries following it and reply
	found := false
	retainFromIndex := 0
	for i := 0; i < len(rf.log); i++ {
		if rf.log[i].Index == lastIncludedIndex && rf.log[i].Term == lastIncludedTerm {
			found = true
			retainFromIndex = i + 1
			break
		}
	}
	if found {
		rf.log = rf.log[retainFromIndex:]
	} else {
		// 7. Discard the entire log
		rf.log = []LogEntry{}
	}
	rf.persist()
	return true
}

// the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (2D).
	go func() {
		rf.mu.Lock()
		defer rf.mu.Unlock()
		DP("Snapshot", rf.me, rf.snapshot, rf.log, index, snapshot)

		// This snapshot is late .... its inside our current snapshot
		if index < rf.lastIncludedIndex {
			return
		}

		// set snapshot
		rf.snapshot = snapshot
		rf.lastIncludedIndex = index
		rf.lastIncludedTerm = rf.logEntry(index).Term
		// trim log
		rf.log = rf.log[rf.getAdjustedIndex(index)+1:]
		rf.persist()
	}()
}

func (rf *Raft) InstallSnapshot(args *InstallSnapshotArgs, reply *InstallSnapshotReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	DP("InstallSnapshot", rf.me, rf.snapshot, rf.log, args)

	// Rules for Servers (§5.1)
	if args.Term > rf.currentTerm {
		rf.becomeFollower(args.Term)
	}

	// 1. Reply immediately if term < currentTerm
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		return
	}

	// 8. Reset state machine using snapshot contents (and load
	// snapshot’s cluster configuration)
	message := ApplyMsg{
		SnapshotValid: true,
		Snapshot:      args.Snapshot,
		SnapshotIndex: args.LastIncludedIndex,
		SnapshotTerm:  args.LastIncludedTerm,
	}
	rf.apply <- message

	reply.Term = rf.currentTerm
	return
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

	// up-to-date =
	// a) If the logs have last entries with different terms, then the log with the later term is more up-to-date.
	// b) If the logs end with the same term, then whichever log is longer is more up-to-date.
	atLeastUpToDate := args.LastLogTerm > rf.lastLogEntryTerm() || (args.LastLogTerm == rf.lastLogEntryTerm() && args.LastLogIndex >= rf.lastLogEntryIndex())

	// 2. If votedFor is null or candidateId, and candidate’s log is at least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)
	// "first come first serve"
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
	DP("AppendEntries", rf.me, rf.lastIncludedIndex, rf.snapshot, rf.log, args)

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

	// viable leader found, update heartbeat
	rf.heartbeat <- true

	prevLogEntryTerm := rf.lastIncludedTerm
	if rf.lastIncludedIndex < args.PrevLogIndex && args.PrevLogIndex <= rf.lastLogEntryIndex() {
		prevLogEntryTerm = rf.logEntry(args.PrevLogIndex).Term
	}

	// 2. Reply false if log doesn’t contain an entry at prevLogIndex whose term matches prevLogTerm (§5.3)
	if !(args.PrevLogIndex <= rf.lastLogEntryIndex()) || prevLogEntryTerm != args.PrevLogTerm {
		reply.Term = rf.currentTerm
		reply.Success = false

		// we find the next index to try by walking back terms
		start := min(args.PrevLogIndex, rf.lastLogEntryIndex())
		end := start
		if len(rf.log) > 0 {
			end = rf.firstLogEntry().Index
		}

		for i := start; i > end; i-- {
			if rf.logEntry(i).Term != rf.logEntry(start).Term {
				reply.NextIndex = i
				break
			}
		}
		return
	}

	// 3. If an existing entry conflicts with a new one (same index but different terms), delete the existing entry and all that follow it (§5.3)
	// 4. Append any new entries not already in the log
	if len(rf.log) == 0 {
		rf.log = append([]LogEntry{}, args.Entries...)
	} else {
		rf.log = append(rf.log[:rf.getAdjustedIndex(args.PrevLogIndex)+1], args.Entries...)
	}
	rf.persist()

	// 5. If leaderCommit > commitIndex, set commitIndex = min(leaderCommit, index of last new entry)
	if args.LeaderCommit > rf.commitIndex {
		rf.commitIndex = min(args.LeaderCommit, rf.lastLogEntryIndex())
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

func (rf *Raft) sendInstallSnapshot(server int, args *InstallSnapshotArgs, reply *InstallSnapshotReply) bool {
	ok := rf.peers[server].Call("Raft.InstallSnapshot", args, reply)

	rf.mu.Lock()
	defer rf.mu.Unlock()
	// Rules for Servers (§5.1)
	if reply.Term > rf.currentTerm {
		rf.becomeFollower(reply.Term)
	}
	return ok
}

// send request votes, in parallel goroutines
func (rf *Raft) sendRequestVotes(electionWon chan bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	args := &RequestVoteArgs{
		Term:         rf.currentTerm,
		CandidateId:  rf.me,
		LastLogIndex: rf.lastLogEntryIndex(),
		LastLogTerm:  rf.lastLogEntryTerm(),
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
						nextIndex := 0
						if len(rf.log) > 0 {
							nextIndex = rf.lastLogEntry().Index + 1
						}
						for i := 0; i < len(rf.peers); i++ {
							rf.nextIndex = append(rf.nextIndex, nextIndex)
							rf.matchIndex = append(rf.matchIndex, 0)
						}
						rf.state = Leader

						electionWon <- true
					}
				}
			}(i)
		}
	}
}

func (rf *Raft) sendHeartbeat() {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	DP("sendHeartbeat", rf.me, rf.log, rf.nextIndex)
	// replicate logs onto followers
	// once majority replicated, commit
	for i := 0; i < len(rf.peers); i++ {
		if i != rf.me {
			go func(peer int) {
				rf.mu.Lock()
				// If peer needs entries from AFTER my snapshot... can appendEntries
				appendEntries := rf.lastIncludedIndex < rf.nextIndex[peer]

				if appendEntries {
					prevLogTerm := rf.lastIncludedTerm
					if rf.lastIncludedIndex+1 != rf.nextIndex[peer] {
						prevLogTerm = rf.logEntry(rf.nextIndex[peer] - 1).Term
					}

					entries := []LogEntry{}
					// Only slice if not heartbeat and log exists
					if len(rf.log) != 0 && rf.getAdjustedIndex(rf.nextIndex[peer]) < len(rf.log) {
						// clone because shared mem leads to race conditions
						entries = append(entries, rf.log[rf.getAdjustedIndex(rf.nextIndex[peer]):]...)
					}

					reply := &AppendEntriesReply{}
					args := &AppendEntriesArgs{
						Term:         rf.currentTerm,
						LeaderId:     rf.me,
						PrevLogIndex: rf.nextIndex[peer] - 1,
						PrevLogTerm:  prevLogTerm,
						Entries:      entries,
						LeaderCommit: rf.commitIndex,
					}
					newNextIndex := rf.nextIndex[peer] + len(args.Entries)
					newMatchIndex := rf.nextIndex[peer] + len(args.Entries) - 1
					rf.mu.Unlock()

					rf.sendAppendEntries(peer, args, reply)

					rf.mu.Lock()
					defer rf.mu.Unlock()
					if reply.Success {
						// Peer has been caught up!
						rf.nextIndex[peer] = newNextIndex
						rf.matchIndex[peer] = newMatchIndex
					} else {
						// backoff
						rf.nextIndex[peer] = max(reply.NextIndex, 1)
					}

					// If there exists an N such that N > commitIndex, a majority of matchIndex[i] ≥ N, and log[N].term == currentTerm: set commitIndex = N (§5.3, §5.4)
					if len(rf.log) > 0 {
						for n := rf.lastLogEntry().Index; n > rf.commitIndex && rf.lastLogEntry().Term == rf.currentTerm; n-- {
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
					}
				} else {
					// We need to send snapshot
					reply := &InstallSnapshotReply{}
					args := &InstallSnapshotArgs{
						Term:              rf.currentTerm,
						LeaderId:          rf.me,
						LastIncludedIndex: rf.lastIncludedIndex,
						LastIncludedTerm:  rf.lastIncludedTerm,
						Snapshot:          append([]byte{}, rf.snapshot...),
					}
					newNextIndex := rf.lastIncludedIndex + 1
					newMatchIndex := rf.lastIncludedIndex
					rf.mu.Unlock()

					rf.sendInstallSnapshot(peer, args, reply)

					rf.mu.Lock()
					defer rf.mu.Unlock()
					rf.nextIndex[peer] = newNextIndex
					rf.matchIndex[peer] = newMatchIndex
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
	index := rf.lastIncludedIndex + 1
	if len(rf.log) > 0 {
		index = rf.lastLogEntry().Index + 1
	}
	term := rf.currentTerm
	isLeader := rf.state == Leader

	if isLeader {
		// If command received from client: append entry to local log, respond after entry applied to state machine (§5.3)
		newEntry := LogEntry{
			Command: command,
			Term:    term,
			Index:   index,
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

		// We dont want to apply in snapshot...
		if rf.lastIncludedIndex < rf.lastApplied {
			message := ApplyMsg{
				CommandValid: true,
				CommandIndex: rf.lastApplied,
				Command:      rf.logEntry(rf.lastApplied).Command,
			}
			rf.apply <- message
		}
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

			// begin election
			electionWon := make(chan bool)
			go rf.sendRequestVotes(electionWon)

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
	rf.lastIncludedIndex = 0
	rf.lastIncludedTerm = 0
	rf.snapshot = []byte{}

	// everyone gets an empty entry... this is for the uptodate
	emptyEntry := LogEntry{
		Term:  0,
		Index: 0,
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
