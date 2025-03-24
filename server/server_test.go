package server

import (
	"bufio"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"go-mcp/pkg"
	"go-mcp/protocol"
	"go-mcp/transport"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

func TestServer(t *testing.T) {
	readerIn, writerIn := io.Pipe()
	readerOut, writerOut := io.Pipe()

	var (
		in = bufio.NewReadWriter(
			bufio.NewReader(readerIn),
			bufio.NewWriter(writerIn),
		)
		out = bufio.NewReadWriter(
			bufio.NewReader(readerOut),
			bufio.NewWriter(writerOut),
		)
	)

	server, err := NewServer(transport.NewMockServerTransport(in, out))
	if err != nil {
		t.Errorf("NewServer: %+v", err)
	}
	testTool := &protocol.Tool{
		Name:        "test_tool",
		Description: "test_tool",
		InputSchema: map[string]interface{}{
			"a": "int",
		},
	}
	server.AddTool(testTool)

	go func() {
		if err := server.Start(); err != nil {
			t.Errorf("server start: %+v", err)
		}
	}()

	tests := []struct {
		name             string
		method           protocol.Method
		request          protocol.ClientRequest
		expectedResponse protocol.ServerResponse
	}{
		{
			name:             "test_list_tool",
			method:           protocol.ToolsList,
			request:          protocol.ListToolsRequest{},
			expectedResponse: protocol.ListToolsResult{Tools: []*protocol.Tool{testTool}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuid, _ := uuid.NewUUID()
			req := protocol.NewJSONRPCRequest(uuid, tt.method, tt.request)
			reqBytes, err := sonic.Marshal(req)
			if err != nil {
				t.Errorf("json Marshal: %+v", err)
			}
			if _, err := in.Write(append(reqBytes, "\n"...)); err != nil {
				t.Errorf("in Write: %+v", err)
			}

			respBytes, err := out.ReadBytes('\n')
			if err != nil {
				t.Errorf("out read: %+v", err)
			}

			var respMap map[string]interface{}
			if err := pkg.JsonUnmarshal(respBytes, &respMap); err != nil {
				t.Errorf("json Unmarshal: %+v", err)
			}

			expectedResp := protocol.NewJSONRPCSuccessResponse(uuid, tt.expectedResponse)
			expectedRespBytes, err := json.Marshal(expectedResp)
			if err != nil {
				t.Errorf("json Marshal: %+v", err)
			}
			var expectedRespMap map[string]interface{}
			if err := pkg.JsonUnmarshal(expectedRespBytes, &expectedRespMap); err != nil {
				t.Errorf("json Unmarshal: %+v", err)
			}

			if !reflect.DeepEqual(respMap, expectedRespMap) {
				t.Errorf("response not as expected.\ngot  = %v\nwant = %v", respMap, expectedRespMap)
			}
		})
	}
}
