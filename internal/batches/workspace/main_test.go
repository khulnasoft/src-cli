package workspace

import (
	"os"
	"testing"

	"github.com/khulnasoft/src-cli/internal/exec/expect"
)

func TestMain(m *testing.M) {
	code := expect.Handle(m)
	os.Exit(code)
}
