# apifs
[![Go Reference](https://pkg.go.dev/badge/github.com/rmatsuoka/apifs.svg)](https://pkg.go.dev/github.com/rmatsuoka/apifs)

Package apifs is a framework for creating file system style APIs
like the Plan 9 software.

This package defines the Node type which manipulates programs from
the file system.  For example, a Node named Val that holds a value
can refer to and change that value in the program as
 a variable. In addition, the value can also be referenced and
 modified by "read" and "write" to `(*Val).Open()` file.

# Example
Let's create a file system with two files `name` and `hello`.
``` Go
root := apifs.NewDir()

name := apifs.NewVal[string]("glenda", func(p []byte) (string, error) {
	return string(p), nil
})
root.Mknod("name", name)

hello := apifs.NewEvent(func() (io.Reader, error) {
	return strings.NewReader(fmt.Sprintf("Hello, %s!\n", name.Get()), nil
})
root.Mknod("hello", hello)

fsys := apifs.NewFS(root)
```

First, the contents of two file are as follows
``` Bash
$ ls
name hello
$ cat name
glenda$ cat hello
Hello, glenda
```

Next, let's change the content of name to `gopher`.
``` Bash
$ echo -n gopher > name
```

Then, the contents of hello is changed.
``` Bash
$ cat hello
Hello, gopher
```

How does it work? First, `hello` is an Event. When this node is
opened, The contents of file hello is io.Reader returned by this
function f.  Next, `name` is a Val. This node holds a value, which
can be referenced and changed by the Set and Get methods.  Furthermore,
this value can also be referenced and changed by an open file.
