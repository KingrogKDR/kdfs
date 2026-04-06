package main

import (
	"fmt"
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

	sb := &metadata.Superblock{
		MagicNumber:      0xDEADAAA,
		BlockSize:        4096,
		InodeCount:       128,
		BitmapOffset:     4096,
		InodeTableOffset: 8192,
		DataRegionOffset: 16384,
	}

	nWrite, err := metadata.WriteSuperblock(f, sb, 24)
	custom_error.Check(err)

	fmt.Println("bytes written:", nWrite)
	fmt.Printf("%+v\n", metadata.ReadSuperblock(f, 0, 24))

}
