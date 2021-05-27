## Cache Consistency: Frangipani
  - Distributed File System

  - Cache Coherenece "I see what you write, even though we both have cached copies"
    - I Aquire lock before read/write
    - If another client has lock, they flush data, release lock
    - I aquire lock, and shared disk is updated
  - Atomicity "Our actions are resolved, they don't overwrite eachother" 
    - Lock Server (Replicated for fault tolerence)
    - Must aquire lock before read/write 
  - Crash Recovery "When I crash, you can still use the file system"
    - Write-ahead logging on shared disk
    - locks have leases, and lock server relizes if a client has died
    - Lock server requests another client's "recovery demon" to finish based on write-ahead log

## Distributed Transactions
  - Concurrency Control => 2 Phase Locking
    - 1. Lock when transacation sees variable
    - 2. Release when transacation complete

    - Means our transactions have some serializable order (T1 cannot happen during T2..)
    - Can produce deadlock
    - "Pessimistic", locking causes overhead

  - Atomic Commit => 2 Phase Commit
    - 1. Coordinator sends all participants a PREPARE message
    - 2. Participants reply with Yes/No
    - 3a. If all yes: 
      - Coordinator sends COMMIT, 
    - 3b. else: ABORT
    - 4. Participants COMMIT or ABORT

    - After participants replying YES to prepare must hold lock until COMMIT / ABORT is complete. Coordinator may have crashed but other participants may be committing ..
    - After Coordinator sends COMMIT , it must complete, again because participants may have commited.
    - PREPARE phase is safe to backout
    - During important phases, use logging to recover from failure and continue execution

    - Slow, locking, 2 rounds of messages, disk writes