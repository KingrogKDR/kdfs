package custom_error

import (
	"errors"
	"fmt"
	"io"
	"log"
)

var (
	ErrInvalidSuperblock = errors.New("invalid superblock")
	ErrDiskNotFound      = errors.New("disk not found")
	ErrCorruptData       = errors.New("corrupt data")
)

type DiskError struct {
	Op   string
	Path string
	Err  error
}

func (e *DiskError) Error() string {
	return fmt.Sprintf("%s on %s: %v", e.Op, e.Path, e.Err)
}

func (e *DiskError) Unwrap() error {
	return e.Err
}

func Wrap(op, path string, err error) error {
	if err == nil {
		return nil
	}
	return &DiskError{
		Op:   op,
		Path: path,
		Err:  err,
	}
}

func Check(err error) {
	if err == nil || errors.Is(err, io.EOF) {
		return
	}
	log.Fatalf("fatal error: %+v", err)
}

func FileSystemCheck(magic, expected uint32) error {
	if magic != expected {
		return fmt.Errorf("superblock validation failed: %w", ErrInvalidSuperblock)
	}
	return nil
}
