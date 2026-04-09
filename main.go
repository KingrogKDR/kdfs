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

	dataBitmapSize := (sb.BlockCount + 7) / 8

	kBitmapSize := sb.BlockSize - dataBitmapSize // padding it to cover the block (4096-320)

	dataBitmap := make([]byte, dataBitmapSize)
	kBitmap := make([]byte, kBitmapSize)

	dm := metadata.Bitmap{
		Name:     "Data-Bitmap",
		Data:     dataBitmap,
		StartOff: sb.DataStart,
		EndOff:   sb.BlockCount - 1,
	}

	km := metadata.Bitmap{
		Name:     "Knode-Bitmap",
		Data:     kBitmap,
		StartOff: 0,
		EndOff:   knodeCount - 1,
	}

	dataBitmapOffset := sb.BitmapStart * sb.BlockSize
	kBitmapOffset := dataBitmapOffset + dataBitmapSize

	_, err = dm.Read(f, dataBitmapOffset)
	custom_error.Check(err)

	_, err = km.Read(f, kBitmapOffset)
	custom_error.Check(err)

	metadata.ReserveMetaBlocks(&dm, sb.DataStart)

	_, err = dm.Write(f, dataBitmapOffset)
	custom_error.Check(err)

	_, err = km.Write(f, kBitmapOffset)
	custom_error.Check(err)

	fmt.Println("Bitmap written to disk")
}
