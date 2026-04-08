package metadata

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/KingrogKDR/kdfs/custom_error"
)

type KnodeType uint32

const (
	File KnodeType = iota
	Dir
)

type Knode struct {
	Typ       KnodeType
	Size      uint32
	LinkCount uint32
	Blocks    [12]uint32
}

func (kn *Knode) WriteKnode(f *os.File, knodeStartOff, knodeIndex, knodeSize uint32) error {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, kn); err != nil {
		return custom_error.WrapIO("encode knode", f.Name(), err)
	}

	data := buf.Bytes()
	if len(data) > int(knodeSize) {
		return custom_error.Corrupt("knode", "exceeds knode size")
	}

	padded := make([]byte, knodeSize) // makes sure Knode is equal to knodeSize
	copy(padded, data)

	offset := knodeStartOff + knodeIndex*knodeSize
	_, err := f.WriteAt(padded, int64(offset))
	if err != nil {
		return custom_error.WrapIO("write knode", f.Name(), err)
	}

	return nil
}

func (kn *Knode) ReadKnode(f *os.File, knodeStartOff, knodeIndex, knodeSize uint32) error {
	buf := make([]byte, knodeSize)

	offset := knodeStartOff + knodeIndex*knodeSize
	_, err := f.ReadAt(buf, int64(offset))
	if err != nil {
		return custom_error.WrapIO("read knode", f.Name(), err)
	}

	reader := bytes.NewReader(buf)
	if err := binary.Read(reader, binary.LittleEndian, kn); err != nil {
		return custom_error.WrapIO("decode knode", f.Name(), err)
	}

	return nil
}
