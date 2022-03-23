package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/rmatsuoka/apifs"
	"github.com/rmatsuoka/ya9p"
)

func main() {
	root := apifs.NewDir()
	name := apifs.NewVal[string]("glenda", func(p []byte) (string, error) {
		return string(p), nil
	})
	root.Mknod("name", name)

	hello := apifs.NewEvent(func() (io.Reader, error) {
		return strings.NewReader(fmt.Sprintf("Hello, %s!\n%v\n", name.Get(), time.Now())), nil
	})
	root.Mknod("hello", hello)

	dir1 := apifs.NewDir()
	dir1.Mknod("name", name)
	root.Mknod("dir1", dir1)

	fsys := apifs.NewFS(root)
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
		}
		go ya9p.ServeFS(conn, fsys)
	}
}
