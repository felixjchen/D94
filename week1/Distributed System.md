A distributed system is a collection of computers on different networks, which coordinate to achieve a common computation.

Using distributed systems:
- we can utilize physically seperate machines for our tasks
- we can increase computation power by adding more computers (scale horizontally)
- we can tolerate faults using replication
- we can physically isolate each server, which provides security

Cloud computing has given everyone access to large amounts of compute making distributed systems more accesible. 

Some challenges for distrbuted computing are:
- Harder to justify correctness
- Quantify performance, overhead for distributed computing, some nodes might have to wait etc..
- Fault tolerance, if a node dies, the computation must continue

Some topics to consider in distributed computing are:
- Fault tolerance, how can we make our application available always and recover from failures (replication)
- Consistency, how can we create well-defined behavior 
- Performance, do we scale efficiently 
- Fault tolerence, consistency, performance are enemies, we often make tradeoffs in practice
