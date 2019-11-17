package packets

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytesToString(t *testing.T) {
	b := []byte{'a', 'b', 'c'}
	require.Equal(t, "abc", bytesToString(b))
}

func TestDecodeBytes(t *testing.T) {
	expect := []struct {
		rawBytes   []byte
		result     []uint8
		next       int
		offset     int
		shouldFail bool
	}{
		{
			rawBytes: []byte{0, 4, 77, 81, 84, 84, 4, 194, 0, 50, 0, 36, 49, 53, 52}, // ... truncated connect packet (clean session)
			result:   []uint8([]byte{0x4d, 0x51, 0x54, 0x54}),
			next:     6,
			offset:   0,
		},
		{
			rawBytes: []byte{0, 4, 77, 81, 84, 84, 4, 192, 0, 50, 0, 36, 49, 53, 52, 50}, // ... truncated connect packet, only checking start
			result:   []uint8([]byte{0x4d, 0x51, 0x54, 0x54}),
			next:     6,
			offset:   0,
		},
		{
			rawBytes:   []byte{0, 4, 77, 81},
			result:     []uint8([]byte{0x4d, 0x51, 0x54, 0x54}),
			offset:     0,
			shouldFail: true,
		},
		{
			rawBytes:   []byte{0, 4, 77, 81},
			result:     []uint8([]byte{0x4d, 0x51, 0x54, 0x54}),
			offset:     8,
			shouldFail: true,
		},
		{
			rawBytes:   []byte{0, 4, 77, 81},
			result:     []uint8([]byte{0x4d, 0x51, 0x54, 0x54}),
			offset:     0,
			shouldFail: true,
		},
	}

	for i, wanted := range expect {
		result, _, err := decodeBytes(wanted.rawBytes, wanted.offset)
		if wanted.shouldFail {
			require.Error(t, err, "Expected error decoding bytes [i:%d]", i)
			continue
		}

		require.NoError(t, err, "Error decoding bytes [i:%d]", i)
		require.Equal(t, wanted.result, result, "Incorrect decoded value [i:%d]", i)
	}
}

func BenchmarkDecodeBytes(b *testing.B) {
	in := []byte{0, 4, 77, 81, 84, 84, 4, 194, 0, 50, 0, 36, 49, 53, 52}
	for n := 0; n < b.N; n++ {
		decodeBytes(in, 0)
	}
}

func TestDecodeByte(t *testing.T) {
	expect := []struct {
		rawBytes   []byte
		result     uint8
		offset     int
		shouldFail bool
	}{
		{
			rawBytes: []byte{0, 4, 77, 81, 84, 84}, // nonsense slice of bytes
			result:   uint8(0x00),
			offset:   0,
		},
		{
			rawBytes: []byte{0, 4, 77, 81, 84, 84},
			result:   uint8(0x04),
			offset:   1,
		},
		{
			rawBytes: []byte{0, 4, 77, 81, 84, 84},
			result:   uint8(0x4d),
			offset:   2,
		},
		{
			rawBytes: []byte{0, 4, 77, 81, 84, 84},
			result:   uint8(0x51),
			offset:   3,
		},
		{
			rawBytes:   []byte{0, 4, 77, 80, 82, 84},
			result:     uint8(0x00),
			offset:     8,
			shouldFail: true,
		},
	}

	for i, wanted := range expect {
		result, offset, err := decodeByte(wanted.rawBytes, wanted.offset)
		if wanted.shouldFail {
			require.Error(t, err, "Expected error decoding byte [i:%d]", i)
			continue
		}

		require.NoError(t, err, "Error decoding byte [i:%d]", i)
		require.Equal(t, wanted.result, result, "Incorrect decoded value [i:%d]", i)
		require.Equal(t, i+1, offset, "Incorrect offset value [i:%d]", i)
	}
}

func BenchmarkDecodeByte(b *testing.B) {
	in := []byte{0, 4, 77, 81, 84, 84}
	for n := 0; n < b.N; n++ {
		decodeByte(in, 0)
	}
}

func TestDecodeUint16(t *testing.T) {
	expect := []struct {
		rawBytes   []byte
		result     uint16
		offset     int
		shouldFail bool
	}{
		{
			rawBytes: []byte{0, 7, 97, 47, 98, 47, 99, 47, 100, 97},
			result:   uint16(0x07),
			offset:   0,
		},
		{
			rawBytes: []byte{0, 7, 97, 47, 98, 47, 99, 47, 100, 97},
			result:   uint16(0x761),
			offset:   1,
		},
		{
			rawBytes:   []byte{0, 7, 255, 47},
			result:     uint16(0x761),
			offset:     8,
			shouldFail: true,
		},
	}

	for i, wanted := range expect {
		result, offset, err := decodeUint16(wanted.rawBytes, wanted.offset)
		if wanted.shouldFail {
			require.Error(t, err, "Expected error decoding uint16 [i:%d]", i)
			continue
		}

		require.NoError(t, err, "Error decoding uint16 [i:%d]", i)
		require.Equal(t, wanted.result, result, "Incorrect decoded value [i:%d]", i)
		require.Equal(t, i+2, offset, "Incorrect offset value [i:%d]", i)
	}
}

func BenchmarkDecodeUint16(b *testing.B) {
	in := []byte{0, 7, 97, 47, 98, 47, 99, 47, 100, 97}
	for n := 0; n < b.N; n++ {
		decodeUint16(in, 0)
	}
}

func TestDecodeByteBool(t *testing.T) {
	expect := []struct {
		rawBytes   []byte
		result     bool
		offset     int
		shouldFail bool
	}{
		{
			rawBytes: []byte{0x00, 0x00},
			result:   false,
		},
		{
			rawBytes: []byte{0x01, 0x00},
			result:   true,
		},
		{
			rawBytes:   []byte{0x01, 0x00},
			offset:     5,
			shouldFail: true,
		},
	}

	for i, wanted := range expect {
		result, offset, err := decodeByteBool(wanted.rawBytes, wanted.offset)
		if wanted.shouldFail {
			require.Error(t, err, "Expected error decoding byte bool [i:%d]", i)
			continue
		}

		require.NoError(t, err, "Error decoding byte bool [i:%d]", i)
		require.Equal(t, wanted.result, result, "Incorrect decoded value [i:%d]", i)
		require.Equal(t, 1, offset, "Incorrect offset value [i:%d]", i)
	}
}

func BenchmarkDecodeByteBool(b *testing.B) {
	in := []byte{0x00, 0x00}
	for n := 0; n < b.N; n++ {
		decodeByteBool(in, 0)
	}
}

func TestEncodeBool(t *testing.T) {
	result := encodeBool(true)
	require.Equal(t, byte(1), result, "Incorrect encoded value; not true")

	result = encodeBool(false)
	require.Equal(t, byte(0), result, "Incorrect encoded value; not false")

	// Check failure.
	result = encodeBool(false)
	require.NotEqual(t, byte(1), result, "Expected failure, incorrect encoded value")
}

func BenchmarkEncodeBool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		encodeBool(true)
	}
}

func TestEncodeBytes(t *testing.T) {
	result := encodeBytes([]byte("testing"))
	require.Equal(t, []uint8{0, 7, 116, 101, 115, 116, 105, 110, 103}, result, "Incorrect encoded value")

	result = encodeBytes([]byte("testing"))
	require.NotEqual(t, []uint8{0, 7, 113, 101, 115, 116, 105, 110, 103}, result, "Expected failure, incorrect encoded value")
}

func BenchmarkEncodeBytes(b *testing.B) {
	bb := []byte("testing")
	for n := 0; n < b.N; n++ {
		encodeBytes(bb)
	}
}

func TestEncodeUint16(t *testing.T) {
	result := encodeUint16(0)
	require.Equal(t, []byte{0x00, 0x00}, result, "Incorrect encoded value, 0")

	result = encodeUint16(32767)
	require.Equal(t, []byte{0x7f, 0xff}, result, "Incorrect encoded value, 32767")

	result = encodeUint16(65535)
	require.Equal(t, []byte{0xff, 0xff}, result, "Incorrect encoded value, 65535")
}

func BenchmarkEncodeUint16(b *testing.B) {
	for n := 0; n < b.N; n++ {
		encodeUint16(32767)
	}
}

func TestEncodeString(t *testing.T) {
	result := encodeString("testing")
	require.Equal(t, []uint8{0x00, 0x07, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67}, result, "Incorrect encoded value, testing")

	result = encodeString("")
	require.Equal(t, []uint8{0x00, 0x00}, result, "Incorrect encoded value, null")

	result = encodeString("a")
	require.Equal(t, []uint8{0x00, 0x01, 0x61}, result, "Incorrect encoded value, a")

	result = encodeString("b")
	require.NotEqual(t, []uint8{0x00, 0x00}, result, "Expected failure, incorrect encoded value, b")

}

func BenchmarkEncodeString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		encodeString("benchmarking")
	}
}
