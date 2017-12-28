package concurrent

type Worker struct {
	jobs chan *Job
	Proc func(job *Job) error
	done bool
}

func (w *Worker) run() {
	for {
		var job *Job
		select {
		case job = <-w.jobs:
		default: // Do nothing to break
		}
		if job == nil {
			// No job found any more
			w.done = true
			break
		}
		if job.Error != nil {
			continue
		}

		err := w.Proc(job)
		if err != nil {
			job.Error = err
			continue
		}
	}
}
