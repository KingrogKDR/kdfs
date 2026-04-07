package metadata

import (
	"fmt"

	"github.com/KingrogKDR/kdfs/custom_error"
)

type Disk struct {
	BlockSize  uint32
	BlockCount uint32
}

type Layout struct {
	TotalBlocks uint32

	BitmapStart  uint32
	BitmapBlocks uint32

	InodeStart  uint32
	InodeBlocks uint32

	DataStart uint32
}

func ceilDiv(a, b uint32) uint32 {
	return (a + b - 1) / b
}

func ComputeLayout(d Disk, inodeCount, inodeSize uint32) (*Layout, error) {
	if d.BlockSize == 0 || d.BlockCount == 0 {
		return nil, custom_error.Corrupt("compute layout", "no space assigned to disk")
	}

	totalBlocks := d.BlockCount

	bitmapBytes := ceilDiv(totalBlocks, 8)
	bitmapBlocks := ceilDiv(bitmapBytes, d.BlockSize)

	inodeTableBytes := inodeCount * inodeSize
	inodeBlocks := ceilDiv(inodeTableBytes, d.BlockSize)

	superblockBlocks := uint32(1)
	bitmapStart := superblockBlocks
	inodeStart := bitmapStart + bitmapBlocks
	dataStart := inodeStart + inodeBlocks

	if dataStart >= totalBlocks {
		return nil, fmt.Errorf(
			"layout exceeds disk: dataStart=%d totalBlocks=%d: %w",
			dataStart, totalBlocks, custom_error.Corrupt("compute layout", "invalid layout, exceeds disk capacity"),
		)
	}

	return &Layout{
		TotalBlocks: totalBlocks,

		BitmapStart:  bitmapStart,
		BitmapBlocks: bitmapBlocks,

		InodeStart:  inodeStart,
		InodeBlocks: inodeBlocks,

		DataStart: dataStart,
	}, nil
}
