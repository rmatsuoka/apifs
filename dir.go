package apifs

import (
	"fmt"
	"io"
	"io/fs"
	"sync"
)

type Dir struct {
	m *sync.Map
}

func NewDir() *Dir {
	return &Dir{new(sync.Map)}
}

func (d *Dir) Open(name string, mode int) (fs.File, error) {
	return newdirFile(d, name), nil
}

func (d *Dir) Child(name string) (Node, error) {
	n, ok := d.m.Load(name)
	if !ok {
		return nil, fs.ErrNotExist
	}
	return n.(Node), nil
}

func (d *Dir) Add(name string, n Node) error {
	_, ok := d.m.LoadOrStore(name, n)
	if ok {
		return fmt.Errorf("(*Dir).Add: %s has already existed", name)
	}
	return nil
}

type ent struct {
	name string
	n    Node
}

type dirFile struct {
	ents   []*ent
	offset int
	name   string
}

func newdirFile(d *Dir, name string) *dirFile {
	f := &dirFile{name: name}
	d.m.Range(func(k, v any) bool {
		name := k.(string)
		n := v.(Node)
		f.ents = append(f.ents, &ent{name: name, n: n})
		return true
	})
	return f
}

func (f *dirFile) Read(p []byte) (int, error) { return 0, ErrIsDir }
func (f *dirFile) Close() error               { return nil }
func (f *dirFile) Stat() (fs.FileInfo, error) {
	return &info{name: f.name, isDir: true}, nil
}

func (f *dirFile) ReadDir(n int) ([]fs.DirEntry, error) {
	l := len(f.ents) - f.offset
	if n > 0 && n < l {
		l = n
	}
	if l == 0 {
		if n <= 0 {
			return nil, nil
		}
		return nil, io.EOF
	}

	e := make([]fs.DirEntry, 0, l)
	for i := f.offset; i < l+f.offset; i++ {
		_, ok := f.ents[i].n.(DirNode)
		e = append(e, &info{
			name:  f.ents[i].name,
			isDir: ok,
		})
	}
	f.offset += l
	return e, nil
}
