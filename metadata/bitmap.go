package metadata

import (
	"io"
	"os"

	"github.com/KingrogKDR/kdfs/custom_error"
)

// 1 block = 1 bit
// 1 byte = 8 blocks
// km range -> 0-127
// dm range -> 6-2559

type Bitmap struct {
	Name       string
	Data       []byte
	DiskOffset uint32

	BitStart uint32
	BitEnd   uint32
}

func ReserveMetaBlocks(bm *Bitmap, dataStartOff uint32) {
	for i := range dataStartOff {
		bm.setBit(i)
	}
}

func (b *Bitmap) Reset(f *os.File) {
	for i := b.BitStart; i <= b.BitEnd; i++ {
		b.clearBit(i)
	}
	err := b.Write(f)
	if err != nil {
		custom_error.Exception("reset bitmap:", "error writing bitmap to disk")
	}
}

func (b *Bitmap) AllocBit() (uint32, error) {
	for i := b.BitStart; i <= b.BitEnd; i++ {
		if b.getBit(i) == 0 {
			b.setBit(i)
			return i, nil
		}
	}
	return 0, custom_error.NoSpace("Bitmap is full")
}

func (b *Bitmap) FreeBit(block uint32) error {
	if block < b.BitStart || block > b.BitEnd {
		return custom_error.Corrupt("free bit", "out of range")
	}
	if b.getBit(block) == 0 {
		return custom_error.Corrupt("free bit", "block is not allocated")
	}

	b.clearBit(block)

	return nil
}

func (b *Bitmap) Read(f *os.File) error {
	_, err := f.ReadAt(b.Data, int64(b.DiskOffset))
	if err != nil && err != io.EOF {
		return custom_error.WrapIO("read bitmap from file", b.Name, err)
	}
	return nil
}

func (b *Bitmap) Write(f *os.File) error {
	_, err := f.WriteAt(b.Data, int64(b.DiskOffset))
	if err != nil {
		return custom_error.WrapIO("write bitmap from file", b.Name, err)
	}
	return nil
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
