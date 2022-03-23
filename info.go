package apifs

import (
	"io/fs"
	"time"
)

type info struct {
	name string
	mode fs.FileMode
}

func (i *info) Name() string       { return i.name }
func (i *info) Size() int64        { return 0 }
func (i *info) ModTime() time.Time { return time.Unix(0, 0) }
func (i *info) IsDir() bool        { return i.mode.IsDir() }
func (i *info) Mode() fs.FileMode  { return i.mode }
func (i *info) Sys() any           { return nil }

type dirEntry struct {
	name string
	n    Node
}

func (d *dirEntry) Name() string { return d.name }

func (d *dirEntry) IsDir() bool {
	_, ok := d.n.(DirNode)
	return ok
}

func (d *dirEntry) Type() fs.FileMode {
	_, ok := d.n.(DirNode)
	if ok {
		return fs.ModeDir
	}
	return 0
}

func (d *dirEntry) Info() (fs.FileInfo, error) {
	return stat(d.n, d.name)
}
