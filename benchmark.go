package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// https://tikv.org/blog/tikv-3.0ga/

const ecc_prefix = "distributed_cache ecc client"
const raft_prefix = "distributed_cache raft client"

const records = 10000
const keys = records
const values = records

func set_random_key_value() string {
	return "key" + strconv.Itoa(rand.Intn(keys)) + " value" + strconv.Itoa(rand.Intn(values))
}

func get_random_key() string {
	return "key" + strconv.Itoa(rand.Intn(keys))
}

func set_key_values(prefix string, count int) [][]string {
	var res [][]string

	for i := 0; i < count; i++ {
		command := strings.Split(prefix+set_random_key_value(), " ")
		res = append(res, command)
	}
	return res
}

func get_keys(prefix string, count int) [][]string {
	var res [][]string

	for i := 0; i < count; i++ {
		command := strings.Split(prefix+get_random_key(), " ")
		res = append(res, command)
	}
	return res
}

func shuffle_workload(workload [][]string) {

	rand.Shuffle(len(workload), func(i, j int) {
		workload[i], workload[j] = workload[j], workload[i]
	})

}

func get_workload_A(prefix string) [][]string {
	// 50 Read / 50 Read
	reads := 0.5 * records
	writes := 0.5 * records

	workload := set_key_values(prefix+" set ", int(writes))
	workload = append(workload, get_keys(prefix+" get ", int(reads))...)
	shuffle_workload(workload)

	return workload
}

func get_workload_B(prefix string) [][]string {
	// 95 Read / 5 Write
	reads := 0.95 * records
	writes := 0.05 * records

	workload := set_key_values(prefix+" set ", int(writes))
	workload = append(workload, get_keys(prefix+" get ", int(reads))...)
	shuffle_workload(workload)

	return workload
}

func main() {
	workloadA := get_workload_A(raft_prefix)
	fmt.Println("Generated workload")

	start := time.Now()

	for i, command := range workloadA {
		cmd := exec.Command(command[0], command[1:]...)
		stdout, _ := cmd.Output()
		fmt.Println(i, string(stdout))
	}

	elapsed := time.Since(start)
	fmt.Printf("Benchmark took %s", elapsed)
}
