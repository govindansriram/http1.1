package urlparser

import (
	"bytes"
	"encoding/hex"
)

type Host struct {
	Type  string
	Value string
}

type Authority struct {
	UserInformation *string
	Port            uint16
	Host            Host
}

type Uri struct {
	Scheme    string
	Authority Authority
	Path      string
	Query     string
	Fragment  string
}

/*
getUnreservedHexs

# RFC 3986 2.3

	Octets in the range %41-5A, %30-%39 and octets %2D, %2E, %5F, %7E are known
	as unreserved characters.

	this function returns a slice consisting of all unreserved characters in 3 byte
	segements
*/
func getUnreservedHexs() []byte {
	validBytes := make([]byte, 0, 120)
	data := []byte{
		'a', 'b', 'c', 'd', 'e', 'f',
	}

	for i := 1; i < 10; i++ {
		validBytes = append(validBytes, '%', '4', byte('0'+i))
	}

	for i := uint8(0); i < 10; i++ {
		validBytes = append(validBytes, '%', '3', byte('0'+i))
	}

	for _, i := range data {
		validBytes = append(validBytes, '%', '4', i)
	}

	for i := uint8(0); i < 10; i++ {
		validBytes = append(validBytes, '%', '5', byte('0'+i))
	}

	validBytes = append(validBytes, []byte{
		'%', '2', 'd', '%', '2', 'e', '%', '5', 'f', '%', '7', 'e', '%', '5', 'a',
	}...)

	return validBytes
}

/*
normalizeUri

# RFC 3986 2.3

	Octets in the range %41-5A, %30-%39 and octets %2D, %2E, %5F, %7E are known
	as unreserved characters.

	These characters if in hexadecimal form in a uri should
	be automatically converted to their their ascii equivalent.

	This function does just that

	EX: `%41a%4Dbc` -> `AaMbc`
*/
func normalizeUri(uri []byte) ([]byte, error) {
	unreserved := getUnreservedHexs()
	normalizedUri := make([]byte, 0, len(uri))

	seq := [2]byte{}
	index := -1

	for _, data := range uri {

		if data == '%' {
			if index == 0 {
				normalizedUri = append(normalizedUri, '%')
			}
			index = 0
			continue
		}

		if index > -1 {
			seq[index] = data
			index++
		} else {
			normalizedUri = append(normalizedUri, data)
		}

		if index == 2 {
			index = -1
			if bytes.Contains(
				unreserved,
				bytes.ToLower([]byte{'%', seq[0], seq[1]})) {
				data, err := hex.DecodeString(string(seq[:]))

				if err != nil {
					return nil, err
				}

				normalizedUri = append(normalizedUri, data...)
			} else {
				normalizedUri = append(normalizedUri, '%', seq[0], seq[1])
			}
		}
	}

	if index > -1 {
		normalizedUri = append(normalizedUri, '%')
	}

	if index > 0 {
		normalizedUri = append(normalizedUri, seq[0])
	}

	retUri := make([]byte, len(normalizedUri))
	copy(retUri, normalizedUri)

	return retUri, nil
}
