package concurrent

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestPayload struct {
	Input string
}

func NewCopyFunc(t *testing.T) func(job *Job) error {
	return func(job *Job) error {
		payload, ok := job.Payload.(*TestPayload)
		if assert.True(t, ok) {
			if payload.Input == "" {
				return fmt.Errorf("Input is blank")
			}
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	}
}

func TestWorkersProcessEmptyJobs(t *testing.T) {
	workers := NewWorkers(NewCopyFunc(t), 3)
	jobs := Jobs{}
	err := workers.process(jobs)
	assert.NoError(t, err)
}

func TestWorkersProcess1Job(t *testing.T) {
	workers := NewWorkers(NewCopyFunc(t), 3)
	jobs := Jobs{
		&Job{Payload: &TestPayload{Input: "foo"}},
	}
	err := workers.process(jobs)
	assert.NoError(t, err)
}

func TestWorkersProcess1ErrorJob(t *testing.T) {
	workers := NewWorkers(NewCopyFunc(t), 3)
	jobs := Jobs{
		&Job{Payload: &TestPayload{}},
	}
	err := workers.process(jobs)
	assert.Equal(t, "Input is blank", err.Error())
}

func TestWorkersProcess1SuccessAnd1Error(t *testing.T) {
	workers := NewWorkers(NewCopyFunc(t), 3)
	jobs := Jobs{
		&Job{Payload: &TestPayload{Input: "foo"}},
		&Job{Payload: &TestPayload{}},
	}
	err := workers.process(jobs)
	assert.Equal(t, "Input is blank", err.Error())
}

func TestWorkersProcess3SuccessesAnd2Errors(t *testing.T) {
	workers := NewWorkers(NewCopyFunc(t), 3)
	jobs := Jobs{
		&Job{Payload: &TestPayload{}},
		&Job{Payload: &TestPayload{Input: "foo"}},
		&Job{Payload: &TestPayload{Input: "foo"}},
		&Job{Payload: &TestPayload{}},
		&Job{Payload: &TestPayload{Input: "foo"}},
	}
	err := workers.process(jobs)
	assert.Equal(t, "Input is blank\nInput is blank", err.Error())
}

func TestWorkersProcess3Successes(t *testing.T) {
	workers := NewWorkers(NewCopyFunc(t), 3)
	jobs := Jobs{
		&Job{Payload: &TestPayload{Input: "foo"}},
		&Job{Payload: &TestPayload{Input: "foo"}},
		&Job{Payload: &TestPayload{Input: "foo"}},
	}
	err := workers.process(jobs)
	assert.NoError(t, err)
}
