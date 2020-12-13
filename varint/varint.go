package varint

import (
	"encoding/binary"
	"errors"
)

func Read(b []byte) uint64 {
	i := b[0]
	b = b[1:]

	if i == 0xfd {
		return uint64(binary.LittleEndian.Uint16(b[:2]))
	}

	if i == 0xfe {
		return uint64(binary.LittleEndian.Uint32(b[:4]))
	}

	if i == 0xff {
		return uint64(binary.LittleEndian.Uint64(b[:8]))
	}

	return uint64(i)
}

func Encode(i uint64) ([]byte, error) {
	if i < 0xfd {
		return []byte{byte(int(i))}, nil
	}

	if i < 0x10000 {
		d := make([]byte, 2)
		binary.LittleEndian.PutUint16(d, uint16(i))
		return append([]byte{0xfd}, d...), nil
	}

	if i < 0x100000000 {
		d := make([]byte, 4)
		binary.LittleEndian.PutUint32(d, uint32(i))
		return append([]byte{0xfe}, d...), nil
	}

	if i <= 0xffffffffffffffff {
		d := make([]byte, 8)
		binary.LittleEndian.PutUint64(d, uint64(i))
		return append([]byte{0xff}, d...), nil
	}

	return nil, errors.New("int to large")
}
