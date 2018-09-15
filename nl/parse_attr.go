package nl

import (
	"encoding/binary"
	"fmt"
)

type Attribute struct {
	Type  uint16
	Value []byte
}

func ParseAttributes(data []byte) <-chan Attribute {
	native := NativeEndian()
	result := make(chan Attribute)

	go func() {
		i := 0
		for i+4 < len(data) {
			length := int(native.Uint16(data[i : i+2]))

			result <- Attribute{
				Type:  native.Uint16(data[i+2 : i+4]),
				Value: data[i+4 : i+length],
			}
			i += rtaAlignOf(length)
		}
		close(result)
	}()

	return result
}

func PrintAttributes(data []byte) {
	printAttributes(data, 0)
}

func printAttributes(data []byte, level int) {
	for attr := range ParseAttributes(data) {
		for i := 0; i < level; i++ {
			print("> ")
		}
		nested := attr.Type&NLA_F_NESTED != 0
		fmt.Printf("type=%d nested=%v len=%v %v\n", attr.Type&NLA_TYPE_MASK, nested, len(attr.Value), attr.Value)
		if nested {
			printAttributes(attr.Value, level+1)
		}
	}
}

// Uint32 returns the uint32 value respecting the NET_BYTEORDER flag
func (attr *Attribute) Uint32() uint32 {
	if attr.Type&NLA_F_NET_BYTEORDER != 0 {
		return binary.BigEndian.Uint32(attr.Value)
	} else {
		return NativeEndian().Uint32(attr.Value)
	}
}
