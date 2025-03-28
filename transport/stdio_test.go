package transport

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestStdioTransport(t *testing.T) {
	var (
		err    error
		server *stdioServerTransport
		client *stdioClientTransport
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

	clientT, err := NewStdioClientTransport(mockServerTrPath, []string{})
	if err != nil {
		t.Fatalf("NewStdioClientTransport failed: %v", err)
	}

	client = clientT.(*stdioClientTransport)
	server = NewStdioServerTransport().(*stdioServerTransport)

	// Create pipes for communication
	reader1, writer1 := io.Pipe()
	reader2, writer2 := io.Pipe()

	// Set up the communication channels
	server.reader = reader2
	server.writer = writer1
	client.reader = reader1
	client.writer = writer2

	testTransport(t, client, server)
}

func compileMockStdioServerTr(outputPath string) error {
	cmd := exec.Command("go", "build", "-o", outputPath, "../testdata/mock_block_server.go")

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("compilation failed: %v\nOutput: %s", err, output)
	}

	return nil
}
