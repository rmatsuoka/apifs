package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rmatsuoka/apifs"
)

func main() {
	name := apifs.NewVar[string]("name", "glenda", func(p []byte) (string, error) {
		return string(p), nil
	})
	hello := apifs.NewEvent("hello", func() (io.Reader, error) {
		return strings.NewReader(fmt.Sprintf("Hello, %s!\n%v\n", name.Get(), time.Now())), nil
	})
	fs := apifs.NewFS(name, hello)

	f, _ := fs.Open("hello")
	io.Copy(os.Stdout, f)
	f.Close()

	f, _ = fs.Open("name")
	fmt.Fprint(f.(io.Writer), "rmatsuoka")
	f.Close()

	f, _ = fs.Open("hello")
	io.Copy(os.Stdout, f)
	f.Close()
}