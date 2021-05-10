package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"
)

type Coordinator struct {
	// Your definitions here.
	mu      sync.Mutex
	nReduce int

	map_assigned    map[string]bool
	reduce_assigned map[string]bool

	map_completed    map[string]bool
	reduce_completed map[string]bool
}

const (
	OK                 = "OK"
	ErrUnknownTaskType = "ErrUnknownTaskType"
)

// Your code here -- RPC handlers for the worker to call.
func isMapTrue(m map[string]bool) bool {
	mapTrue := true
	for _, v := range m {
		mapTrue = mapTrue && v
	}
	return mapTrue
}

func watchMapTask(c *Coordinator, task_number string) {
	time.Sleep(10 * time.Second)

	c.mu.Lock()
	defer c.mu.Unlock()

	// If not commplete, this needs to be reassigned
	if !c.map_completed[task_number] {
		c.map_assigned[task_number] = false
	}
}

func watchReduceTask(c *Coordinator, task_number string) {
	time.Sleep(10 * time.Second)

	c.mu.Lock()
	defer c.mu.Unlock()

	// If not commplete, this needs to be reassigned
	if !c.reduce_completed[task_number] {
		c.reduce_assigned[task_number] = false
	}
}

func (c *Coordinator) CompleteTask(args *CompleteTaskArgs, reply *CompleteTaskReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	task_type := args.TaskType
	task_number := args.TaskNumber

	if task_type == "m" {
		c.map_completed[task_number] = true
		reply.Err = OK
		return nil
	} else if task_type == "r" {
		c.reduce_completed[task_number] = true
		reply.Err = OK
		return nil
	} else {
		reply.Err = ErrUnknownTaskType
		return nil
	}
}

func (c *Coordinator) GetTask(args *GetTaskArgs, reply *GetTaskReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Delegate work
	// Give first available map task
	for i, assigned := range c.map_assigned {
		if !assigned {
			c.map_assigned[i] = true
			reply.Err = OK
			reply.TaskType = "m"
			reply.TaskNumber = i
			reply.NReduce = c.nReduce
			go watchMapTask(c, i)
			return nil
		}
	}

	// No more map tasks, worker must wait for reduce phase
	mapDone := isMapTrue(c.map_completed)
	if !mapDone {
		reply.Err = OK
		reply.TaskType = "w"
		reply.TaskNumber = ""
		reply.NReduce = c.nReduce
		return nil
	}

	// Give first available reduce task
	for i, assigned := range c.reduce_assigned {
		if !assigned {
			c.reduce_assigned[i] = true
			reply.Err = OK
			reply.TaskType = "r"
			reply.TaskNumber = i
			reply.NReduce = c.nReduce
			go watchReduceTask(c, i)
			return nil
		}
	}

	reduceDone := isMapTrue(c.reduce_completed)
	if !reduceDone {
		//  No more reduce tasks, worker must wait for completion phase
		reply.Err = OK
		reply.TaskType = "w"
		reply.TaskNumber = ""
		reply.NReduce = c.nReduce
		return nil
	}

	// Done
	reply.Err = OK
	reply.TaskType = "DONE"
	reply.TaskNumber = ""
	reply.NReduce = c.nReduce
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
	mapDone := isMapTrue(c.map_completed)
	reduceDone := isMapTrue(c.reduce_completed)

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

	c.nReduce = nReduce

	c.map_assigned = make(map[string]bool)
	c.reduce_assigned = make(map[string]bool)
	c.map_completed = make(map[string]bool)
	c.reduce_completed = make(map[string]bool)
	for _, f := range files {
		c.map_assigned[f] = false
		c.map_completed[f] = false
	}
	for i := 0; i < nReduce; i++ {
		j := strconv.Itoa(i)
		c.reduce_assigned[j] = false
		c.reduce_completed[j] = false
	}

	c.server()
	return &c
}
