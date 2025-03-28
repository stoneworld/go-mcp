package transport

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestStdioTransport(t *testing.T) {
	var (
		err    error
		server ServerTransport
		client ClientTransport
	)

	mockServerTrPath := filepath.Join(os.TempDir(), "mockstdio_server_tr")
	if err := compileMockStdioServerTr(mockServerTrPath); err != nil {
		t.Fatalf("Failed to compile mock server: %v", err)
	}

	defer func(name string) {
		if err = os.Remove(name); err != nil {
			t.Fatalf("Failed to remove mock server: %v", err)
		}
	}(mockServerTrPath)

	server = NewStdioServerTransport()
	if client, err = NewStdioClientTransport(mockServerTrPath); err != nil {
		t.Fatalf("NewStdioClientTransport failed: %v", err)
	}

	testTransport(t, client, server)
}

func compileMockStdioServerTr(outputPath string) error {
	cmd := exec.Command("go", "build", "-o", outputPath, "../testdata/mockstdio_server_tr.go")

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("compilation failed: %v\nOutput: %s", err, output)
	}

	return nil
}
