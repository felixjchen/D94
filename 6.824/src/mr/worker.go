package mr

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

func get_content(filename string) string {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatalf("cannot open %v", filename)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
	}
	return string(content)
}

//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue, reducef func(string, []string) string) {

	// Your worker implementation here.
	fmt.Println("Worker started!")

	// uncomment to send the Example RPC to the coordinator.
	// CallExample()

	task := get_task()
	for task.TaskType == "m" || task.TaskType == "r" || task.TaskType == "w" {
		fmt.Println(task)
		if task.TaskType == "m" {
			// Compute map
			filename := task.TaskNumber
			content := get_content(filename)
			kva := mapf(filename, string(content))

			// Output result, nReduce buckets, "mr-out-X-Y"
			nReduce := task.NReduce
			intermediate_files := []*os.File{}
			for i := 0; i < nReduce; i++ {
				// NEED BETTER TEMP NAMES AT LEAST UNIQUE
				temp_name := fmt.Sprintf("mr-out-%s-%d-tmp", filename, i)
				temp_file, _ := os.Create(temp_name)
				intermediate_files = append(intermediate_files, temp_file)
			}
			// For each KV, append to respective file
			for _, kv := range kva {
				Y := ihash(kv.Key) % nReduce
				fmt.Fprintf(intermediate_files[Y], "%v %v\n", kv.Key, kv.Value)
			}
			// Atomic Rename
			for _, file := range intermediate_files {
				temp_name := file.Name()
				oname := temp_name[:len(temp_name)-4]
				os.Rename(temp_name, oname)
			}

			// Complete Task
			complete_task(task)

		} else if task.TaskType == "r" {
			// Compute reduce
			glob_pattern := fmt.Sprintf("mr-out-*-%s", task.TaskNumber)
			intermediate_files, err := filepath.Glob(glob_pattern)
			if err != nil {
				log.Fatalf("cannot glob %s", glob_pattern)
			}
			intermediate := []KeyValue{}
			// Add all files for these keys into array
			for _, filename := range intermediate_files {
				content := get_content(filename)
				lines := strings.Split(content, "\n")
				for _, line := range lines {
					// Last line can be empty..
					if line != "" {
						s := strings.Split(line, " ")
						k := s[0]
						v := s[1]
						kv := KeyValue{k, v}
						intermediate = append(intermediate, kv)
					}
				}
			}
			sort.Sort(ByKey(intermediate))

			// NEED BETTER TEMP NAMES AT LEAST UNIQUE
			temp_name := fmt.Sprintf("mr-out-%s-tmp", task.TaskNumber)
			temp_file, _ := os.Create(temp_name)
			defer temp_file.Close()

			i := 0
			for i < len(intermediate) {
				j := i + 1
				for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
					j++
				}
				values := []string{}
				for k := i; k < j; k++ {
					values = append(values, intermediate[k].Value)
				}
				output := reducef(intermediate[i].Key, values)

				// this is the correct format for each line of Reduce output.
				fmt.Fprintf(temp_file, "%v %v\n", intermediate[i].Key, output)
				i = j
			}
			// Output result, "mr-out-X"
			// Atomic Rename
			oname := temp_name[:len(temp_name)-4]
			os.Rename(temp_name, oname)

			// Complete Task
			complete_task(task)
		} else {
			// Wait 5 secs
			time.Sleep(5 * time.Second)
		}

		// Get next task
		task = get_task()
	}

	// Done all

}

//
// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
//
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	call("Coordinator.Example", &args, &reply)

	// reply.Y should be 100.
	fmt.Printf("reply.Y %v\n", reply.Y)
}

func get_task() GetTaskReply {
	args := GetTaskArgs{}

	reply := GetTaskReply{}
	call("Coordinator.GetTask", &args, &reply)

	return reply
}

func complete_task(task GetTaskReply) error {
	args := CompleteTaskArgs{}
	args.TaskType = task.TaskType
	args.TaskNumber = task.TaskNumber

	reply := CompleteTaskReply{}
	call("Coordinator.CompleteTask", &args, &reply)

	return nil
}

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
