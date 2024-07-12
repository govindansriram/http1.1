package parse

import (
	"bytes"
	"errors"
)

const version = "HTTP/1.1"

/*
RFC 2616

The Request-Line begins with a method token, followed by the
Request-URI and the protocol version, and ending with CRLF. The
elements are separated by SP characters. No CR or LF is allowed
except in the final CRLF sequence.

     Request-Line   = Method SP Request-URI SP HTTP-Version CRLF
*/

type RequestLine struct {
	method []byte
	uri    []byte
}

func InitRequestLine(reqLine []byte) (*RequestLine, error) {

	if !bytes.HasSuffix(reqLine, []byte{13, 10}) {
		return nil, errors.New("requestline must end with a CRLF")
	}

	reqLine = bytes.TrimSuffix(reqLine, []byte{13, 10})

	if bytes.ContainsAny(reqLine, "\r\n") {
		return nil, errors.New("requestline contains unnessecary CR and LF tokens")
	}

	lineValues := bytes.Split(reqLine, []byte{32})

	if len(lineValues) != 3 {
		return nil, errors.New("received invalid request line parameter count")
	}

	if !bytes.Equal(lineValues[2], []byte(version)) {
		return nil, errors.New("only requests of type HTTP/1.1 are allowed")
	}

	return &RequestLine{
		method: lineValues[0],
		uri:    lineValues[1],
	}, nil
}

func (r RequestLine) GetMethod() string {
	return string(r.method)
}

func (r RequestLine) GetUri() string {
	return string(r.uri)
}
