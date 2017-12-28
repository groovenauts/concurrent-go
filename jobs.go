package concurrent

import (
	"fmt"
	"strings"
)

type Jobs []*Job

func (jobs Jobs) Error() error {
	messages := []string{}
	for _, t := range jobs {
		if t.Error != nil {
			messages = append(messages, t.Error.Error())
		}
	}
	if len(messages) == 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(messages, "\n"))
}
