package format

import (
	"testing"
)

func checkResult(t *testing.T, b []byte, s, r interface{}) {

	switch s.(type) {
	case uint16:
		if s.(uint16) != r.(uint16) {
			t.Fatalf("%v --> expect %d, got %d", b, r, s)
		}

	case uint32:
		if s.(uint32) != r.(uint32) {
			t.Fatalf("%v --> expect %d, got %d", b, r, s)
		}

	default:
		if s != r {
			t.Fatalf("%v --> expect %q, got %q", b, r, s)
		}
	}

}

func TestCanonicalWireFormat1(t *testing.T) {
	b := []byte{3, 'F', 'O', 'O', 3, 'B', 'A', 'R'}
	s := CanonicalWireFormat(b)
	checkResult(t, b, s, "FOO.BAR")
}

func TestCanonicalWireFormat2(t *testing.T) {
	b := []byte{3, 'F', 'O', 'O', 3, 'B', 'A', 'R', 0}
	s := CanonicalWireFormat(b)
	checkResult(t, b, s, "FOO.BAR")
}

func TestCanonicalWireFormat3(t *testing.T) {
	b := []byte{1, 'F', 0}
	s := CanonicalWireFormat(b)
	checkResult(t, b, s, "F")
}

func TestCanonicalWireFormat4(t *testing.T) {
	b := []byte{0}
	s := CanonicalWireFormat(b)
	checkResult(t, b, s, "")
}

func TestCanonicalWireFormatInvalid(t *testing.T) {
	b := []byte{3, 'F', 'O', 'O', 3, 'B', 'A'}
	s := CanonicalWireFormat(b)
	checkResult(t, b, s, "FOO.")
}

func TestUint16B1(t *testing.T) {
	b := []byte{0, 0}
	s := Uint16B(b)
	checkResult(t, b, s, uint16(0))
}

func TestUint16B2(t *testing.T) {
	b := []byte{0xff, 0xfe}
	s := Uint16B(b)
	checkResult(t, b, s, uint16(0xfffe))
}

func TestUint32B1(t *testing.T) {
	b := []byte{0, 0, 0, 0}
	s := Uint32B(b)
	checkResult(t, b, s, uint32(0))
}

func TestUint32B2(t *testing.T) {
	b := []byte{0xff, 0xff, 0xff, 0xfe}
	s := Uint32B(b)
	checkResult(t, b, s, uint32(0xfffffffe))
}

func TestIPv4String(t *testing.T) {
	b := []byte{1, 2, 3, 255}
	s := IPv4String(b)
	checkResult(t, b, s, "1.2.3.255")
}

func TestMACAddressString(t *testing.T) {
	b := []byte{1, 2, 3, 4, 5, 255}
	s := MACAddressString(b)
	checkResult(t, b, s, "01:02:03:04:05:ff")
}
