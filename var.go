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

func NewVar[T any](name string, init T, unmarshal func([]byte) (T, error)) *Var[T] {
	return &Var[T]{v: init, name: name, f: unmarshal}
}

func (v *Var[T]) Open(int) (fs.File, error) {
	return &varFile[T]{
		v: v,
		r: strings.NewReader(fmt.Sprint(v.Get())),
		w: new(bytes.Buffer),
	}, nil
}

func (v *Var[T]) Name() string              { return v.name }
func (v *Var[T]) IsDir() bool               { return false }
func (v *Var[T]) Children() ([]Node, error) { return nil, ErrNotDir }

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
	v *Var[T]
	w *bytes.Buffer
	r *strings.Reader
}

func (f *varFile[T]) Write(p []byte) (int, error) {
	return f.w.Write(p)
}

func (f *varFile[T]) Read(p []byte) (int, error) {
	return f.r.Read(p)
}

func (f *varFile[T]) Stat() (fs.FileInfo, error) {
	return &info{name: f.v.name}, nil
}

func (f *varFile[T]) Close() error {
	n, err := f.v.f(f.w.Bytes())
	if err != nil {
		return err
	}
	f.v.Set(n)
	return nil
}
