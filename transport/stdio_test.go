package transport

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestServerReceive struct {
	receiveMsg []byte
	wg         sync.WaitGroup
}

func NewTestServerReceive() *TestServerReceive {
	return &TestServerReceive{}
}

func (r *TestServerReceive) Receive(ctx context.Context, sessionID string, msg []byte) error {
	r.receiveMsg = msg

	r.wg.Done()

	return nil
}

func TestStdioServerTransport_Run(t *testing.T) {
	t.Parallel()

	t.Run("receive message", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		tr := NewStdioServerTransportWithIO(in, out)

		testServerReceive := NewTestServerReceive()
		testServerReceive.wg.Add(1)

		tr.SetReceiver(testServerReceive)

		err := tr.Run()
		assert.NoError(t, err)

		testMsg := `{"jsonrpc": "2.0", "method": "hello", "params": {}, "id": 1}` + "\n"
		_, err = in.Write([]byte(testMsg))
		assert.NoError(t, err)

		done := make(chan struct{})
		go func() {
			testServerReceive.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Fatal("receive message timeout")
		}

		assert.Equal(t, testMsg, string(testServerReceive.receiveMsg))

		err = tr.Shutdown(context.Background(), context.Background())
		assert.NoError(t, err)
	})

	t.Run("send message", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		tr := NewStdioServerTransportWithIO(in, out)

		err := tr.Run()
		assert.NoError(t, err)

		testMsg := `{"jsonrpc": "2.0", "method": "hello", "params": {}, "id": 1}`
		err = tr.Send(context.Background(), stdioSessionID, Message(testMsg))
		assert.NoError(t, err)

		buf := make([]byte, len(testMsg))
		_, err = out.Read(buf)
		assert.NoError(t, err)

		assert.Equal(t, testMsg, string(buf))

		err = tr.Shutdown(context.Background(), context.Background())
		assert.NoError(t, err)
	})
}

func compileMockStdioServerTr(outputPath string) error {
	cmd := exec.Command("go", "build", "-o", outputPath, "../testdata/mockstdio_server_tr.go")

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("compilation failed: %v\nOutput: %s", err, output)
	}

	return nil
}

type TestClientReceive struct {
	receiveMsg []byte
	wg         sync.WaitGroup
}

func NewTestClientReceive() *TestClientReceive {
	return &TestClientReceive{}
}

func (r *TestClientReceive) Receive(ctx context.Context, msg []byte) error {
	r.receiveMsg = msg

	r.wg.Done()

	return nil
}

func TestStdioClientTransport_Run(t *testing.T) {
	t.Parallel()

	t.Run("send message", func(t *testing.T) {
		mockServerTrPath := filepath.Join(os.TempDir(), "mockstdio_server_tr")
		if err := compileMockStdioServerTr(mockServerTrPath); err != nil {
			t.Fatalf("Failed to compile mock server: %v", err)
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				t.Fatalf("Failed to remove mock server: %v", err)
			}
		}(mockServerTrPath)

		client, err := NewStdioClientTransport(mockServerTrPath)
		assert.NoError(t, err)

		testClientReceive := NewTestClientReceive()
		testClientReceive.wg.Add(1)

		client.SetReceiver(testClientReceive)

		err = client.Start()

		testMsg := `{"jsonrpc": "2.0", "method": "hello", "params": {}, "id": 1}`
		err = client.Send(context.Background(), Message(testMsg))

		done := make(chan struct{})
		go func() {
			testClientReceive.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatal("receive message timeout")
		}

		assert.Equal(t, testMsg+"\n", string(testClientReceive.receiveMsg))

		err = client.Close(context.Background())
		assert.NoError(t, err)
	})
}
