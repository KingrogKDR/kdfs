package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func check(e error) {
	if e != nil && e != io.EOF {
		panic(e)
	}
}

func main() {
	path := filepath.Join("fs_dd.img")

	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	check(err)
	defer f.Close()

	sb := Superblock{
		MagicNumber:      0xDEADAAA,
		BlockSize:        4096,
		InodeCount:       128,
		BitmapOffset:     4096,
		InodeTableOffset: 8192,
		DataRegionOffset: 16384,
	}

	buf := make([]byte, 24)
	binary.LittleEndian.PutUint32(buf[0:], sb.MagicNumber)
	binary.LittleEndian.PutUint32(buf[4:], sb.BlockSize)
	binary.LittleEndian.PutUint32(buf[8:], sb.InodeCount)
	binary.LittleEndian.PutUint32(buf[12:], sb.BitmapOffset)
	binary.LittleEndian.PutUint32(buf[16:], sb.InodeTableOffset)
	binary.LittleEndian.PutUint32(buf[20:], sb.DataRegionOffset)

	nWrite, err := f.WriteAt(buf, 0)
	check(err)

	readBuf := make([]byte, 24)

	_, err = f.ReadAt(readBuf, 0)
	check(err)

	sbRead := Superblock{
		MagicNumber:      binary.LittleEndian.Uint32(buf[0:]),
		BlockSize:        binary.LittleEndian.Uint32(buf[4:]),
		InodeCount:       binary.LittleEndian.Uint32(buf[8:]),
		BitmapOffset:     binary.LittleEndian.Uint32(buf[12:]),
		InodeTableOffset: binary.LittleEndian.Uint32(buf[16:]),
		DataRegionOffset: binary.LittleEndian.Uint32(buf[20:]),
	}

	fmt.Println("bytes written:", nWrite)
	fmt.Printf("%+v\n", sbRead)

}
