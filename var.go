package apifs

import (
	"bytes"
	"fmt"
	"io/fs"
	"strings"
	"sync"
)

type Var[T any] struct {
	v    T
	name string
	f    func([]byte) (T, error)
	sync.RWMutex
}

func NewVar[T any](init T, unmarshal func([]byte) (T, error)) *Var[T] {
	return &Var[T]{v: init, f: unmarshal}
}

func (v *Var[T]) Open(int) (SemiFile, error) {
	return &varFile[T]{
		v: v,
		r: strings.NewReader(fmt.Sprint(v.Get())),
		w: new(bytes.Buffer),
	}, nil
}

func (v *Var[T]) IsDir() bool { return false }
func (v *Var[T]) Walk(path []string) (Node, error) {
	if len(path) != 0 {
		return nil, ErrNotDir
	}
	return v, nil
}

func (v *Var[T]) Get() T {
	v.RLock()
	defer v.RUnlock()
	return v.v
}

func (v *Var[T]) Set(n T) {
	v.Lock()
	defer v.Unlock()
	v.v = n
}

type varFile[T any] struct {
	v     *Var[T]
	w     *bytes.Buffer
	r     *strings.Reader
	dirty bool
}

func (f *varFile[T]) Read(p []byte) (int, error)         { return f.r.Read(p) }
func (f *varFile[T]) ReadDir(int) ([]fs.DirEntry, error) { return nil, ErrNotDir }

func (f *varFile[T]) Write(p []byte) (int, error) {
	f.dirty = true
	return f.w.Write(p)
}
func (f *varFile[T]) Close() error {
	if f.dirty {
		n, err := f.v.f(f.w.Bytes())
		if err != nil {
			return err
		}
		f.v.Set(n)
	}
	return nil
}
