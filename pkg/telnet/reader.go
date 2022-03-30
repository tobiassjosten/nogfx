package telnet

// http://mud-dev.wikidot.com/telnet:negotiation
// https://blog.ikeran.org/?p=129
// http://pcmicro.com/NetFoss/telnet.html
// https://www.ironrealms.com/gmcp-doc
// https://tintin.sourceforge.io/protocols/mssp/
// MCCP?
// http://www.mushclient.com/mushclient/mxp.htm
// https://wiki.mudlet.org/w/Manual:Supported_Protocols

import (
	"bytes"
	"fmt"
	"io"
)

func (stream *Stream) Read(buffer []byte) (count int, err error) {
	count, err = stream.readConnection(buffer)
	if err != nil && (err != io.EOF || count == 0) {
		return 0, err
	}

	for count = 0; count < len(buffer) && len(stream.buffer) > 0; {
		rawByte := stream.buffer[0]
		stream.buffer = stream.buffer[1:]

		if rawByte == IAC || len(stream.command) > 0 {
			stream.command = append(stream.command, rawByte)
		}

		if len(stream.command) == 0 {
			buffer[count] = rawByte
			count += 1
			continue
		}

		// Two consequetive IACs equals one escaped IAC.
		if bytes.Equal(stream.command, []byte{IAC, IAC}) {
			buffer[count] = IAC
			count += 1
			stream.command = []byte{}
			continue
		}

		// We should negotiate away GA but some servers use it to mark prompts
		// and so we treat it as a simple linebreak.
		// Or? Maybe we can use that instead of linebreaks to detect the end of a message.
		if bytes.Equal(stream.command, []byte{IAC, GA}) {
			stream.buffer = append(stream.linebreak, stream.buffer...)
			stream.command = []byte{}
			continue
		}

		if isCompleteCommand(stream.command) {
			if !isValidCommand(stream.command) {
				stream.buffer = append([]byte(fmt.Sprintf(
					"Received invalid command sequence '%x'\n",
					stream.command,
				)), stream.buffer...)
				stream.command = []byte{}
			} else {
				stream.processCommand()
			}
		}
	}

	return count, nil
}

func (stream *Stream) readConnection(buffer []byte) (count int, err error) {
	if len(stream.buffer) >= len(buffer) {
		return 0, nil
	}

	rawBuffer := make([]byte, len(buffer)-len(stream.buffer))
	count, err = stream.connection.Read(rawBuffer)
	if err != nil && err != io.EOF {
		return count, err
	}

	stream.buffer = append(stream.buffer, rawBuffer[:count]...)

	return len(stream.buffer), err
}

func (stream *Stream) Write(buffer []byte) (int, error) {
	return stream.connection.Write(buffer)
}

func (stream *Stream) Close() error {
	return stream.connection.Close()
}
