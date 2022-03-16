package apifs

import (
	"io/fs"
	"time"
)

type info struct {
	name  string
	isDir bool
}

func (i *info) Name() string               { return i.name }
func (i *info) Size() int64                { return 0 }
func (i *info) ModTime() time.Time         { return time.Unix(0, 0) }
func (i *info) IsDir() bool                { return i.isDir }
func (i *info) Info() (fs.FileInfo, error) { return i, nil }
func (i *info) Sys() any                   { return nil }

func (i *info) Type() fs.FileMode {
	if i.isDir {
		return fs.ModeDir
	}
	return 0
}

func (i *info) Mode() fs.FileMode {
	if i.isDir {
		return fs.ModeDir | 0444
	}
	return 0666
}
