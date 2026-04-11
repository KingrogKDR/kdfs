package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/KingrogKDR/kdfs/custom_error"
	"github.com/KingrogKDR/kdfs/metadata"
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

	layout, err := metadata.ComputeLayout(disk, metadata.KnodeCount, metadata.KnodeSize)
	custom_error.Check(err)

	sb := metadata.Superblock{
		MagicNumber: metadata.MagicFs,
		BlockSize:   disk.BlockSize,
		BlockCount:  disk.BlockCount,
		KnodeCount:  metadata.KnodeCount,

		BitmapStart: layout.BitmapStart,
		KnodeStart:  layout.KnodeStart,
		DataStart:   layout.DataStart,
	}

	err = custom_error.FileSystemCheck(sb.MagicNumber, metadata.MagicFs)
	custom_error.Check(err)

	_, err = metadata.WriteSuperblock(f, &sb)
	custom_error.Check(err)

	dataBitmapSize := (sb.BlockCount + 7) / 8

	kBitmapSize := (metadata.KnodeCount + 7) / 8

	dataBitmap := make([]byte, dataBitmapSize)
	kBitmap := make([]byte, kBitmapSize)

	dataBitmapOffset := sb.BitmapStart * sb.BlockSize
	kBitmapOffset := dataBitmapOffset + dataBitmapSize

	dm := &metadata.Bitmap{
		Name:       "Data-Bitmap",
		Data:       dataBitmap,
		DiskOffset: dataBitmapOffset,
		BitStart:   sb.DataStart,
		BitEnd:     sb.BlockCount - 1,
	}

	km := &metadata.Bitmap{
		Name:       "Knode-Bitmap",
		Data:       kBitmap,
		DiskOffset: kBitmapOffset,
		BitStart:   0,
		BitEnd:     metadata.KnodeCount - 1,
	}

	err = dm.Read(f)
	custom_error.Check(err)

	err = km.Read(f)
	custom_error.Check(err)

	metadata.ReserveMetaBlocks(dm, sb.DataStart)

	// knodeStartOff := sb.KnodeStart * sb.BlockSize
	// fileops.FsCreateFile(f, km, knodeStartOff)
	// fmt.Println("File created on disk")

	// km.Reset(f)
	// dm.Reset(f)

	fmt.Println(km)
	fmt.Println(dm)

}
