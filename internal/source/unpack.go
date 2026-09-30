// Package source decodes the original intro's artwork and animation data.
package source

import (
	"encoding/binary"
	"fmt"
)

// UnpackPRG follows the executable's backward word-bitstream decoder. The
// result is a complete, unrelocated TOS executable, not an emulator snapshot.
func UnpackPRG(data []byte) ([]byte, error) {
	if len(data) < 40 || binary.BigEndian.Uint16(data) != 0x601a {
		return nil, fmt.Errorf("source: invalid TOS executable header")
	}
	end64 := uint64(28) + uint64(binary.BigEndian.Uint32(data[2:]))
	if end64 > uint64(len(data)) || end64 < 40 {
		return nil, fmt.Errorf("source: truncated packed executable")
	}
	end := int(end64)
	size := int(binary.BigEndian.Uint32(data[end-4:]))
	if size < 28 || size > 16<<20 {
		return nil, fmt.Errorf("source: invalid unpacked size %d", size)
	}
	input := end - 12
	word := binary.BigEndian.Uint32(data[input:])
	checksum := binary.BigEndian.Uint32(data[end-8:]) ^ word
	var readErr error
	bit := func() uint32 {
		carry := word & 1
		word >>= 1
		if word == 0 {
			input -= 4
			if input < 0x2a0 {
				readErr = fmt.Errorf("source: packed bitstream underflow")
				return 0
			}
			word = binary.BigEndian.Uint32(data[input:])
			checksum ^= word
			carry = word & 1
			word = word>>1 | 0x80000000
		}
		return carry
	}
	bits := func(count int) int {
		value := uint32(0)
		for range count {
			value = value<<1 | bit()
		}
		return int(value)
	}
	out := make([]byte, size)
	position := size
	for position > 0 {
		literal, count, distanceBits := false, 0, 0
		if bit() == 0 {
			if bit() == 0 {
				literal, count = true, bits(3)+1
			} else {
				count, distanceBits = 2, 8
			}
		} else {
			code := bits(2)
			switch code {
			case 0, 1:
				count, distanceBits = code+3, code+9
			case 2:
				count, distanceBits = bits(8)+1, 12
			case 3:
				literal, count = true, bits(8)+9
			}
		}
		if count > position {
			return nil, fmt.Errorf("source: unpacked run exceeds destination")
		}
		if literal {
			for range count {
				position--
				out[position] = byte(bits(8))
			}
		} else {
			distance := bits(distanceBits)
			for range count {
				position--
				if distance < 1 || position+distance >= size {
					return nil, fmt.Errorf("source: invalid backward match at %d, distance %d", position, distance)
				}
				out[position] = out[position+distance]
			}
		}
		if readErr != nil {
			return nil, readErr
		}
	}
	if checksum != 0 {
		return nil, fmt.Errorf("source: packed checksum mismatch: %08x", checksum)
	}
	if binary.BigEndian.Uint16(out) != 0x601a {
		return nil, fmt.Errorf("source: unpacked data has no TOS header")
	}
	return out, nil
}
