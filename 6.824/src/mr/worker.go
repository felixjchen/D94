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

// Open file with filename, return file contents
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
	task := call_get_task()

	// while task type is map, reduce or wait (DONE==task_type exits)
	for task.TaskType == "m" || task.TaskType == "r" || task.TaskType == "w" {
		fmt.Println(task)

		if task.TaskType == "m" {
			// Get file content, compute map
			ifile := task.TaskNumber
			content := get_content(ifile)
			kva := mapf(ifile, string(content))

			// Output result, nReduce buckets, "map-out-X-Y"
			nReduce := task.NReduce
			intermediate_files := []*os.File{}
			for i := 0; i < nReduce; i++ {

				// trim path to just input file name
				t := filepath.Base(ifile)
				t = strings.Split(t, ".txt")[0]

				temp_name := fmt.Sprintf("map-out-%s-%d.temp*", t, i)
				temp_file, err := ioutil.TempFile("", temp_name)
				defer temp_file.Close()
				if err != nil {
					log.Fatalf("error creating temp file %s", temp_name)
				}
				intermediate_files = append(intermediate_files, temp_file)
			}
			// For each KV, compute hash and append to respective file
			for _, kv := range kva {
				Y := ihash(kv.Key) % nReduce
				fmt.Fprintf(intermediate_files[Y], "%v %v\n", kv.Key, kv.Value)
			}
			// Atomic Rename
			for _, file := range intermediate_files {
				temp_name := file.Name()
				file_name := filepath.Base(temp_name)
				file_name = strings.Split(file_name, ".temp")[0]
				os.Rename(temp_name, file_name)
			}

			// Complete Task
			call_complete_task(task)

		} else if task.TaskType == "r" {

			// Get all input files for this reduce task, bring KV into memory
			glob_pattern := fmt.Sprintf("map-out-*-%s", task.TaskNumber)
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

			// We group all keys together, then call reduce on k, [v], where v all share same key
			sort.Sort(ByKey(intermediate))

			// Output reduce to temp file, then atomic rename
			temp_name := fmt.Sprintf("mr-out-%s.temp*", task.TaskNumber)
			temp_file, err := ioutil.TempFile("", temp_name)
			defer temp_file.Close()
			if err != nil {
				log.Fatalf("error creating temp file %s", temp_name)
			}

			// Compute reduce
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
			temp_name = temp_file.Name()
			file_name := filepath.Base(temp_name)
			file_name = strings.Split(file_name, ".temp")[0]
			os.Rename(temp_name, file_name)

			// Complete Task
			call_complete_task(task)

		} else {
			// We recieved a "w" task, or wait. This happens when other works are doing map or reduce tasks and we need to wait for further instruction.
			// Wait 5 secs
			time.Sleep(5 * time.Second)
		}

		// Get next task
		task = call_get_task()
	}

	// Done all
}

func call_get_task() GetTaskReply {
	args := GetTaskArgs{}
	reply := GetTaskReply{}
	call_result := call("Coordinator.GetTask", &args, &reply)

	// If call failed, server dead and we can assume all MR tasks arae done.
	if !call_result {
		os.Exit(1)
	}
	return reply
}

func call_complete_task(task GetTaskReply) error {
	args := CompleteTaskArgs{}
	args.TaskType = task.TaskType
	args.TaskNumber = task.TaskNumber
	reply := CompleteTaskReply{}
	call_result := call("Coordinator.CompleteTask", &args, &reply)

	// If call failed, server dead and we can assume all MR tasks arae done.
	if !call_result {
		os.Exit(1)
	}
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
