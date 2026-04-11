package metadata

import (
	"encoding/binary"
	"os"

	"github.com/KingrogKDR/kdfs/custom_error"
)

const (
	magicOffset      = 0
	blockSizeOffset  = 4
	blockCountOffset = 8
	knodeCountOffset = 12
	bitmapStartOff   = 16
	knodeStartOff    = 20
	dataStartOff     = 24

	MagicFs = 0xDEADAAA
)

// encoding/binary ignores unexported fields
type Superblock struct {
	MagicNumber uint32

	BlockSize  uint32
	BlockCount uint32

	KnodeCount uint32

	BitmapStart uint32
	KnodeStart  uint32
	DataStart   uint32
}

func ReadSuperblock(f *os.File, offset int64, blockSize uint32) (Superblock, error) {
	buf := make([]byte, blockSize)

	_, err := f.ReadAt(buf, offset)
	if err != nil {
		return Superblock{}, custom_error.WrapIO("read superblock", f.Name(), err)
	}

	sb := Superblock{
		MagicNumber: binary.LittleEndian.Uint32(buf[magicOffset:]),
		BlockSize:   binary.LittleEndian.Uint32(buf[blockSizeOffset:]),
		BlockCount:  binary.LittleEndian.Uint32(buf[blockCountOffset:]),
		KnodeCount:  binary.LittleEndian.Uint32(buf[knodeCountOffset:]),
		BitmapStart: binary.LittleEndian.Uint32(buf[bitmapStartOff:]),
		KnodeStart:  binary.LittleEndian.Uint32(buf[knodeStartOff:]),
		DataStart:   binary.LittleEndian.Uint32(buf[dataStartOff:]),
	}

	return sb, nil
}

func WriteSuperblock(f *os.File, sb *Superblock) (int, error) {
	buf := make([]byte, sb.BlockSize)

	binary.LittleEndian.PutUint32(buf[magicOffset:], sb.MagicNumber)
	binary.LittleEndian.PutUint32(buf[blockSizeOffset:], sb.BlockSize)
	binary.LittleEndian.PutUint32(buf[blockCountOffset:], sb.BlockCount)
	binary.LittleEndian.PutUint32(buf[knodeCountOffset:], sb.KnodeCount)
	binary.LittleEndian.PutUint32(buf[bitmapStartOff:], sb.BitmapStart)
	binary.LittleEndian.PutUint32(buf[knodeStartOff:], sb.KnodeStart)
	binary.LittleEndian.PutUint32(buf[dataStartOff:], sb.DataStart)

	nBytes, err := f.WriteAt(buf, 0)
	if err != nil {
		return 0, custom_error.WrapIO("write superblock", f.Name(), err)
	}
	return nBytes, nil
}
