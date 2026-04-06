package main

import (
	"os"
	"path/filepath"

	"github.com/KingrogKDR/kdfs/custom_error"
	"github.com/KingrogKDR/kdfs/metadata"
)

func main() {
	path := filepath.Join("fs_dd.img")

	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	custom_error.Check(err)
	defer f.Close()

	sb := metadata.Superblock{
		MagicNumber:      0xDEADAAA,
		BlockSize:        4096,
		InodeCount:       128,
		BitmapOffset:     4096,
		InodeTableOffset: 8192,
		DataRegionOffset: 16384,
	}

	_, err = metadata.WriteSuperblock(f, &sb)
	custom_error.Check(err)

	_ = metadata.ReadSuperblock(f, 0)

}
