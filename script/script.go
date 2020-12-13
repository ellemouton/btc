package script

import (
	"encoding/binary"

	"github.com/ellemouton/btc/varint"
)

type Script []elem

type elem struct {
	value byte
	data  []byte
}

func read(b []byte, n int) ([]byte, []byte) {
	res := b[:n]
	b = b[n:]
	return res, b
}

func Parse(b []byte) Script {
	length, b := read(b, 1)
	script := *new(Script)
	count := 0
	var (
		current []byte
		data    []byte
	)

	for count < int(length[0]) {
		current, b = read(b, 1)
		count += 1
		if current[0] >= 0x01 && current[0] <= 0x4b {
			n := current[0]
			data, b = read(b, int(n))
			script = append(script, elem{
				value: n,
				data:  data,
			})
			count += int(n)
		} else if toOpcode(current[0]) == OP_PUSHDATA1 {
			data, b = read(b, 1)
			dataLen := binary.LittleEndian.Uint16(data)
			data, b = read(b, int(dataLen))
			script = append(script, elem{
				value: current[0],
				data:  data,
			})
			count += int(dataLen) + 1
		} else if toOpcode(current[0]) == OP_PUSHDATA2 {
			data, b = read(b, 2)
			dataLen := binary.LittleEndian.Uint16(data)
			data, b = read(b, int(dataLen))
			script = append(script, elem{
				value: current[0],
				data:  data,
			})
			count += int(dataLen) + 2
		} else if toOpcode(current[0]) == OP_PUSHDATA4 {
			data, b = read(b, 4)
			dataLen := binary.LittleEndian.Uint16(data)
			data, b = read(b, int(dataLen))
			script = append(script, elem{
				value: current[0],
				data:  data,
			})
			count += int(dataLen) + 4
		} else {
			script = append(script, elem{
				value: current[0],
			})
		}
	}

	return script
}

func (s Script) Serialize() ([]byte, error) {
	b := *new([]byte)

	for _, e := range s {
		if toOpcode(e.value) == OP_PUSHDATA1 {
			//TODO
		} else if toOpcode(e.value) == OP_PUSHDATA2 {
			//TODO
		} else if toOpcode(e.value) == OP_PUSHDATA4 {
			//TODO
		} else {
			b = append(b, e.value)
			b = append(b, e.data...)
		}
	}

	length, err := varint.Encode(uint64(len(b)))
	if err != nil {
		return nil, err
	}

	return append(length, b...), nil
}

type opcode byte

func toOpcode(b byte) opcode {
	return opcode(b)
}

const (
	OP_0         opcode = 0x00
	OP_FALSE     opcode = OP_0
	OP_PUSHDATA1 opcode = 0x4c
	OP_PUSHDATA2 opcode = 0x4d
	OP_PUSHDATA4 opcode = 0x4e
	OP_1NEGATE   opcode = 0x4f
	OP_RESERVED  opcode = 0x50
	OP_1         opcode = 0x51
	OP_TRUE      opcode = 0x51
	OP_2         opcode = 0x52
	OP_3         opcode = 0x53
	OP_4         opcode = 0x54
	OP_5         opcode = 0x55
	OP_6         opcode = 0x56
	OP_7         opcode = 0x57
	OP_8         opcode = 0x58
	OP_9         opcode = 0x59
	OP_10        opcode = 0x5a
	OP_11        opcode = 0x5b
	OP_12        opcode = 0x5c
	OP_13        opcode = 0x5d
	OP_14        opcode = 0x5e
	OP_15        opcode = 0x5f
	OP_16        opcode = 0x60
	OP_NOP       opcode = 0x61
	OP_VER       opcode = 0x62
	OP_IF        opcode = 0x63
)

var opcodes = map[string]opcode{
	"OP_0":         OP_0,
	"OP_FALSE":     OP_FALSE,
	"OP_PUSHDATA1": OP_PUSHDATA1,
	"OP_PUSHDATA2": OP_PUSHDATA2,
	"OP_PUSHDATA4": OP_PUSHDATA4,
	"OP_1NEGATE":   OP_1NEGATE,
	"OP_RESERVED":  OP_RESERVED,
	"OP_1":         OP_1,
	"OP_TRUE":      OP_TRUE,
	"OP_2":         OP_2,
	"OP_3":         OP_3,
	"OP_4":         OP_4,
	"OP_5":         OP_5,
	"OP_6":         OP_6,
	"OP_7":         OP_7,
	"OP_8":         OP_8,
	"OP_9":         OP_9,
	"OP_10":        OP_10,
	"OP_11":        OP_11,
	"OP_12":        OP_12,
	"OP_13":        OP_13,
	"OP_14":        OP_14,
	"OP_15":        OP_15,
	"OP_16":        OP_16,
}
