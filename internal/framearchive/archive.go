package framearchive

type Archive struct {
	frames [][]byte
}

func (a *Archive) Append(frame []byte) {
	a.frames = append(a.frames, frame)
}

func (a *Archive) Frames() [][]byte {
	return a.frames
}
