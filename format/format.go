package format

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

func CanonicalWireFormat(b []byte) string {
	var buf bytes.Buffer
	i := 0
	for {
		length := int(b[i])
		if length == 0 {
			break
		}
		i++

		if i+length > len(b) {
			break
		}
		buf.Write(b[i : i+length])

		i += length
		if i >= len(b) {
			break
		}

		if b[i] != 0 {
			buf.WriteString(".")
		}
	}
	return buf.String()
}

func Uint16B(b []byte) uint16 {
	buf := bytes.NewBuffer(b)
	var x uint16
	binary.Read(buf, binary.BigEndian, &x)
	return x
}

func Uint32B(b []byte) uint32 {
	buf := bytes.NewBuffer(b)
	var x uint32
	binary.Read(buf, binary.BigEndian, &x)
	return x
}

func IPv4String(b []byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}

func MACAddrString(b []byte) string {
	var buf bytes.Buffer
	for i := range b {
		if i > 0 {
			buf.WriteString(":")
		}
		buf.WriteString(fmt.Sprintf("%02x", b[i]))
	}

	return buf.String()
}

func YesNo(b []byte) string {
	// yes or no
	if b[0] == 0 {
		return "no"
	}
	return "yes"
}

func DurationString(b []byte) string {
	t := time.Duration(Uint32B(b)) * time.Second
	s := int(t.Seconds()) % 60
	m := int(t.Minutes()) % 60
	h := int(t.Hours()) % 24
	d := int(t.Hours()) / 24
	return fmt.Sprintf("%dd%dh%dm%ds", d, h, m, s)
}

func String(b []byte) string {
	return fmt.Sprintf("%q", string(b))
}

// Types according to RFC 1700
func RFC1700Types(b []byte) string {
	switch b[0] {
	case 1:
		return MACAddrString(b[1:7])
	default:
		return fmt.Sprintf("type %d (len %d)", b[0], len(b)-1)
	}
}
