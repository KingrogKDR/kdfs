package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/KingrogKDR/kdfs/custom_error"
	"github.com/KingrogKDR/kdfs/metadata"
)

const (
	inodeCount uint32 = 128
	inodeSize  uint32 = 128
)

func main() {
	path := filepath.Join("disk.img")

	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	custom_error.Check(custom_error.Wrap("open", path, err))
	defer f.Close()

	disk := metadata.Disk{
		BlockSize:  4096, // 4kb
		BlockCount: 2560, // ~10mb
	}

	layout, err := metadata.ComputeLayout(disk, inodeCount, inodeSize)
	custom_error.Check(err)

	sb := metadata.Superblock{
		MagicNumber: metadata.MagicFs,
		BlockSize:   disk.BlockSize,
		BlockCount:  disk.BlockCount,
		InodeCount:  inodeCount,

		BitmapStart: layout.BitmapStart,
		InodeStart:  layout.InodeStart,
		DataStart:   layout.DataStart,
	}

	err = custom_error.FileSystemCheck(sb.MagicNumber, metadata.MagicFs)
	custom_error.Check(err)

	_, err = metadata.WriteSuperblock(f, &sb)
	custom_error.Check(err)

	sbRead, err := metadata.ReadSuperblock(f, 0)
	custom_error.Check(err)

	fmt.Printf("%+v\n", sbRead)

}
