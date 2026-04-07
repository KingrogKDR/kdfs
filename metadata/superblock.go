package metadata

import (
	"encoding/binary"
	"os"

	"github.com/KingrogKDR/kdfs/custom_error"
)

// encoding/binary ignores unexported fields

const (
	SuperblockSize = 28 // 7 fields × 4 bytes (uint32)

	MagicOffset      = 0
	BlockSizeOffset  = 4
	BlockCountOffset = 8
	InodeCountOffset = 12
	BitmapStartOff   = 16
	InodeStartOff    = 20
	DataStartOff     = 24

	MagicFs = 0xDEADAAA
)

type Superblock struct {
	MagicNumber uint32

	BlockSize  uint32
	BlockCount uint32

	InodeCount uint32

	BitmapStart uint32
	InodeStart  uint32
	DataStart   uint32
}

func ReadSuperblock(f *os.File, offset int64) (Superblock, error) {
	buf := make([]byte, SuperblockSize)

	_, err := f.ReadAt(buf, offset)
	if err != nil {
		return Superblock{}, custom_error.WrapIO("read superblock", f.Name(), err)
	}

	sb := Superblock{
		MagicNumber: binary.LittleEndian.Uint32(buf[MagicOffset:]),
		BlockSize:   binary.LittleEndian.Uint32(buf[BlockSizeOffset:]),
		BlockCount:  binary.LittleEndian.Uint32(buf[BlockCountOffset:]),
		InodeCount:  binary.LittleEndian.Uint32(buf[InodeCountOffset:]),
		BitmapStart: binary.LittleEndian.Uint32(buf[BitmapStartOff:]),
		InodeStart:  binary.LittleEndian.Uint32(buf[InodeStartOff:]),
		DataStart:   binary.LittleEndian.Uint32(buf[DataStartOff:]),
	}

	return sb, nil
}

func WriteSuperblock(f *os.File, sb *Superblock) (int, error) {
	buf := make([]byte, SuperblockSize)

	binary.LittleEndian.PutUint32(buf[MagicOffset:], sb.MagicNumber)
	binary.LittleEndian.PutUint32(buf[BlockSizeOffset:], sb.BlockSize)
	binary.LittleEndian.PutUint32(buf[BlockCountOffset:], sb.BlockCount)
	binary.LittleEndian.PutUint32(buf[InodeCountOffset:], sb.InodeCount)
	binary.LittleEndian.PutUint32(buf[BitmapStartOff:], sb.BitmapStart)
	binary.LittleEndian.PutUint32(buf[InodeStartOff:], sb.InodeStart)
	binary.LittleEndian.PutUint32(buf[DataStartOff:], sb.DataStart)

	nBytes, err := f.WriteAt(buf, 0)
	if err != nil {
		return 0, custom_error.WrapIO("write superblock", f.Name(), err)
	}
	return nBytes, nil
}
