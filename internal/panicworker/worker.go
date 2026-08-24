package panicworker

import "fmt"

type Worker struct{ poisoned bool }

func (w *Worker) Execute(payload string) (result string, err error) {
	if w.poisoned {
		return "", fmt.Errorf("worker remains poisoned")
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			w.poisoned = true
			result = ""
			err = fmt.Errorf("worker panic: %v", recovered)
		}
	}()
	if payload == "panic" {
		panic("decoder exploded")
	}
	return "decoded:" + payload, nil
}
