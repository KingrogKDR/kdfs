package fileops

import (
	"io"
	"os"
	"time"

	"github.com/KingrogKDR/kdfs/custom_error"
	"github.com/KingrogKDR/kdfs/metadata"
)

func FsCreateFile(disk *os.File, km *metadata.Bitmap, knodeStartOff uint32) (uint32, error) {
	kn := metadata.NewKnode(metadata.File)
	kIndex, err := kn.AllocKnode(disk, km, knodeStartOff)
	if err != nil {
		return metadata.InvalidKnode, custom_error.New(custom_error.TypeInvalid, "file create:", "file_ops/basic.go", "couldn't alloc knode", err)
	}

	return kIndex, nil
}

func FsWriteFile(disk *os.File, data []byte, km, dm *metadata.Bitmap, knodeStartOff, dataStartOff, kIndex, blocksize uint32) error {
	// check if file exists
	var kn metadata.Knode

	err := kn.ReadKnode(disk, knodeStartOff, kIndex)

	if err != nil {
		return custom_error.WrapIO("file write:", "file_ops/basic.go\ncouldn't read knode", err)
	}

	blocksNeeded := (uint32(len(data)) + blocksize - 1) / blocksize
	if blocksNeeded > 12 {
		return custom_error.NoSpace("file too large")
	}

	// append to file

	blockIndex := kn.Size / blocksize
	blockOffset := kn.Size % blocksize

	written := uint32(0)
	remaining := uint32(len(data))

	if blockOffset != 0 {
		if blockIndex >= 12 {
			return custom_error.NoSpace("cannot append to file, it becomes too large")
		}

		// allocate if block doesn't exist
		if kn.Direct[blockIndex] == 0 {
			b, err := dm.AllocBit()
			if err != nil {
				return err
			}
			kn.Direct[blockIndex] = b
		}

		blockNum := kn.Direct[blockIndex]

		writeStart := dataStartOff + blockNum*blocksize + blockOffset

		if _, err := disk.Seek(int64(writeStart), io.SeekStart); err != nil {
			return err
		}

		space := blocksize - blockOffset
		toWrite := min(space, remaining)

		if _, err := disk.Write(data[written : written+toWrite]); err != nil {
			return err
		}

		written += toWrite
		remaining -= toWrite
		blockIndex++
	}

	for remaining > 0 {
		if blockIndex >= 12 {
			return custom_error.NoSpace("cannot append to file, it becomes too large")
		}

		b, err := dm.AllocBit()
		if err != nil {
			return err
		}
		kn.Direct[blockIndex] = b

		blockOff := int64(dataStartOff) + int64(b)*int64(blocksize)

		if _, err := disk.Seek(blockOff, io.SeekStart); err != nil {
			return err
		}

		toWrite := min(blocksize, remaining)

		if _, err := disk.Write(data[written : written+toWrite]); err != nil {
			return err
		}

		written += toWrite
		remaining -= toWrite
		blockIndex++
	}

	// update knode/inode

	kn.Size += uint32(len(data))
	kn.Mtime = time.Now().Unix()
	kn.Atime = kn.Mtime

	err = kn.WriteKnode(disk, knodeStartOff, kIndex)

	if err != nil {
		return err
	}

	return nil
}
