package main

import "time"

type rawTime []byte

func (t rawTime) unify() (time.Time, error) {
	time, err := time.Parse("2006-01-02 15:04:05", string(t))
	if err != nil {
		return time, err
	}

	return time, nil
}
