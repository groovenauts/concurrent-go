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

func TestWorkersProcess(t *testing.T) {
	workers := NewWorkers(NewCopyFunc(t), 3)

	// Empty jobs
	jobs0 := Jobs{}
	err := workers.process(jobs0)
	assert.NoError(t, err)

	// 1 job
	jobs1 := Jobs{
		&Job{Payload: &TestPayload{Input: "foo"}},
	}
	err = workers.process(jobs1)
	assert.NoError(t, err)

	// 1 error job
	jobs1Error := Jobs{
		&Job{Payload: &TestPayload{}},
	}
	err = workers.process(jobs1Error)
	assert.Equal(t, "Input is blank", err.Error())

	// 1 success and 1 error
	jobs1Success1Error := Jobs{
		&Job{Payload: &TestPayload{Input: "foo"}},
		&Job{Payload: &TestPayload{}},
	}
	err = workers.process(jobs1Success1Error)
	assert.Equal(t, "Input is blank", err.Error())

	// 3 success and 2 error
	jobs3Success2Error := Jobs{
		&Job{Payload: &TestPayload{}},
		&Job{Payload: &TestPayload{Input: "foo"}},
		&Job{Payload: &TestPayload{Input: "foo"}},
		&Job{Payload: &TestPayload{}},
		&Job{Payload: &TestPayload{Input: "foo"}},
	}
	err = workers.process(jobs3Success2Error)
	assert.Equal(t, "Input is blank\nInput is blank", err.Error())

	// 3 success
	jobs3Success := Jobs{
		&Job{Payload: &TestPayload{Input: "foo"}},
		&Job{Payload: &TestPayload{Input: "foo"}},
		&Job{Payload: &TestPayload{Input: "foo"}},
	}
	err = workers.process(jobs3Success)
	assert.NoError(t, err)
}
