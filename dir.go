package apifs

import (
	"io"
	"io/fs"
)

type Dir map[string]Node

func (d Dir) Open(int) (SemiFile, error) { return newdirFile(d), nil }
func (d Dir) IsDir() bool                { return true }
func (d Dir) Walk(path []string) (Node, error) {
	if len(path) == 0 {
		return d, nil
	}
	n, ok := d[path[0]]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return n.Walk(path[1:])
}

type ent struct {
	name string
	n    Node
}

type dirFile struct {
	ents   []ent
	offset int
}

func newdirFile(d Dir) *dirFile {
	f := new(dirFile)
	f.ents = make([]ent, 0, len(d))
	for k, v := range d {
		f.ents = append(f.ents, ent{name: k, n: v})
	}
	return f
}
func (f *dirFile) Read(p []byte) (int, error)  { return 0, ErrIsDir }
func (f *dirFile) Write(p []byte) (int, error) { return 0, ErrIsDir }
func (f *dirFile) Close() error                { return nil }

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
		e = append(e, &info{
			name:  f.ents[i].name,
			isDir: f.ents[i].n.IsDir(),
		})
	}
	f.offset += l
	return e, nil
}
