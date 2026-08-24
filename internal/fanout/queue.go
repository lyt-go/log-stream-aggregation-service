package fanout

type Queue struct {
	frames [][]byte
}

func (q *Queue) Push(frame []byte) {
	q.frames = append(q.frames, frame)
}

func (q *Queue) Peek() [][]byte {
	return q.frames
}

func (q *Queue) Drain() [][]byte {
	frames := q.frames
	q.frames = nil
	return frames
}
