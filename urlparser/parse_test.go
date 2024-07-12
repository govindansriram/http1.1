package urlparser

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_getUnreservedHexs(t *testing.T) {
	data := getUnreservedHexs()
	if len(data) != cap(data) {
		t.Error("capacity must equal length")
	}
}

func Test_normalizeUri(t *testing.T) {

	type test struct {
		uri      []byte
		expected []byte
		isErr    bool
		name     string
	}

	tests := []test{
		{
			uri: []byte{
				'%', '4', '1', 'a', '%', '4', 'D', 'b', 'c',
			},
			expected: []byte{
				'A', 'a', 'M', 'b', 'c',
			},
			isErr: false,
			name:  "standard",
		},
		{
			uri: []byte{
				'%', '4', '1', 'a', '%', '4', 'D', 'b', 'c', '%', 'E',
			},
			expected: []byte{
				'A', 'a', 'M', 'b', 'c', '%', 'E',
			},
			isErr: false,
			name:  "two encoded letters",
		},
		{
			uri: []byte{
				'%', '4', '1', 'a', '%', '4', 'D', 'b', 'c', '%', 'E', 'F',
			},
			expected: []byte{
				'A', 'a', 'M', 'b', 'c', '%', 'E', 'F',
			},
			isErr: false,
			name:  "reserved encoded data",
		},
		{
			uri: []byte{
				'%', '4', '1', 'a', '%', '%', '4', 'D', 'b', 'c',
			},
			expected: []byte{
				'A', 'a', '%', 'M', 'b', 'c',
			},
			isErr: false,
			name:  "extra percent",
		},
		{
			uri: []byte{
				'%', '4', '1', 'a', '%', '%', '4', 'D', '%', 'E',
			},
			expected: []byte{
				'A', 'a', '%', 'M', '%', 'E',
			},
			isErr: false,
			name:  "extra percent and unfinished chars",
		},
		{
			uri: []byte{
				'%', '%', '4', '1', 'a', '%', '%', '4', 'D', 'b', 'c', '%',
			},
			expected: []byte{
				'%', 'A', 'a', '%', 'M', 'b', 'c', '%',
			},
			isErr: false,
			name:  "percent add the end",
		},
	}

	for _, tst := range tests {

		t.Run(tst.name, func(t *testing.T) {
			res, err := normalizeUri(tst.uri)

			if err == nil && tst.isErr {
				t.Errorf("function should have errored")
			}

			if err != nil && !tst.isErr {
				t.Errorf("got unexpected error %v", err)
			}

			if !bytes.Equal(res, tst.expected) {
				fmt.Println(string(res))
				t.Errorf("invalid result recieved %s", res)
			}

		})
	}
}
