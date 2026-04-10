package metadata

import (
	"bytes"
	"encoding/binary"
	"os"
	"time"

	"github.com/KingrogKDR/kdfs/custom_error"
)

type KnodeType uint32

const (
	Free KnodeType = iota
	File
	Dir
)

type Knode struct {
	Typ       KnodeType
	Size      uint32
	LinkCount uint32

	Atime int64
	Mtime int64
	Ctime int64

	Direct [12]uint32
}

const (
	KnodeCount   uint32 = 128
	KnodeSize    uint32 = 128
	InvalidKnode        = ^uint32(0)
)

func (kn *Knode) AllocKnode(km *Bitmap, f *os.File, typ KnodeType, knodeStartOff uint32) (uint32, error) {
	i, err := km.AllocBit()
	if err != nil {
		return InvalidKnode, custom_error.WrapIO("alloc knode: can't alloc bit in knode-bitmap", f.Name(), err)
	}

	if err := kn.ReadKnode(f, knodeStartOff, i); err != nil {
		km.FreeBit(i)
		return InvalidKnode, custom_error.WrapIO("alloc knode: can't read knode", f.Name(), err)
	}

	if kn.Typ != Free {
		km.FreeBit(i)
		return InvalidKnode, custom_error.Corrupt("alloc knode", "bitmap/inode mismatch")
	}

	now := time.Now().Unix()

	kn.Typ = typ
	kn.Size = 0
	kn.LinkCount = 1
	kn.Atime = now
	kn.Mtime = now
	kn.Ctime = now

	if err := kn.WriteKnode(f, knodeStartOff, i); err != nil {
		km.FreeBit(i)
		return InvalidKnode, custom_error.WrapIO("alloc knode: can't write knode", f.Name(), err)
	}

	return i, nil
}

func (kn *Knode) FreeKnode(km *Bitmap, f *os.File, knodeIndex, knodeStartOff uint32) error {
	if err := kn.ReadKnode(f, knodeStartOff, knodeIndex); err != nil {
		return custom_error.WrapIO("free knode", f.Name(), err)
	}

	*kn = Knode{}

	err := kn.WriteKnode(f, knodeStartOff, knodeIndex)

	if err != nil {
		return custom_error.Corrupt("free knode", "couldn't update knode state")
	}

	km.FreeBit(knodeIndex - 1)

	return nil
}

func (kn *Knode) WriteKnode(f *os.File, knodeStartOff, knodeIndex uint32) error {
	if knodeIndex >= KnodeCount {
		return custom_error.Corrupt("write knode:", "knode index is out of bounds")
	}

	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, kn); err != nil {
		return custom_error.WrapIO("encode knode", f.Name(), err)
	}

	data := buf.Bytes()

	if len(data) > int(KnodeSize) {
		return custom_error.Corrupt("can't write knode:", "data exceeds knode size")
	}

	padded := make([]byte, KnodeSize) // makes sure Knode is equal to knodeSize
	copy(padded, data)

	offset := knodeStartOff + knodeIndex*KnodeSize
	_, err := f.WriteAt(padded, int64(offset))
	if err != nil {
		return custom_error.WrapIO("write knode", f.Name(), err)
	}

	return nil
}

func (kn *Knode) ReadKnode(f *os.File, knodeStartOff, knodeIndex uint32) error {
	buf := make([]byte, KnodeSize)

	offset := knodeStartOff + knodeIndex*KnodeSize
	_, err := f.ReadAt(buf, int64(offset))
	if err != nil {
		return custom_error.WrapIO("read knode", f.Name(), err)
	}

	reader := bytes.NewReader(buf)
	if err := binary.Read(reader, binary.LittleEndian, kn); err != nil {
		return custom_error.WrapIO("decode knode", f.Name(), err)
	}

	if kn.Typ > Dir {
		return custom_error.Corrupt("can't read knode:", "invalid type")
	}

	return nil
}
