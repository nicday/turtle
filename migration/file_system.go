package migration

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

// FS is the active FileSystem
var FS FileSystem = osFS{}

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

// osFS implements FileSystem using by wrapping the os package.
type osFS struct{}

func (osFS) Open(name string) (File, error)        { return os.Open(name) }
func (osFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }
func (osFS) ReadFile(name string) ([]byte, error)  { return ioutil.ReadFile(name) }

// MockFS is a concrete implementation of the File interface.
type MockFS struct {
	Files map[string]MockFile
}

// NewMockFS initializes a new MockFS.
func NewMockFS() *MockFS {
	return &MockFS{
		Files: map[string]MockFile{},
	}
}

// AddFiles adds a collection of File to the MockFS.
func (m *MockFS) AddFiles(path string, files ...MockFile) error {
	for _, file := range files {
		info, err := file.Stat()
		if err != nil {
			return err
		}

		// If the current file is a directory
		if len(file.Files) > 0 {
			// Add the directory
			m.Files[fmt.Sprintf("%s%s", path, info.Name())] = file

			// Then add each file within the directory
			info, err := file.Stat()
			if err != nil {
				return err
			}
			m.AddFiles(fmt.Sprintf("%s%s/", path, info.Name()), file.Files...)
			return nil
		}

		m.Files[fmt.Sprintf("%s%s", path, info.Name())] = file
	}

	return nil
}

// Open satisfies the File interface. It will return a File or raise an error.
func (m MockFS) Open(name string) (File, error) {
	file, ok := m.Files[name]
	if !ok {
		return &MockFile{}, errors.New("file not found")
	}

	return file, nil
}

// Stat satisfies the File interface. It will return the file's FileInfo or raise and error.
func (m MockFS) Stat(name string) (os.FileInfo, error) {
	file, ok := m.Files[name]
	if !ok {
		return MockFileInfo{}, errors.New("file not found")
	}

	return file.Stat()
}

// ReadFile satisfies the File interface. It will return a File contents or raise an error.
func (m MockFS) ReadFile(name string) ([]byte, error) {
	file, ok := m.Files[name]
	if !ok {
		return []byte{}, errors.New("file not found")
	}

	return file.Content(), nil
}

// MockFile is a concrete implementation of the File interface
type MockFile struct {
	Info    os.FileInfo
	Files   []MockFile
	content []byte

	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
}

// NewMockFile initializes a new MockFile.
func NewMockFile(name string, content []byte, files ...MockFile) MockFile {
	// If any files have been passed in, this is a directory.
	var isDir bool
	if len(files) != 0 {
		isDir = true
	}

	info := NewMockFileInfo(name, isDir)

	m := MockFile{
		Files:   files,
		Info:    info,
		content: content,
	}

	return m
}

// Stat satisfies the File interface. It will return MockFile.Info
func (m MockFile) Stat() (os.FileInfo, error) {
	return m.Info, nil
}

// Readdir satisfies the File interface. It will return a slice of FileInfo for the contained files.
func (m MockFile) Readdir(n int) ([]os.FileInfo, error) {
	infos := []os.FileInfo{}

	for _, file := range m.Files {
		info, err := file.Stat()
		if err != nil {
			return infos, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

// Content returns the file contents.
func (m MockFile) Content() []byte {
	return m.content
}

// MockFileInfo is a concrete implementation of os.FileInfo.
type MockFileInfo struct {
	name string
	dir  bool
}

// NewMockFileInfo initializes a new MockFileInfo.
func NewMockFileInfo(name string, isDir bool) MockFileInfo {
	return MockFileInfo{
		name: name,
		dir:  isDir,
	}
}

// Name satisfies the os.FileInfo interface. It will return the MockFileInfo.name.
func (m MockFileInfo) Name() string {
	return m.name
}

// IsDir satisfies the os.FileInfo interface. It will return MockFileInfo.dir
func (m MockFileInfo) IsDir() bool {
	return m.dir
}

// Mode satisfies the os.FileInfo interface. It will always return full executable permissions.
func (m MockFileInfo) Mode() os.FileMode {
	return 0777
}

// ModTime satisfies the os.FileInfo interface. It will always return the current time.
func (m MockFileInfo) ModTime() time.Time {
	return time.Now()
}

// Size satisfies the os.FileInfo interface. It will always return 0.
func (m MockFileInfo) Size() int64 {
	return 0
}

// Sys satisfies the os.FileInfo interface. It will always return nil.
func (m MockFileInfo) Sys() interface{} {
	return nil
}
