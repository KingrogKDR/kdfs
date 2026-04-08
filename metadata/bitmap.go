package metadata

import (
	"github.com/KingrogKDR/kdfs/custom_error"
)

// 1 block = 1 bit
// 1 byte = 8 blocks
// knode_bitmap : 0 -> 5
// data_bitmap : 6 -> end
// both bitmaps exist in the same block

type Bitmap struct {
	Data []byte
	Base uint32
}

func (b *Bitmap) Reset() {
	for i := range b.Data {
		b.Data[i] = 0
	}
}

func (b *Bitmap) AllocBit() (uint32, error) {
	maxBits := uint32(len(b.Data)) * 8
	for i := range maxBits {
		idx := b.Base + i

		if b.getBit(idx) == 0 {
			b.setBit(idx)
			return idx, nil
		}
	}

	return 0, custom_error.NoSpace("bitmap full")
}

func (b *Bitmap) FreeBit(block uint32) error {
	if block < b.Base {
		return custom_error.Corrupt("free bit", "out of range")
	}
	if b.getBit(block) == 0 {
		return custom_error.Corrupt("free bit from bitmap", "block is not allocated")
	}

	b.clearBit(block)

	return nil
}

func (b *Bitmap) getBit(index uint32) byte {
	index -= b.Base
	byteIndex := index / 8
	bitOffset := index % 8
	if byteIndex >= uint32(len(b.Data)) {
		panic("byte index out of range")
	}
	return (b.Data[byteIndex] >> bitOffset) & 1
}

func (b *Bitmap) setBit(index uint32) {
	index -= b.Base
	byteIndex := index / 8
	bitOffset := index % 8
	if byteIndex >= uint32(len(b.Data)) {
		panic("byte index out of range")
	}

	b.Data[byteIndex] |= (1 << bitOffset)
}

func (b *Bitmap) clearBit(index uint32) {
	index -= b.Base
	byteIndex := index / 8
	bitOffset := index % 8
	if byteIndex >= uint32(len(b.Data)) {
		panic("byte index out of range")
	}
	b.Data[byteIndex] &^= (1 << bitOffset) // first `not` the mask then `and` both
}
