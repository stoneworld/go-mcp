package transport

import "testing"

func TestSSETransport(t *testing.T) {
	serverTransport, err := NewSSEServerTransport()
	if err != nil {
		t.Errorf("NewSSEServerTransport: %+v", err)
	}
	clientTransport, err := NewSSEClientTransport()
	if err != nil {
		t.Errorf("NewSSEClientTransport: %+v", err)
	}
	testClient2Server(t, clientTransport, serverTransport)
	testServer2Client(t, clientTransport, serverTransport)
}
