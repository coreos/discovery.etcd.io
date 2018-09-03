package e2e

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Unsetenv("ETCDCTL_API")
	os.Exit(m.Run())
}
