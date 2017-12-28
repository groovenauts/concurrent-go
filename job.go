package concurrent

type Job struct {
	Payload interface{}
	Error   error
}
