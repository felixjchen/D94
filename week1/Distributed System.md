A distributed system is a collection of computers on different networks, which coordinate to achieve a common computation.

Using distributed systems:
- we can utilize physically seperate machines for our tasks
- we can increase computation power by adding more computers (scale)
- we can tolerate faults using replication
- we can physically isolate each server, which provides security

Cloud computing has given everyone access to large amounts of compute, making distributed systems more accesible. 

Challenges:
- Harder to justify correctness
- Quantify performance, overhead for distributed computing, some jobs not parallelizable 
- Fault tolerance, if a node dies, the computation must continue

Considerations:
- Fault tolerance, how can we make our application highly available and recover from failures
- Consistency, how can we guarentee well-defined behavior 
- Performance, do we scale efficiently
- Fault tolerence, consistency, performance clash, we often make tradeoffs in practice
