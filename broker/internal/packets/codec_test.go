package packets

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytesToString(t *testing.T) {
	b := []byte{'a', 'b', 'c'}
	require.Equal(t, "abc", bytesToString(b))
}

func BenchmarkBytesToString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		bytesToString([]byte{'a', 'b', 'c'})
	}
}

func TestDecodeString(t *testing.T) {
	expect := []struct {
		name       string
		rawBytes   []byte
		result     string
		offset     int
		shouldFail error
	}{
		{
			offset:   0,
			rawBytes: []byte{0, 7, 97, 47, 98, 47, 99, 47, 100, 97},
			result:   "a/b/c/d",
		},
		{
			offset: 14,
			rawBytes: []byte{
				byte(Connect << 4), 17, // Fixed header
				0, 6, // Protocol Name - MSB+LSB
				'M', 'Q', 'I', 's', 'd', 'p', // Protocol Name
				3,     // Protocol Version
				0,     // Packet Flags
				0, 30, // Keepalive
				0, 3, // Client ID - MSB+LSB
				'h', 'e', 'y', // Client ID "zen"},
			},
			result: "hey",
		},
		{
			offset:   2,
			rawBytes: []byte{0, 0, 0, 23, 49, 47, 50, 47, 51, 47, 52, 47, 97, 47, 98, 47, 99, 47, 100, 47, 101, 47, 94, 47, 64, 47, 33, 97},
			result:   "1/2/3/4/a/b/c/d/e/^/@/!",
		},
		{
			offset:   0,
			rawBytes: []byte{0, 5, 120, 47, 121, 47, 122, 33, 64, 35, 36, 37, 94, 38},
			result:   "x/y/z",
		},
		{
			offset:     0,
			rawBytes:   []byte{0, 9, 'a', '/', 'b', '/', 'c', '/', 'd', 'z'},
			shouldFail: ErrOffsetBytesOutOfRange,
		},
		{
			offset:     5,
			rawBytes:   []byte{0, 7, 97, 47, 98, 47, 'x'},
			shouldFail: ErrOffsetBytesOutOfRange,
		},
		{
			offset:     9,
			rawBytes:   []byte{0, 7, 97, 47, 98, 47, 'y'},
			shouldFail: ErrOffsetUintOutOfRange,
		},
		{
			offset: 17,
			rawBytes: []byte{
				byte(Connect << 4), 0, // Fixed header
				0, 4, // Protocol Name - MSB+LSB
				'M', 'Q', 'T', 'T', // Protocol Name
				4,     // Protocol Version
				0,     // Flags
				0, 20, // Keepalive
				0, 3, // Client ID - MSB+LSB
				'z', 'e', 'n', // Client ID "zen"
				0, 6, // Will Topic - MSB+LSB
				'l',
			},
			shouldFail: ErrOffsetBytesOutOfRange,
		},
		{
			offset:     0,
			rawBytes:   []byte{0, 7, 0xc3, 0x28, 98, 47, 99, 47, 100},
			shouldFail: ErrOffsetStrInvalidUTF8,
		},
	}

	for i, wanted := range expect {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, _, err := decodeString(wanted.rawBytes, wanted.offset)
			if wanted.shouldFail != nil {
				require.True(t, errors.Is(err, wanted.shouldFail), "want %v to be a %v", err, wanted.shouldFail)
				return
			}

			require.NoError(t, err)
			require.Equal(t, wanted.result, result)
		})
	}
}

func BenchmarkDecodeString(b *testing.B) {
	in := []byte{0, 7, 97, 47, 98, 47, 99, 47, 100, 97}
	for n := 0; n < b.N; n++ {
		decodeString(in, 0)
	}
}

func TestDecodeBytes(t *testing.T) {
	expect := []struct {
		rawBytes   []byte
		result     []uint8
		next       int
		offset     int
		shouldFail error
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
			offset:     0,
			shouldFail: ErrOffsetBytesOutOfRange,
		},
		{
			rawBytes:   []byte{0, 4, 77, 81},
			offset:     8,
			shouldFail: ErrOffsetUintOutOfRange,
		},
	}

	for i, wanted := range expect {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, _, err := decodeBytes(wanted.rawBytes, wanted.offset)
			if wanted.shouldFail != nil {
				require.True(t, errors.Is(err, wanted.shouldFail), "want %v to be a %v", err, wanted.shouldFail)
				return
			}

			require.NoError(t, err)
			require.Equal(t, wanted.result, result)
		})
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
		shouldFail error
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
			offset:     8,
			shouldFail: ErrOffsetByteOutOfRange,
		},
	}

	for i, wanted := range expect {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, offset, err := decodeByte(wanted.rawBytes, wanted.offset)
			if wanted.shouldFail != nil {
				require.True(t, errors.Is(err, wanted.shouldFail), "want %v to be a %v", err, wanted.shouldFail)
				return
			}

			require.NoError(t, err)
			require.Equal(t, wanted.result, result)
			require.Equal(t, i+1, offset)
		})
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
		shouldFail error
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
			offset:     8,
			shouldFail: ErrOffsetUintOutOfRange,
		},
	}

	for i, wanted := range expect {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, offset, err := decodeUint16(wanted.rawBytes, wanted.offset)
			if wanted.shouldFail != nil {
				require.True(t, errors.Is(err, wanted.shouldFail), "want %v to be a %v", err, wanted.shouldFail)
				return
			}

			require.NoError(t, err)
			require.Equal(t, wanted.result, result)
			require.Equal(t, i+2, offset)
		})
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
		shouldFail error
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
			shouldFail: ErrOffsetBoolOutOfRange,
		},
	}

	for i, wanted := range expect {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, offset, err := decodeByteBool(wanted.rawBytes, wanted.offset)
			if wanted.shouldFail != nil {
				require.True(t, errors.Is(err, wanted.shouldFail), "want %v to be a %v", err, wanted.shouldFail)
				return
			}

			require.NoError(t, err)
			require.Equal(t, wanted.result, result)
			require.Equal(t, 1, offset)
		})
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
