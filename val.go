package apifs

import (
	"bytes"
	"fmt"
	"io/fs"
	"strings"
	"sync"
)

type Val[T any] struct {
	val T
	f   func([]byte) (T, error)
	sync.RWMutex
}

func NewVal[T any](init T, unmarshal func([]byte) (T, error)) *Val[T] {
	return &Val[T]{val: init, f: unmarshal}
}

func (v *Val[T]) Open(name string, mode int) (fs.File, error) {
	return &valFile[T]{
		v:    v,
		r:    strings.NewReader(fmt.Sprint(v.Get())),
		w:    new(bytes.Buffer),
		name: name,
	}, nil
}

func (v *Val[T]) Get() T {
	v.RLock()
	defer v.RUnlock()
	return v.val
}

func (v *Val[T]) Set(n T) {
	v.Lock()
	defer v.Unlock()
	v.val = n
}

type valFile[T any] struct {
	v     *Val[T]
	w     *bytes.Buffer
	r     *strings.Reader
	name  string
	dirty bool
}

func (f *valFile[T]) Read(p []byte) (int, error) {
	return f.r.Read(p)
}

func (f *valFile[T]) Write(p []byte) (int, error) {
	f.dirty = true
	return f.w.Write(p)
}

func (f *valFile[T]) Stat() (fs.FileInfo, error) {
	return &info{name: f.name}, nil
}

func (f *valFile[T]) Close() error {
	if f.dirty {
		n, err := f.v.f(f.w.Bytes())
		if err != nil {
			return err
		}
		f.v.Set(n)
	}
	return nil
}
