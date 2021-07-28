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

// const records = 10000000
const records = 20
const keys = 1000
const values = 1000

func get_random_key_value() string {
	return "key" + strconv.Itoa(rand.Intn(keys)) + " value" + strconv.Itoa(rand.Intn(values))
}

func get_key_values(prefix string, count int) []string {

	var res []string

	for i := 0; i < count; i++ {
		res = append(res, prefix+get_random_key_value())
	}

	return res
}

func shuffle_workload(workload []string) {

	rand.Shuffle(len(workload), func(i, j int) {
		workload[i], workload[j] = workload[j], workload[i]
	})

}

func get_workload_A() []string {
	// 50 Read / 50 Read
	reads := 0.5 * records
	writes := 0.5 * records

	workload := get_key_values(ecc_prefix+" set ", int(writes))
	workload = append(workload, get_key_values(ecc_prefix+" get ", int(reads))...)

	shuffle_workload(workload)

	return workload
}

func get_workload_B() []string {
	// 95 Read / 5 Write
	reads := 0.95 * records
	writes := 0.05 * records

	workload := get_key_values(ecc_prefix+" set ", int(writes))
	workload = append(workload, get_key_values(ecc_prefix+" get ", int(reads))...)

	shuffle_workload(workload)

	return workload
}

func main() {
	workloadA := get_workload_A()

	start := time.Now()

	for _, command := range workloadA {
		command := strings.Split(command, " ")

		cmd := exec.Command(command[0], command[1:]...)
		stdout, _ := cmd.Output()

		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	return
		// }

		// Print the output
		fmt.Println(string(stdout))
	}

	elapsed := time.Since(start)
	fmt.Printf("Benchmark took %s", elapsed)
}
