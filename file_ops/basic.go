package fileops

import (
	"os"

	"github.com/KingrogKDR/kdfs/custom_error"
	"github.com/KingrogKDR/kdfs/metadata"
)

func FsCreateFile(f *os.File, km *metadata.Bitmap, knodeStartOff uint32) (uint32, error) {
	kn := metadata.NewKnode(metadata.File)
	kIndex, err := kn.AllocKnode(f, km, knodeStartOff)
	if err != nil {
		return metadata.InvalidKnode, custom_error.New(custom_error.TypeInvalid, "file create:", "file_ops/basic.go", "couldn't alloc knode", err)
	}

	return kIndex, nil
}

// func FsCreateDir(f *os.File, km *metadata.Bitmap, knodeStartOff uint32) (uint32, error) {
// 	kn := metadata.NewKnode(metadata.File)
// 	kIndex, err := kn.AllocKnode(f, km, knodeStartOff)
// 	if err != nil {
// 		return metadata.InvalidKnode, custom_error.New(custom_error.TypeInvalid, "file create:", "file_ops/basic.go", "couldn't alloc knode", err)
// 	}

// 	return kIndex, nil
// }
