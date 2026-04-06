package metadata

import (
	"encoding/binary"
	"os"

	"github.com/KingrogKDR/kdfs/custom_error"
)

// encoding/binary ignores unexported fields

const (
	SuperblockSize = 24

	MagicOffset     = 0
	BlockSizeOffset = 4
	InodeOffset     = 8
	BitmapOffsetOff = 12
	InodeTableOff   = 16
	DataRegionOff   = 20
)

type Superblock struct {
	MagicNumber      uint32
	BlockSize        uint32
	InodeCount       uint32
	BitmapOffset     uint32
	InodeTableOffset uint32
	DataRegionOffset uint32
}

func ReadSuperblock(f *os.File, offset int64) Superblock {
	readBuf := make([]byte, SuperblockSize)

	_, err := f.ReadAt(readBuf, offset)
	custom_error.Check(err)

	sbRead := Superblock{
		MagicNumber:      binary.LittleEndian.Uint32(readBuf[0:4]),
		BlockSize:        binary.LittleEndian.Uint32(readBuf[4:8]),
		InodeCount:       binary.LittleEndian.Uint32(readBuf[8:12]),
		BitmapOffset:     binary.LittleEndian.Uint32(readBuf[12:16]),
		InodeTableOffset: binary.LittleEndian.Uint32(readBuf[16:20]),
		DataRegionOffset: binary.LittleEndian.Uint32(readBuf[20:24]),
	}

	return sbRead
}

func WriteSuperblock(f *os.File, sb *Superblock) (int, error) {
	buf := make([]byte, SuperblockSize)

	binary.LittleEndian.PutUint32(buf[MagicOffset:], sb.MagicNumber)
	binary.LittleEndian.PutUint32(buf[BlockSizeOffset:], sb.BlockSize)
	binary.LittleEndian.PutUint32(buf[InodeOffset:], sb.InodeCount)
	binary.LittleEndian.PutUint32(buf[BitmapOffsetOff:], sb.BitmapOffset)
	binary.LittleEndian.PutUint32(buf[InodeTableOff:], sb.InodeTableOffset)
	binary.LittleEndian.PutUint32(buf[DataRegionOff:], sb.DataRegionOffset)

	nBytes, err := f.WriteAt(buf, 0)
	return nBytes, err
}
