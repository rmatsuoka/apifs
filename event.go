package apifs

import (
	"io"
	"io/fs"
	"sync"
)

type Event struct {
	f func() (io.Reader, error)
}

func NewEvent(handler func() (io.Reader, error)) *Event {
	return &Event{f: handler}
}

func (e *Event) Open(name string, mode int) (fs.File, error) {
	return &eventFile{e: e, name: name}, nil
}

type eventFile struct {
	e    *Event
	o    sync.Once
	rc   io.Reader
	err  error
	name string
}

func (f *eventFile) Read(p []byte) (int, error) {
	f.o.Do(func() { f.rc, f.err = f.e.f() })
	if f.err != nil {
		return 0, f.err
	}
	return f.rc.Read(p)
}
func (f *eventFile) Stat() (fs.FileInfo, error) {
	return &info{name: f.name, mode: 0444}, nil
}

func (f *eventFile) Close() error {
	if f.rc == nil {
		return nil
	}
	if c, ok := f.rc.(io.Closer); ok {
		return c.Close()
	}
	f.rc = nil
	return nil
}
