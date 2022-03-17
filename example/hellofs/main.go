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
	name := apifs.NewVar[string]("name", "glenda", func(p []byte) (string, error) {
		return string(p), nil
	})
	hello := apifs.NewEvent("hello", func() (io.Reader, error) {
		return strings.NewReader(fmt.Sprintf("Hello, %s!\n%v\n", name.Get(), time.Now())), nil
	})
	fs := apifs.NewFS(name, hello)

	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
		}
		go ya9p.ServeFS(conn, fs)
	}
}
