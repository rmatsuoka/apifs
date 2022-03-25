package apifs

import (
	"io"
	"io/fs"
)

type Event struct {
	f func() (io.Reader, error)
}

func NewEvent(handler func() (io.Reader, error)) *Event {
	return &Event{f: handler}
}

func (e *Event) Open(name string, mode int) (fs.File, error) {
	var rc io.Reader
	var err error
	if e.f != nil {
		rc, err = e.f()
	}
	return &eventFile{rc: rc, name: name}, err
}

func (e *Event) Stat(name string) (fs.FileInfo, error) {
	return &info{name: name, mode: 0444}, nil
}

type eventFile struct {
	rc   io.Reader
	name string
}

func (f *eventFile) Read(p []byte) (int, error) {
	if f.rc == nil {
		return 0, io.EOF
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
