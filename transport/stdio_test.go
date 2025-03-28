package transport

import (
	"bytes"
	"context"
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
