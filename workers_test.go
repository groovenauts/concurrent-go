package concurrent

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkersProcess(t *testing.T) {
	type TestPayload struct {
		Data string
	}
	f := func(job *Job) error {
		payload, ok := job.Payload.(*TestPayload)
		if assert.True(t, ok) {
			if payload.Data == "" {
				return fmt.Errorf("Data is blank")
			}
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	workers := Workers{
		&Worker{Proc: f},
		&Worker{Proc: f},
		&Worker{Proc: f},
	}

	// Empty jobs
	jobs0 := Jobs{}
	err := workers.process(jobs0)
	assert.NoError(t, err)

	// 1 job
	jobs1 := Jobs{
		&Job{Payload: &TestPayload{Data: "foo"}},
	}
	err = workers.process(jobs1)
	assert.NoError(t, err)

	// 1 error job
	jobs1Error := Jobs{
		&Job{Payload: &TestPayload{}},
	}
	err = workers.process(jobs1Error)
	assert.Equal(t, "Data is blank", err.Error())

	// 1 success and 1 error
	jobs1Success1Error := Jobs{
		&Job{Payload: &TestPayload{Data: "foo"}},
		&Job{Payload: &TestPayload{}},
	}
	err = workers.process(jobs1Success1Error)
	assert.Equal(t, "Data is blank", err.Error())

	// 3 success and 2 error
	jobs3Success2Error := Jobs{
		&Job{Payload: &TestPayload{}},
		&Job{Payload: &TestPayload{Data: "foo"}},
		&Job{Payload: &TestPayload{Data: "foo"}},
		&Job{Payload: &TestPayload{}},
		&Job{Payload: &TestPayload{Data: "foo"}},
	}
	err = workers.process(jobs3Success2Error)
	assert.Equal(t, "Data is blank\nData is blank", err.Error())

	// 3 success
	jobs3Success := Jobs{
		&Job{Payload: &TestPayload{Data: "foo"}},
		&Job{Payload: &TestPayload{Data: "foo"}},
		&Job{Payload: &TestPayload{Data: "foo"}},
	}
	err = workers.process(jobs3Success)
	assert.NoError(t, err)
}
