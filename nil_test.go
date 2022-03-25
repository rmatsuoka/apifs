// Check that calling methods with nil does not panic
package apifs_test

import (
	"testing"

	"github.com/rmatsuoka/apifs"
)

func TestNilFS(t *testing.T) {
	var fsys *apifs.FS = nil
	fsys.Open("foo")

	fsys = apifs.NewFS(nil)
	fsys.Open("foo")
}
