package parse

import (
	"testing"
)

func Test_unfold(t *testing.T) {

	type tester struct {
		barr   []byte
		answer []byte
		name   string
	}

	tests := []tester{
		{
			barr:   []byte{97, 97, 13, 10, 32, 9, 13, 10, 9},
			answer: []byte{97, 97, 32},
			name:   "standard byte",
		},
		{
			barr:   []byte{97, 97, 13, 10, 32, 97, 97, 10, 13, 10, 32, 97, 13, 97, 97},
			answer: []byte{97, 97, 32, 97, 97, 10, 32, 97, 13, 97, 97},
			name:   "complex byte",
		},
		{
			barr:   []byte{9, 13, 10, 97, 97, 10, 13, 32, 13, 10, 97, 13, 97, 97, 13, 10},
			answer: []byte{32, 13, 10, 97, 97, 10, 13, 32, 13, 10, 97, 13, 97, 97, 13, 10},
			name:   "same length byte",
		},
		{
			barr:   []byte{9, 13, 10, 32, 13, 10, 9},
			answer: []byte{32},
			name:   "single byte",
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			newBarr := unfold(tst.barr)

			if len(newBarr) != len(tst.answer) {
				t.Error("result did not have proper length")
				return
			}

			for index, val := range newBarr {
				if val != tst.answer[index] {
					t.Error("received invalid result")
					break
				}
			}
		})
	}
}
