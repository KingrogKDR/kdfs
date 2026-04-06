package metadata

import (
	"encoding/binary"
	"os"

	"github.com/KingrogKDR/kdfs/custom_error"
)

// encoding/binary ignores unexported fields

type Superblock struct {
	MagicNumber      uint32
	BlockSize        uint32
	InodeCount       uint32
	BitmapOffset     uint32
	InodeTableOffset uint32
	DataRegionOffset uint32
}

func ReadSuperblock(f *os.File, offset int64, bufSize int) Superblock {
	readBuf := make([]byte, bufSize)

	_, err := f.ReadAt(readBuf, offset)
	custom_error.Check(err)

	sbRead := Superblock{
		MagicNumber:      binary.LittleEndian.Uint32(readBuf[0:]),
		BlockSize:        binary.LittleEndian.Uint32(readBuf[4:]),
		InodeCount:       binary.LittleEndian.Uint32(readBuf[8:]),
		BitmapOffset:     binary.LittleEndian.Uint32(readBuf[12:]),
		InodeTableOffset: binary.LittleEndian.Uint32(readBuf[16:]),
		DataRegionOffset: binary.LittleEndian.Uint32(readBuf[20:]),
	}

	return sbRead
}

func WriteSuperblock(f *os.File, sb *Superblock, bufSize int) (int, error) {
	buf := make([]byte, 24)
	binary.LittleEndian.PutUint32(buf[0:], sb.MagicNumber)
	binary.LittleEndian.PutUint32(buf[4:], sb.BlockSize)
	binary.LittleEndian.PutUint32(buf[8:], sb.InodeCount)
	binary.LittleEndian.PutUint32(buf[12:], sb.BitmapOffset)
	binary.LittleEndian.PutUint32(buf[16:], sb.InodeTableOffset)
	binary.LittleEndian.PutUint32(buf[20:], sb.DataRegionOffset)

	nBytes, err := f.WriteAt(buf, 0)
	return nBytes, err
}
