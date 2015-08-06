package fs

import (
	"io"
	"os"
)

// FileSystem is a generic interface for a file system.
type FileSystem interface {
	Open(name string) (File, error)
	Stat(name string) (os.FileInfo, error)
	ReadFile(name string) ([]byte, error)
}

// File is a generic interface for a file.
type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
	Readdir(n int) ([]os.FileInfo, error)
}
