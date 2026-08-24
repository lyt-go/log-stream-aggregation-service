package decodeplugin

import "errors"

var ErrCorruptFrame = errors.New("decoder rejected corrupt frame")

type Plugin struct{}

func (Plugin) Decode(payload string) ([]string, error) {
	if payload == "corrupt" {
		panic(ErrCorruptFrame)
	}
	return []string{payload}, nil
}
