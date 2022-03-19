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
	rwc, err := e.f()
	if err != nil {
		return nil, err
	}
	return &eventFile{e: e, rwc: rwc, name: name}, nil
}

type eventFile struct {
	e    *Event
	rwc  io.Reader
	name string
}

func (f *eventFile) Read(p []byte) (int, error) { return f.rwc.Read(p) }
func (f *eventFile) Stat() (fs.FileInfo, error) {
	return &info{name: f.name}, nil
}

func (f *eventFile) Write(p []byte) (int, error) {
	if w, ok := f.rwc.(io.Writer); ok {
		return w.Write(p)
	} else {
		return 0, ErrNoWrite
	}
}

func (f *eventFile) Close() error {
	if c, ok := f.rwc.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
