package concurrent

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestPayload struct {
	Input  string
	Output string
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
		// Just Copy Input to Output
		payload.Output = payload.Input
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
	payload1 := &TestPayload{Input: "foo"}
	jobs := Jobs{
		&Job{Payload: payload1},
	}
	err := workers.process(jobs)
	assert.NoError(t, err)
	assert.Equal(t, "foo", payload1.Output)
}

func TestWorkersProcess1ErrorJob(t *testing.T) {
	workers := NewWorkers(NewCopyFunc(t), 3)
	payload1 := &TestPayload{}
	jobs := Jobs{
		&Job{Payload: payload1},
	}
	err := workers.process(jobs)
	assert.Equal(t, "Input is blank", err.Error())
	assert.Zero(t, payload1.Output)
}

func TestWorkersProcess1SuccessAnd1Error(t *testing.T) {
	workers := NewWorkers(NewCopyFunc(t), 3)
	payload1 := &TestPayload{Input: "foo"}
	payload2 := &TestPayload{}
	jobs := Jobs{
		&Job{Payload: payload1},
		&Job{Payload: payload2},
	}
	err := workers.process(jobs)
	assert.Equal(t, "Input is blank", err.Error())
	assert.Equal(t, "foo", payload1.Output)
	assert.Zero(t, payload2.Output)
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
	payloads := []*TestPayload{
		&TestPayload{Input: "foo"},
		&TestPayload{Input: "foo"},
		&TestPayload{Input: "foo"},
	}
	jobs := Jobs{}
	for _, payload := range payloads {
		jobs = append(jobs, &Job{Payload: payload})
	}
	err := workers.process(jobs)
	assert.NoError(t, err)
	for _, payload := range payloads {
		assert.Equal(t, "foo", payload.Output)
	}
}
