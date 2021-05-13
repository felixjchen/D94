MapReduce is a distributed algorithm, where workers perform a map followed by a reduce function. MapReduce abstracts away all details of distributed computing, and a programmer will only need to implement a map and reduce function. 

Psuedo Code:
1. Split input into M files
2. Workers apply map for each input file
	1. Split intermediate results to N files, using a hash function on the key % N  (Ultimately NxM intermediate results)
3. Wait for all maps to complete
4. Workers apply reduce for each intermediate file
5. Concat all reduce outputs for output

MapReduce is fault tolerant because it reassigns map/reduce tasks that are taking too long (stragglers). Since map and reduce are deterministic functions, recomputed tasks are equally correct. 

MapReduce scales well because additional machines can run map or reduce tasks in parallel. 

Some limitations include:
- Limited computation (must be expressed in a map and one reduce)
- No real-time processing
- No iteration
- Master is single point of failure, if it dies the entire computation must restart