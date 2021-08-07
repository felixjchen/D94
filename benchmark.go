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

func get_set_command(key_number int) string {
	set_command := "key" + strconv.Itoa(key_number) + " value" + strconv.Itoa(rand.Intn(records))
	return set_command
}

func get_get_command(key_number int) string {
	return "key" + strconv.Itoa(key_number)
}

func get_commands(prefix string, count int, command_fn func(int) string) [][]string {
	var res [][]string
	for i := 0; i < count; i++ {
		command := strings.Split(prefix+command_fn(i), " ")
		res = append(res, command)
	}
	return res
}

func shuffle_workload(workload [][]string) {
	rand.Shuffle(len(workload), func(i int, j int) {
		workload[i], workload[j] = workload[j], workload[i]
	})
}

func get_workload_A(prefix string) [][]string {
	// 50 Read / 50 Read
	reads := 0.5 * records
	writes := 0.5 * records

	write_workload := get_commands(prefix+" set ", int(writes), get_set_command)
	read_workload := get_commands(prefix+" get ", int(reads), get_get_command)
	workload := append(write_workload, read_workload...)
	shuffle_workload(workload)

	return workload
}

func get_workload_B(prefix string) [][]string {
	// 95 Read / 5 Write
	reads := 0.95 * records
	writes := 0.05 * records

	write_workload := get_commands(prefix, int(writes), get_set_command)
	read_workload := get_commands(prefix, int(reads), get_get_command)
	workload := append(write_workload, read_workload...)
	shuffle_workload(workload)

	return workload
}

func main() {
	workloadA := get_workload_A(ecc_prefix)
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
