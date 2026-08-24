package batchparser

import "errors"

var ErrMalformed = errors.New("malformed log batch")

type Parser struct{}

func (Parser) Parse(payload string) ([]string, error) {
	if payload == "malformed" {
		panic(ErrMalformed)
	}
	return []string{payload}, nil
}
