# Error Correcting Codes for Distributed Caching

🗞️ [Report](https://github.com/felixjchen/Distributed-Cache/blob/main/report/report.pdf)

📊 [Benchmarks](https://github.com/felixjchen/Distributed-Cache/tree/main/report/benchmarks)

🦀 [Rust Implementation](https://github.com/felixjchen/Distributed-Cache)

## Motivation
Reduce redundancy in distributed caching by avoiding data replication and using error correction codes. 

## Result

Created two distributed key value stores, using two strategies:
1. Raft replication
2. Error Correcting Codes (Reed Solomon)

### ECC cache uses 3.5 MiB
![](https://user-images.githubusercontent.com/31393977/129127326-b744db92-29ca-4881-8aee-98c308f8b958.png)
### Raft based cache uses 6.0 MiB
![](https://user-images.githubusercontent.com/31393977/129127327-3d3aedab-76d6-4240-8225-d92d7a13cc78.png)

## Idea

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
| 5    | Basic definitions: finite fields, rate, distance & linear codes | [Raft 2D](https://pdos.csail.mit.edu/6.824/labs/lab-raft.html) (MIT 6.824)                                                                                                |
| 6    | parity check matrix, generator matrix & Hamming code                                      | [In memory k/v store with CLI](https://github.com/pingcap/talent-plan/tree/master/courses/rust/projects/project-1) (TP 201)                                               |
| 7    | Reed Solomon codes                                                         | [Persistant k/v store with CLI](https://github.com/pingcap/talent-plan/tree/master/courses/rust/projects/project-2) (TP 201)                                              |
| 8    | Project Planning                       | [Single-threaded, persistant k/v store server, client with custom protocol ](https://github.com/pingcap/talent-plan/tree/master/courses/rust/projects/project-3) (TP 201) |
| 9    | Project  (Raft Cache)                                                                            |                                                                                                                                                                           |
| 10   | Project  (ECC Cache)                                                                            |                                                                                                                                                                           |
| 11   | Project  (Benchmarks)                                                                            |                                                                                                                                                                           |
| 12   | Project  (Report)                                                                            |                                                                                                                                                                           |

## Project

Using Go or Rust, create a distributed cache (k/v store or nosql) that uses error correcting codes to rebuild failed nodes. 


## 
Introuce error
networking benchmarks
timing benchmarks

## Introduction
- intuition : save space while keeping perofrmance
- approach : related works 
- my apparoach: implementation
- expirement : spec of expirements
- results : analyze what we obtained
- conclusion : does it match hypothesiss
- future work : if we have time / money whats next

- dont be shy 
