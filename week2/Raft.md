Raft is a concensus algorithm for managing a replicated log. Raft is used to build replicated state machines, since logs entries are eventually guarenteed to be in order on all servers, each state machine can build an identical state. 

## Raft Guarantees
Election Safety:  One leader in each term
Leader Append-Only: Leader only adds new entries
Log Matching: If two logs contain an entry with the same index and term, then all previous entries are identical.
Leader Completeness:  If a log entry is commited, then that entry will be present in the leader's logs, for all higher terms.
State Machine Safety: If a server applys a log entry at a given index, then no other server will apply a different log entry at the index.

## Basics
Raft cluster contains an odd number of servers, where decisions are completed by majority votes. Raft clusters container three roles: leaders, candidates and followers. 

- Leaders are responsible for broacasting new logs, and constantly sending heartbeats
- Candidates are responsible for becoming leaders or stepping down to followers
- Followers are responsible for replying to all RPCs and becoming candidates.

Each term has a unique leader, and terms begin whenever a server becomes a candidate. 

Raft uses only two types of RPCs, RequestVote for election voting and AppendEntries for log replication.

## Leader Elections
Followers use a random election timeout ([150ms, 300ms] typically), if no heartbeats are recieved within this timeout, they convert to a candidate.

Candidates immediately begins an election in a new term (term++ and vote for itself) and sends RequestVote RPCs to all other servers, requesting a vote. Whichever happens first:

1. A majority vote is granted, candidate becomes leader
2. Another leader is dicovered by heartbeats, candidate steps down to follower
3. An election timeout expires ([150ms, 300ms] typically), a new election begins.

Leaders broadcast hearbeats, to prevent new elections.


## Log Replication 
All clients send new server requests to the sole leader. Leaders then broadcast AppendEntries RPCs to all other servers, when the majority of clients acknowledge the AppendEntries RPC, the leader commits the log entry and applies the entry to the state machine. The state machine finally responds to the client request.

A log entry stores the command for the RSM, but also a term number that matches the leaders term. 

Inconsitent logs arise when leaders crash, to solve this, leaders force their logs onto followers *SAFETY. Leaders find the earliest matching log with the follower, deletes all logs on follower after that point, and copies the rest of the log after that point.

Instead of backing off one index at a time, we can back off terms at a time.


## Safety

Only logs from leader's current terms are commited by majority vote, once a log has been commited this way, then all prior entries are good to be over written (Log Matching Property).

## Log Compaction
Raft longs become long, we can use periodic snapshotting to truncate logs. The idea is that the state machine's state is sufficient instead of all the logs. Each server takes its own snapshots periodically to truncate it's logs.

We introduce a new InstallSnapshot RPC, for when the follower requires logs from the leader that has been already snapshotted. 

e.g.

L: | snapshot | entry | entry | ....
F: 