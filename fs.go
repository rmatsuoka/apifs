package apifs

import (
	"errors"
	"io/fs"
	"os"
	"strings"
)

var (
	ErrIsDir   = errors.New("is a directory")
	ErrNotDir  = errors.New("not a directory")
	ErrNoWrite = errors.New("write prohibited")
)

type Node interface {
	Open(int) (fs.File, error)
	Name() string
	IsDir() bool
	Children() ([]Node, error)
}

type FS struct {
	root *Dir
}

func NewFS(n ...Node) *FS {
	return &FS{&Dir{DirName: ".", Nodes: n}}
}

func (f *FS) Open(name string) (fs.File, error) {
	return f.OpenFile(name, os.O_RDONLY, 0)
}

func (f *FS) OpenFile(name string, mode int, perm fs.FileMode) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, fs.ErrInvalid
	}
	if name == "." {
		return f.root.Open(mode)
	}
	n, err := namen(f.root, strings.Split(name, "/"))
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}
	return n.Open(mode)
}

func namen(n Node, path []string) (Node, error) {
	if len(path) == 0 {
		return n, nil
	}
	chd, err := n.Children()
	if err != nil {
		return nil, err
	}
	for _, c := range chd {
		if c.Name() == path[0] {
			return namen(c, path[1:])
		}
	}
	return nil, fs.ErrNotExist
}