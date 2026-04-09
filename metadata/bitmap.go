package metadata

import (
	"io"
	"os"

	"github.com/KingrogKDR/kdfs/custom_error"
)

// 1 block = 1 bit
// 1 byte = 8 blocks

type Bitmap struct {
	Name     string
	Data     []byte
	StartOff uint32
	EndOff   uint32
}

func ReserveMetaBlocks(bm *Bitmap, dataStartOff uint32) {
	for i := range dataStartOff {
		bm.setBit(i)
	}
}

func (b *Bitmap) Reset() {
	for i := b.StartOff; i < uint32(len(b.Data)*8); i++ {
		b.clearBit(i)
	}
}

func (b *Bitmap) AllocBit() (uint32, error) {
	for i := b.StartOff; i <= b.EndOff; i++ {
		if b.getBit(i) == 0 {
			b.setBit(i)
			return i, nil
		}
	}
	return 0, custom_error.NoSpace("Bitmap is full")
}

func (b *Bitmap) FreeBit(block uint32) error {
	if block < b.StartOff || block > b.EndOff {
		return custom_error.Corrupt("free bit", "out of range")
	}
	if b.getBit(block) == 0 {
		return custom_error.Corrupt("free bit", "block is not allocated")
	}

	b.clearBit(block)

	return nil
}

func (b *Bitmap) Read(f *os.File, offset uint32) (int, error) {
	n, err := f.ReadAt(b.Data, int64(offset))
	if err != nil && err != io.EOF {
		return 0, custom_error.WrapIO("read bitmap from file", b.Name, err)
	}
	return n, nil
}

func (b *Bitmap) Write(f *os.File, offset uint32) (int, error) {
	n, err := f.WriteAt(b.Data, int64(offset))
	if err != nil {
		return 0, custom_error.WrapIO("write bitmap from file", b.Name, err)
	}
	return n, nil
}

func (b *Bitmap) getBit(block uint32) byte {
	byteIndex := block / 8
	bitOffset := block % 8
	if byteIndex >= uint32(len(b.Data)) {
		custom_error.Exception("get bit", "byteIndex outside of range")
	}

	return (b.Data[byteIndex] >> bitOffset) & 1
}

func (b *Bitmap) setBit(block uint32) {
	byteIndex := block / 8
	bitOffset := block % 8
	if byteIndex >= uint32(len(b.Data)) {
		custom_error.Exception("set bit", "byteIndex outside of range")
	}

	b.Data[byteIndex] |= (1 << bitOffset)
}

func (b *Bitmap) clearBit(block uint32) {
	byteIndex := block / 8
	bitOffset := block % 8
	if byteIndex >= uint32(len(b.Data)) {
		custom_error.Exception("clear bit", "byteIndex outside of range")
	}
	b.Data[byteIndex] &^= (1 << bitOffset) // first `not` the mask then `and` both
}
