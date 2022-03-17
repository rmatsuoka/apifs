package apifs

import (
	"errors"
	"io"
	"io/fs"
	"os"
	pathpkg "path"
	"strings"
)

var (
	ErrIsDir   = errors.New("is a directory")
	ErrNotDir  = errors.New("not a directory")
	ErrNoWrite = errors.New("write prohibited")
)

type Node interface {
	Open(int) (SemiFile, error)
	IsDir() bool
	Walk([]string) (Node, error)
}

type SemiFile interface {
	io.Reader
	io.Writer
	io.Closer
	ReadDir(int) ([]fs.DirEntry, error)
}

type FS Dir

func (f FS) Open(name string) (fs.File, error) {
	return f.OpenFile(name, os.O_RDONLY, 0)
}

func (f FS) OpenFile(name string, mode int, perm fs.FileMode) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	if name == "." {
		s, _ := Dir(f).Open(mode)
		return &file{SemiFile: s, name: ".", isDir: true}, nil
	}

	path := strings.Split(name, "/")
	n, err := Dir(f).Walk(path)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}

	var s SemiFile
	s, err = n.Open(mode)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}
	return &file{SemiFile: s, name: pathpkg.Base(name), isDir: n.IsDir()}, nil
}

type file struct {
	SemiFile
	name  string
	isDir bool
}

func (f *file) Stat() (fs.FileInfo, error) {
	return &info{name: f.name, isDir: f.isDir}, nil
}
