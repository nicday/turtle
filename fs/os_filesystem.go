package fs

import (
	"io/ioutil"
	"os"
)

// OS implements FileSystem using by wrapping the os package.
type OS struct{}

// Open is a wrapper around os.Open
func (OS) Open(name string) (File, error) { return os.Open(name) }

// Stat is a wrapper around os.Stat
func (OS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

// ReadFile is a wrapper around ioutil.ReadFile
func (OS) ReadFile(name string) ([]byte, error) { return ioutil.ReadFile(name) }
