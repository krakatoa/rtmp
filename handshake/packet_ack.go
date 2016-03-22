package handshake

import (
	"bytes"
	"encoding/binary"
	"io"
)

// AckPacket represents the RTMP packet used in communicating and agreeing upon
// data during the handshake process. It encapsulates three fields, two of which
// are _usually_ used for sending timestamps (epochs), and the third which is
// used to exchange challenge sequences.
type AckPacket struct {
	// Time1 is the first four-byte field in the AckPacket. See RTMP
	// specification for more details.
	Time1 uint32
	// Time2 is second first four-byte field in the AckPacket. See RTMP
	// specification for more details.
	Time2 uint32
	// Payload is the challenge payload that is sent back and forth between
	// the client and the server to agree that they are speaking securely.
	Payload [1528]byte
}

// Read reads the encoded data representing an AckPacket from the given
// io.Reader and serializes it into the struct.
func (a *AckPacket) Read(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &a.Time1); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &a.Time2); err != nil {
		return err
	}

	if _, err := r.Read(a.Payload[:]); err != nil {
		return err
	}

	return nil
}

// Write writes the data contained in the struct out over the given io.Writer.
//
// If an error is encountered during any of these writes, the error is returned
// immediately.
func (a *AckPacket) Write(w io.Writer) error {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, a.Time1); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.BigEndian, a.Time2); err != nil {
		return err
	}

	buf.Write(a.Payload[:])

	if _, err := buf.WriteTo(w); err != nil {
		return err
	}

	return nil
}