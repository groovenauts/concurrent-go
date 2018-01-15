package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	concurrent ".."
)

type SumPayload struct {
	Input  int64
	Output int64
}

func sum(job *concurrent.Job) error {
	payload, ok := job.Payload.(*SumPayload)
	if !ok {
		return fmt.Errorf("Unsupported payload %v\n", job.Payload)
	}
	var i int64
	var r int64
	for i = 1; i <= payload.Input; i++ {
		r = r + i
	}
	time.Sleep(100)
	payload.Output = r
	fmt.Fprintf(os.Stderr, "SUCCESS payload: %v => %v\n", payload.Input, payload.Output)
	return nil
}

func main() {
	fmt.Fprintf(os.Stderr, "runtime.NumCPU() : %v\n", runtime.NumCPU())

	numProc, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid number proc: %v\n", os.Args[1])
		os.Exit(1)
	}
	oldMaxProcs := runtime.GOMAXPROCS(numProc)
	fmt.Fprintf(os.Stderr, "GOMAXPROCS from %v to %v\n", oldMaxProcs, numProc)

	numWorkers, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid number of workers : %v\n", os.Args[2])
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "numWorkers : %v\n", numWorkers)

	numStrs := os.Args[3:]
	jobs := concurrent.Jobs{}
	for _, numStr := range numStrs {
		num, err := strconv.ParseInt(numStr, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid number argument: %v\n", numStr)
			os.Exit(1)
		}
		jobs = append(jobs, &concurrent.Job{Payload: &SumPayload{Input: num}})
	}
	workers := concurrent.NewWorkers(sum, numWorkers)
	workers.Process(jobs)
	err = jobs.Error()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	}
}
