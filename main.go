package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/KingrogKDR/kdfs/custom_error"
	"github.com/KingrogKDR/kdfs/metadata"
)

const (
	knodeCount uint32 = 128
	knodeSize  uint32 = 128
)

func main() {
	path := filepath.Join("disk.img")

	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		custom_error.Check(custom_error.WrapIO("open disk", path, err))
	}
	defer f.Close()

	disk := metadata.Disk{
		BlockSize:  4096, // 4kb
		BlockCount: 2560, // ~10mb
	}

	layout, err := metadata.ComputeLayout(disk, knodeCount, knodeSize)
	custom_error.Check(err)

	sb := metadata.Superblock{
		MagicNumber: metadata.MagicFs,
		BlockSize:   disk.BlockSize,
		BlockCount:  disk.BlockCount,
		KnodeCount:  knodeCount,

		BitmapStart: layout.BitmapStart,
		KnodeStart:  layout.KnodeStart,
		DataStart:   layout.DataStart,
	}

	err = custom_error.FileSystemCheck(sb.MagicNumber, metadata.MagicFs)
	custom_error.Check(err)

	_, err = metadata.WriteSuperblock(f, &sb)
	custom_error.Check(err)

	bitmapSize := (sb.BlockCount + 7) / 8
	bitmapBlock := make([]byte, bitmapSize)

	_, err = f.ReadAt(bitmapBlock, int64(sb.BitmapStart*sb.BlockSize))
	if err != nil && err.Error() != "EOF" {
		custom_error.Check(err)
	}

	knodeBitmapSize := (knodeCount + 7) / 8

	_ = metadata.Bitmap{
		Data: bitmapBlock[:knodeBitmapSize],
		Base: 0,
	}

	_ = metadata.Bitmap{
		Data: bitmapBlock[knodeBitmapSize:bitmapSize],
		Base: sb.DataStart,
	}

	_ = metadata.Knode{
		Typ:       metadata.File,
		Size:      1024,
		LinkCount: 10,
		Blocks:    [12]uint32{},
	}

	_, err = f.WriteAt(bitmapBlock, int64(sb.BitmapStart*sb.BlockSize))
	custom_error.Check(err)

	fmt.Println("Bitmap written to disk")
}
