package parse

/*
   HTTP/1.1 header field values can be folded onto multiple lines if the
   continuation line begins with a space or horizontal tab. All linear
   white space, including folding, has the same semantics as SP. A
   recipient MAY replace any linear white space with a single SP before
   interpreting the field value or forwarding the message downstream.

   LWS            = [CRLF] 1*( SP | HT )
   CRLF           = CR LF
   CR             = <US-ASCII CR, carriage return (13)>
   LF             = <US-ASCII LF, linefeed (10)>
   SP             = <US-ASCII SP, space (32)>
   HT             = <US-ASCII HT, horizontal-tab (9)>
*/

/*
unfold

replaces LWS with SP, as permitted by RFC above. Please use for headers,
helps with readability.
*/
func unfold(asciiEncoding []byte) []byte {
	encodingQueue := make([]byte, 0, len(asciiEncoding))
	lwsQueue := make([]byte, 0, 2)

	flush := func() {
		encodingQueue = append(encodingQueue, lwsQueue...)
		lwsQueue = make([]byte, 0, 2)
	}

	checkIndex := func(expected byte) bool {
		if len(lwsQueue) == 0 {
			return false
		}
		return lwsQueue[len(lwsQueue)-1] == expected
	}

	for _, char := range asciiEncoding {
		if char == 9 || char == 32 {
			if checkIndex(10) {
				lwsQueue = make([]byte, 0, 2)
			}

			if checkIndex(13) {
				flush()
			}

			if len(lwsQueue) == 0 {
				lwsQueue = append(lwsQueue, 32)
			}
		} else if char == 13 {
			lwsQueue = append(lwsQueue, char)
		} else if char == 10 {
			if checkIndex(13) {
				lwsQueue = append(lwsQueue, char)
			} else {
				flush()
				encodingQueue = append(encodingQueue, char)
			}
		} else {
			flush()
			encodingQueue = append(encodingQueue, char)
		}
	}

	flush()
	fixedEncoding := make([]byte, len(encodingQueue))
	copy(fixedEncoding, encodingQueue)
	return encodingQueue
}
