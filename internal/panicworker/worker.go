package panicworker

import "fmt"

// Worker 是无状态的：每次 Execute 调用相互独立。
// 某个批次的 decoder panic 被恢复后，不得污染后续批次的执行，
// 因此不再保留跨调用的 poisoned 状态。
type Worker struct{}

func (w *Worker) Execute(payload string) (result string, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = ""
			err = fmt.Errorf("worker panic: %v", recovered)
		}
	}()
	if payload == "panic" {
		panic("decoder exploded")
	}
	return "decoded:" + payload, nil
}
