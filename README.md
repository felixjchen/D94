# Error Correcting Codes for Distributed Caching

## Project Idea

In distributed caching, how can we leverage error correcting codes for node failure tolerance.

## Motivation

Redis uses a master-slave architecture, which allows for all master nodes to fail, but it requires at least 50% redundant storage.

A Hamming(15,11) code uses ~27% redundent bits and will allow us to correct 1 bit errors.

Reed-Solomon codes will allow us to correct an arbitrary number of error bits, depending on how many parity bits we use. RS(n,k) can correct (n-k)/2 symbols ([CMU](https://www.cs.cmu.edu/~guyb/realworld/reedsolomon/reed_solomon_codes.html)).

## Sources

| Course                                                                                                                                                                    |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| MIT 6.824: [Distributed Systems](https://pdos.csail.mit.edu/6.824/schedule.html)                                                                                          |
| Talent-Plan 201: [Practical Networked Applications in Rust](https://github.com/pingcap/talent-plan/tree/master/courses/rust)                                              |
| University of Washington CSE533: [Error-Correcting Codes: Constructions and Algorithms](https://courses.cs.washington.edu/courses/cse533/06au/)                           |
| University of Buffalo CSE545: [Error Correcting Codes: Combinatorics, Algorithms and Applications](https://cse.buffalo.edu/faculty/atri/courses/coding-theory/spr09.html) |
| Penn State CSE554 (Video Lectures): [Error-Correcting Codes](https://goo.gl/63Kc29) & [Notes](http://www.ee.psu.edu/viveck/EE564_s2016/)                                  |
| Mario Blaum: [A Short Course on Error-Correcting Codes](https://arxiv.org/abs/1908.09903)                                                                                 |

## Schedule

I propose following MIT 6.824: Distributed Systems closely for 4 weeks, followed by studying selective topics in error correcting codes from a few sources. The last 4 weeks will be dedicated to the project.

Material is subject to change, based on project relevance.

| Week | Material                                                                             | Code                                                                                                                                                                      |
| ---- | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1    | Introduction & RPC and Threads (MIT 6.824)                                           | [MapReduce 1](https://pdos.csail.mit.edu/6.824/labs/lab-mr.html) (MIT 6.824)                                                                                              |
| 2    | Raft (MIT 6.824)                                                                     | [Raft 2A](https://pdos.csail.mit.edu/6.824/labs/lab-raft.html) (MIT 6.824)                                                                                                |
| 3    | Cache Consistency & Distributed Transactions (MIT 6.824)                             | [Raft 2B](https://pdos.csail.mit.edu/6.824/labs/lab-raft.html) (MIT 6.824)                                                                                                |
| 4    | MapReduce 1 & Raft 2A 2B, Q&A (MIT 6.824)                                            | [Raft 2C](https://pdos.csail.mit.edu/6.824/labs/lab-raft.html) (MIT 6.824)                                                                                                |
| 5    | Basic definitions: finite fields, rate, distance, linear codes & parity check matrix | [Raft 2D](https://pdos.csail.mit.edu/6.824/labs/lab-raft.html) (MIT 6.824)                                                                                                |
| 6    | Hamming code & Binary Symmetric Channel (BSC)                                        | [In memory k/v store with CLI](https://github.com/pingcap/talent-plan/tree/master/courses/rust/projects/project-1) (TP 201)                                               |
| 7    | Golay code & Bounds on codes                                                         | [Persistant k/v store with CLI](https://github.com/pingcap/talent-plan/tree/master/courses/rust/projects/project-2) (TP 201)                                              |
| 8    | Reed Solomon (decoding)                                                              | [Single-threaded, persistant k/v store server, client with custom protocol ](https://github.com/pingcap/talent-plan/tree/master/courses/rust/projects/project-3) (TP 201) |
| 9    | Project                                                                              |                                                                                                                                                                           |
| 10   | Project                                                                              |                                                                                                                                                                           |
| 11   | Project                                                                              |                                                                                                                                                                           |
| 12   | Project                                                                              |                                                                                                                                                                           |

## Project

Using Go or Rust, create a distributed cache (k/v store or nosql) that uses error correcting codes to rebuild failed nodes. Clients should use a custom TCP protocol.
