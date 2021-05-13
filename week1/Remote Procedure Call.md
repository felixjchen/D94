A Remote Procedure Call (RPC) is when a program executes procedures on a seperate machine by communicating function names, arguments and return values over the network.

Psuedo Code:
1. Client creates a stub, with a function name and argument, marshalls stub 
2. Stub is passed to server
3. Server recieves stub and unmarshalls stub
4. Server computes function with arguments (dispatch), puts return value in stub and marshalls stub
5. Stub returned to client
6. Client unmarshalls stub for return value

Marshalling is the process of transforming data with the intent of moving it, serialization does not necessarily have this intent (e.g. store on disk). Marshalling may involve serialization. 

Failure strategies:
- "At least once" / "best effort" is when a client will retry a stub if the server does not acknowledge with a response stub. This is problematic with non idempotent operations, since the server may perform multiple operations with differing outcomes, but only manages to respond once because of networking issues. (e.g. append(v))
- "At most once" is when the server will execute a stub at most once, one way to guarentee this is to sign each stub with a UUID. If the server sees the same UUID more then once, it responds with the same response stub (caching). 
- "Exactly once" involves unbounded retries, duplicate request detection and fault-tolerant service