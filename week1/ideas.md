## Lectures

### Distributed System
- Everyone has cloud => everyone can build a distributed system
- What kind of tradeoffs affect consistency, and what do we need to give up for a very consistant model

### MapReduce
- Why not Foldl ? This function is more powershell then map and reduce, but it might be harder to parrallelize....
- MapReduce retries work because both are deterministic functions => Can we create a programming language that runs its programs as a distributed computation ? Feel like some sequential functions will be difficult to compute, but ones like map and reduce can be done in parralel. 

### RPC
- serializing != marshalling, although they are similar
- marshalling has intent for moving data
- serializing may not have this intent, may want to store on disk serialized
- Marshalling may involve serializing the data (JSON encoding)

## Lab
- fault tolerence, 10 seconds for completion seems naive, are there betterways to do this
- a supervisor mental model seemed helpful, how can I rally all workers to accomplish my computation
- intermediate file naming is hard
- Atomic rename on work completion is a good way to solve failed work but half done intermediate fiels
