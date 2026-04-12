package fileops

import (
	"io"
	"os"
	"testing"

	"github.com/KingrogKDR/kdfs/metadata"
)

func newTestBitmap(name string, size uint32, start uint32) *metadata.Bitmap {
	byteSize := (size + 7) / 8

	return &metadata.Bitmap{
		Name:       name,
		Data:       make([]byte, byteSize),
		DiskOffset: 0,
		BitStart:   start,
		BitEnd:     start + size - 1,
	}
}

func initKnodeTable(disk *os.File, knodeStart uint32) error {
	var empty metadata.Knode

	for i := range metadata.KnodeCount {
		err := empty.WriteKnode(disk, knodeStart, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestWriteReadSingleBlock(t *testing.T) {
	disk, err := os.Create("test.img")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.img")
	defer disk.Close()

	// setup
	km := newTestBitmap("knode", metadata.KnodeCount, 0)
	dm := newTestBitmap("data", 3000, 6)

	knodeStart := uint32(0)
	dataStart := uint32(6 * 4096) // example
	blockSize := uint32(4096)

	// create file
	err = disk.Truncate(10 * 1024 * 1024) // 10 MB
	if err != nil {
		t.Fatal(err)
	}

	if err := initKnodeTable(disk, knodeStart); err != nil {
		t.Fatal(err)
	}

	kIndex, err := FsCreateFile(disk, km, knodeStart)
	if err != nil {
		t.Fatal(err)
	}

	input := []byte("hello filesystem")

	err = FsWriteFile(disk, input, km, dm, knodeStart, dataStart, kIndex, blockSize)
	if err != nil {
		t.Fatal(err)
	}

	// read back manually
	var kn metadata.Knode
	_ = kn.ReadKnode(disk, knodeStart, kIndex)

	block := kn.Direct[0]

	buf := make([]byte, len(input))

	offset := int64(dataStart) + int64(block)*int64(blockSize)
	if _, err := disk.Seek(offset, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	if _, err := disk.Read(buf); err != nil {
		t.Fatal(err)
	}

	if string(buf) != string(input) {
		t.Fatalf("data mismatch: got %s, want %s", buf, input)
	}
}

func TestAppend(t *testing.T) {
	disk, err := os.Create("test.img")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.img")
	defer disk.Close()

	// setup
	km := newTestBitmap("knode", metadata.KnodeCount, 0)
	dm := newTestBitmap("data", 3000, 6)

	knodeStart := uint32(0)
	dataStart := uint32(6 * 4096) // example
	blockSize := uint32(4096)

	err = disk.Truncate(10 * 1024 * 1024) // 10 MB
	if err != nil {
		t.Fatal(err)
	}

	if err := initKnodeTable(disk, knodeStart); err != nil {
		t.Fatal(err)
	}

	// create file
	kIndex, err := FsCreateFile(disk, km, knodeStart)
	if err != nil {
		t.Fatal(err)
	}

	first := []byte("hello ")
	second := []byte("world")

	if err := FsWriteFile(disk, first, km, dm, knodeStart, dataStart, kIndex, blockSize); err != nil {
		t.Fatal(err)
	}
	if err := FsWriteFile(disk, second, km, dm, knodeStart, dataStart, kIndex, blockSize); err != nil {
		t.Fatal(err)
	}

	expected := "hello world"
	// read full file
	var kn metadata.Knode
	kn.ReadKnode(disk, knodeStart, kIndex)

	buf := make([]byte, kn.Size)

	var read uint32 = 0

	for i := 0; i < 12 && read < kn.Size; i++ {
		if kn.Direct[i] == 0 {
			continue
		}

		offset := int64(dataStart) + int64(kn.Direct[i])*int64(blockSize)
		if _, err := disk.Seek(offset, io.SeekStart); err != nil {
			t.Fatal(err)
		}

		n := min(blockSize, kn.Size-read)
		if _, err := disk.Read(buf[read : read+n]); err != nil {
			t.Fatal(err)
		}
		read += n
	}

	if string(buf) != expected {
		t.Fatalf("append failed: got %s", buf)
	}
}
