## Replication
- RAFT provides tolerence through replication... this is a bit against error correcting codes, where we allow a percentage of data to fail, and use a percentage as parity. 
- perhaps there is a nice middle ground with replication + parity 

## Implementing RAFT
- select, chan pattern is powerful for implementing timeouts
- still leader based, wonder if there a p2p model with same guarenetees 

