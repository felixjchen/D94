package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"
)

type Coordinator struct {
	// Your definitions here.
	mu sync.Mutex
	// m for map, r for reduce
	state   string
	mapped  map[string]bool
	reduced map[string]bool
}

const (
	OK              = "OK"
	ErrUnknownState = "ErrUnknownState"
)

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

func (c *Coordinator) GetTask(args *GetTaskArgs, reply *GetTaskReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Delegate work
	if c.state == "m" {
		// Give first available map task
		for f, assigned := range c.mapped {
			if !assigned {
				reply.Err = OK
				reply.TaskType = "m"
				reply.TaskNumber = f
				return nil
			}
		}

		// No more map tasks, worker must wait for reduce phase
		reply.Err = OK
		reply.TaskType = "w"
		reply.TaskNumber = ""
		return nil

	} else if c.state == "r" {
		// Give first available reduce task
		for i, assigned := range c.reduced {
			if !assigned {
				reply.Err = OK
				reply.TaskType = "r"
				reply.TaskNumber = i
				return nil
			}
		}

		//  No more reduce tasks, worker must wait for completion phase
		reply.Err = OK
		reply.TaskType = "w"
		reply.TaskNumber = ""
		return nil
	}

	// ??
	reply.Err = ErrUnknownState
	reply.TaskType = c.state
	reply.TaskNumber = ""
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.
	c.mu.Lock()
	defer c.mu.Unlock()

	// Is all map and reduce done ?
	mapDone := true
	for _, v := range c.mapped {
		mapDone = mapDone && v
	}

	reduceDone := true
	for _, v := range c.reduced {
		reduceDone = reduceDone && v
	}

	// fmt.Println("mapDone", mapDone, "reduceDone", reduceDone)
	ret = reduceDone && mapDone

	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	// Your code here.
	// INIT
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = "m"
	c.mapped = make(map[string]bool)
	c.reduced = make(map[string]bool)
	for _, f := range files {
		c.mapped[f] = false
	}
	for i := 0; i < nReduce; i++ {
		c.reduced[strconv.Itoa(i)] = false
	}

	c.server()
	return &c
}
