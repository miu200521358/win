package win

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	COMPRESS_ALGORITHM_MSZIP       = 2         // MSZIP 圧縮アルゴリズム
	COMPRESS_ALGORITHM_XPRESS      = 3         // XPRESS 圧縮アルゴリズム
	COMPRESS_ALGORITHM_XPRESS_HUFF = 4         // Huffman エンコードを使用した XPRESS 圧縮アルゴリズム
	COMPRESS_ALGORITHM_LZMS        = 5         // LZMS 圧縮アルゴリズム
	COMPRESS_RAW                   = 536870912 // 「CAB 形式」として圧縮されている
)

var (
	// Library
	libcabinet32 *windows.LazyDLL

	createDecompressor *windows.LazyProc
	closeDecompressor  *windows.LazyProc
	decompress         *windows.LazyProc
)

func init() {
	// Library
	libcabinet32 = windows.NewLazySystemDLL("cabinet.dll")

	// Functions
	createDecompressor = libcabinet32.NewProc("CreateDecompressor")
	closeDecompressor = libcabinet32.NewProc("CloseDecompressor")
	decompress = libcabinet32.NewProc("Decompress")
}

func CreateDecompressor(Algorithm int32, decompressorHandlePtr *HWND) (bool, error) {
	ret, _, errno := syscall.Syscall(createDecompressor.Addr(), 3,
		uintptr(Algorithm),
		0,
		uintptr(unsafe.Pointer(decompressorHandlePtr)),
	)

	return ret != 0, errno
}

func CloseDecompressor(decompressorHandle HWND) bool {
	ret, _, _ := syscall.Syscall(closeDecompressor.Addr(), 1,
		uintptr(decompressorHandle),
		0,
		0)

	return ret != 0
}

func Decompress(
	decompressorHandle HWND,
	compressedData *byte,
	compressedDataSize uint32,
	uncompressedBuffer *byte,
	uncompressedBufferSize uint32,
	uncompressedDataSize *uint32,
) (bool, error) {
	ret, _, errno := syscall.Syscall6(decompress.Addr(), 6,
		uintptr(decompressorHandle),
		uintptr(unsafe.Pointer(compressedData)),
		uintptr(compressedDataSize),
		uintptr(unsafe.Pointer(uncompressedBuffer)),
		uintptr(uncompressedBufferSize),
		uintptr(unsafe.Pointer(uncompressedDataSize)),
	)

	return ret != 0, errno
}
