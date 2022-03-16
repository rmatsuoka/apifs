package apifs

import (
	"io"
	"io/fs"
)

type Dir struct {
	DirName string
	Nodes   []Node
}

func (d *Dir) Open(int) (fs.File, error) { return &dirFile{d: d}, nil }
func (d *Dir) Name() string              { return d.DirName }
func (d *Dir) IsDir() bool               { return true }
func (d *Dir) Children() ([]Node, error) { return d.Nodes, nil }

type dirFile struct {
	d      *Dir
	offset int
}

func (f *dirFile) Stat() (fs.FileInfo, error) {
	return &info{name: f.d.DirName, isDir: true}, nil
}

func (f *dirFile) Read(p []byte) (int, error) {
	return 0, &fs.PathError{"read", f.d.DirName, ErrIsDir}
}

func (f *dirFile) Close() error { return nil }

func (f *dirFile) ReadDir(n int) ([]fs.DirEntry, error) {
	l := len(f.d.Nodes) - f.offset
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
			name:  f.d.Nodes[i].Name(),
			isDir: f.d.Nodes[i].IsDir(),
		})
	}
	f.offset += l
	return e, nil
}
