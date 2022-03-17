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

func (e *Event) Open(int) (SemiFile, error) {
	rwc, err := e.f()
	if err != nil {
		return nil, err
	}
	return &eventFile{e: e, rwc: rwc}, nil
}

func (e *Event) IsDir() bool { return false }
func (e *Event) Walk(path []string) (Node, error) {
	if len(path) != 0 {
		return nil, ErrNotDir
	}
	return e, nil
}

type eventFile struct {
	e   *Event
	rwc io.Reader
}

func (f *eventFile) ReadDir(int) ([]fs.DirEntry, error) { return nil, ErrNotDir }
func (f *eventFile) Read(p []byte) (int, error)         { return f.rwc.Read(p) }

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
