package apifs

import (
	"io"
	"io/fs"
)

type Event struct {
	name string
	f    func() (io.Reader, error)
}

func NewEvent(name string, handler func() (io.Reader, error)) *Event {
	return &Event{name: name, f: handler}
}

func (e *Event) Open(int) (fs.File, error) {
	rwc, err := e.f()
	if err != nil {
		return nil, err
	}
	return &eventFile{e: e, rwc: rwc}, nil
}

func (e *Event) Name() string              { return e.name }
func (e *Event) IsDir() bool               { return false }
func (e *Event) Children() ([]Node, error) { return nil, ErrNotDir }

type eventFile struct {
	e   *Event
	rwc io.Reader
}

func (f *eventFile) Write(p []byte) (int, error) {
	if w, ok := f.rwc.(io.Writer); ok {
		return w.Write(p)
	} else {
		return 0, ErrNoWrite
	}
}

func (f *eventFile) Read(p []byte) (int, error) {
	return f.rwc.Read(p)
}

func (f *eventFile) Close() error {
	if c, ok := f.rwc.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func (f *eventFile) Stat() (fs.FileInfo, error) {
	return &info{name: f.e.name}, nil
}