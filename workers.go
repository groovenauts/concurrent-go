package concurrent

import (
	"time"
)

type Workers []*Worker

func NewWorkers(proc func(job *Job) error, workers int) Workers {
	r := Workers{}
	for i := 0; i < workers; i++ {
		r = append(r, &Worker{Proc: proc})
	}
	return r
}

func (ws Workers) process(jobs Jobs) error {
	c := make(chan *Job, len(jobs))
	for _, job := range jobs {
		c <- job
	}

	for _, w := range ws {
		w.done = false
		w.jobs = c
		go w.run()
	}

	for {
		time.Sleep(100 * time.Millisecond)
		if ws.done() {
			break
		}
	}

	return jobs.error()
}

func (ws Workers) done() bool {
	for _, w := range ws {
		if !w.done {
			return false
		}
	}
	return true
}
