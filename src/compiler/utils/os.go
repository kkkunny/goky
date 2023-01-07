package utils

import (
	stlos "github.com/kkkunny/stl/os"
	"golang.org/x/exp/constraints"
	"os"
	"path/filepath"
	"unsafe"
)

// PtrByte 指针大小
var PtrByte = uint(unsafe.Sizeof(uintptr(0)))

// AlignByte 对齐
var AlignByte = 4

// AlignTo 对齐
func AlignTo[T constraints.Integer | constraints.Float](n, align T) T {
	return (n + align - 1) / align * align
}

// GetRootPath 获取语言根目录
func GetRootPath() (stlos.Path, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return stlos.Path(filepath.Dir(exe)), err
}
