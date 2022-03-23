package apifs

import (
	"errors"
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
	Open(basename string, mode int) (fs.File, error)
}

type DirNode interface {
	Node
	Child(name string) (Node, error)
}

type StatNode interface {
	Node
	Stat(basename string) (fs.FileInfo, error)
}

func stat(n Node, basename string) (fs.FileInfo, error) {
	if n, ok := n.(StatNode); ok {
		return n.Stat(basename)
	}
	f, err := n.Open(basename, os.O_RDONLY)
	if err != nil {
		return nil, err
	}
	return f.Stat()
}

type FS struct {
	root DirNode
}

func NewFS(root DirNode) *FS {
	return &FS{root}
}

func (fsys *FS) Open(name string) (fs.File, error) {
	return fsys.OpenFile(name, os.O_RDONLY, 0)
}

func (fsys *FS) OpenFile(name string, mode int, perm fs.FileMode) (fs.File, error) {
	n, err := fsys.namen(name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}

	var f fs.File
	f, err = n.Open(pathpkg.Base(name), mode)
	if err != nil {
		return f, &fs.PathError{Op: "open", Path: name, Err: err}
	}
	return f, nil
}

func (fsys *FS) Stat(name string) (fs.FileInfo, error) {
	n, err := fsys.namen(name)
	if err != nil {
		return nil, &fs.PathError{Op: "stat", Path: name, Err: err}
	}
	return stat(n, pathpkg.Base(name))
}

func (fsys *FS) namen(name string) (Node, error) {
	if !fs.ValidPath(name) {
		return nil, fs.ErrInvalid
	}
	if name == "." {
		return fsys.root, nil
	}

	path := strings.Split(name, "/")
	var n Node = fsys.root
	for _, p := range path {
		d, ok := n.(DirNode)
		if !ok {
			return nil, fs.ErrNotExist
		}

		var err error
		n, err = d.Child(p)
		if err != nil {
			return nil, err
		}
	}
	return n, nil
}
