package custom_error

import (
	"errors"
	"fmt"
	"io"
	"log"
)

type ErrorType string

const (
	TypeIO      ErrorType = "IO_ERROR"
	TypeCorrupt ErrorType = "CORRUPT_DATA"
	TypeInvalid ErrorType = "INVALID_STATE"
	TypeNoSpace ErrorType = "NO_SPACE"
)

type DiskError struct {
	Type ErrorType

	Op   string
	Path string

	Msg string
	Err error
}

func (e *DiskError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s on %s: %s: %v",
			e.Type, e.Op, e.Path, e.Msg, e.Err)
	}
	return fmt.Sprintf("[%s] %s on %s: %s",
		e.Type, e.Op, e.Path, e.Msg)
}

func (e *DiskError) Unwrap() error {
	return e.Err
}

func New(t ErrorType, op, path, msg string, err error) error {
	return &DiskError{
		Type: t,
		Op:   op,
		Path: path,
		Msg:  msg,
		Err:  err,
	}
}

func WrapIO(op, path string, err error) error {
	return New(TypeIO, op, path, "io failure", err)
}

func Corrupt(op, msg string) error {
	return New(TypeCorrupt, op, "", msg, nil)
}

func NoSpace(op string) error {
	return New(TypeNoSpace, op, "", "no free blocks available", nil)
}

func Check(err error) {
	if err == nil || errors.Is(err, io.EOF) {
		return
	}

	var dErr *DiskError
	if errors.As(err, &dErr) {
		log.Fatalf("%s", dErr.Error())
	}

	log.Fatalf("fatal: %v", err)
}

func FileSystemCheck(magic, expected uint32) error {
	if magic != expected {
		return Corrupt("superblock", "invalid magic number")
	}
	return nil
}
