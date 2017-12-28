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

func (ws Workers) Process(jobs Jobs) {
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
		if ws.Done() {
			break
		}
	}
}

func (ws Workers) Done() bool {
	for _, w := range ws {
		if !w.done {
			return false
		}
	}
	return true
}
