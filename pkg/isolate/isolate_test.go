package isolate

import (
	"testing"
)

func TestSandboxLifecycle(t *testing.T) {
	sandbox, err := NewSandbox(0, true)
	if err != nil {
		t.Fatalf("NewSandbox failed: %v", err)
	}

	err = sandbox.Init()
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	err = sandbox.Cleanup()
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}
}
