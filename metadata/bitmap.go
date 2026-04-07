package metadata

import (
	"fmt"
	"os"

	"github.com/KingrogKDR/kdfs/custom_error"
)

// 1 block = 1 bit
// 1 byte = 8 blocks

type Bitmap struct {
	Data      []byte
	DataStart uint32
}

func (b *Bitmap) WriteToFile(f *os.File, bitmapOffset uint32) error {
	_, err := f.WriteAt(b.Data, int64(bitmapOffset))
	if err != nil {
		return custom_error.WrapIO("write bitmap", f.Name(), err)
	}
	return nil
}

func (b *Bitmap) AllocBit(totalBlocks uint32) (uint32, error) {
	for i := b.DataStart; i < totalBlocks; i++ {
		if b.getBit(i) == 0 {
			b.setBit(i)
			return i, nil
		}
	}
	return 0, custom_error.NoSpace("bitmap alloc")
}

func (b *Bitmap) FreeBit(block uint32) error {
	if block < b.DataStart {
		return custom_error.Corrupt("free bit from bitmap", "metadata block cannot be freed")
	}

	if b.getBit(block) == 0 {
		return custom_error.Corrupt("free bit from bitmap", "block is not allocated")
	}

	b.clearBit(block)

	return nil
}

func (b *Bitmap) SetMetaBlocks() {
	for i := uint32(0); i < b.DataStart; i++ {
		b.setBit(i)
	}

	fmt.Println("Metadata blocks set to 1 in bitmap")
}

func (b *Bitmap) getBit(block uint32) byte {
	byteIndex := block / 8
	bitOffset := block % 8
	if byteIndex >= uint32(len(b.Data)) {
		panic("bitmap out of range")
	}
	return (b.Data[byteIndex] >> bitOffset) & 1
}

func (b *Bitmap) setBit(block uint32) {
	byteIndex := block / 8
	bitOffset := block % 8
	if byteIndex >= uint32(len(b.Data)) {
		panic("bitmap out of range")
	}

	b.Data[byteIndex] |= (1 << bitOffset)
}

func (b *Bitmap) clearBit(block uint32) {
	byteIndex := block / 8
	bitOffset := block % 8
	if byteIndex >= uint32(len(b.Data)) {
		panic("bitmap out of range")
	}
	b.Data[byteIndex] &^= (1 << bitOffset) // first `not` the mask then `and` both
}
