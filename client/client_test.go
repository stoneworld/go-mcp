package client

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"go-mcp/pkg"
	"go-mcp/protocol"
	"go-mcp/transport"

	"github.com/bytedance/sonic"
)

func TestClient(t *testing.T) {
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

	client, err := NewClient(transport.NewMockClientTransport(in, out))
	if err != nil {
		t.Errorf("NewServer: %+v", err)
	}

	tests := []struct {
		name             string
		f                func(client *Client, request protocol.ClientRequest) (protocol.ServerResponse, error)
		request          protocol.ClientRequest
		expectedResponse protocol.ServerResponse
	}{
		{
			name: "test_call_tool",
			f: func(client *Client, request protocol.ClientRequest) (protocol.ServerResponse, error) {
				return client.CallTool(context.Background(), request.(*protocol.CallToolRequest))
			},
			request: protocol.NewCallToolRequest("test_call_tool", map[string]interface{}{}),
			expectedResponse: protocol.NewCallToolResponse([]protocol.Content{protocol.TextContent{
				Type: "text",
				Text: "test_call_tool",
			}}, true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				reqBytes, err := out.ReadBytes('\n')
				if err != nil {
					t.Errorf("out read: %+v", err)
				}
				jsonrpcReq := &protocol.JSONRPCRequest{}
				if err := pkg.JsonUnmarshal(reqBytes, &jsonrpcReq); err != nil {
					t.Errorf("Json Unmarshal: %+v", err)
				}

				request := make(map[string]interface{})
				if err := pkg.JsonUnmarshal(jsonrpcReq.RawParams, &request); err != nil {
					t.Errorf("Json Unmarshal: %+v", err)
				}

				expectedReqBytes, err := json.Marshal(tt.request)
				if err != nil {
					t.Errorf("json Marshal: %+v", err)
				}
				var expectedReqMap map[string]interface{}
				if err := pkg.JsonUnmarshal(expectedReqBytes, &expectedReqMap); err != nil {
					t.Errorf("json Unmarshal: %+v", err)
				}

				if !reflect.DeepEqual(request, expectedReqMap) {
					t.Errorf("response not as expected.\ngot  = %v\nwant = %v", request, expectedReqMap)
				}

				respBytes, err := sonic.Marshal(protocol.NewJSONRPCSuccessResponse(jsonrpcReq.ID, tt.expectedResponse))
				if err != nil {
					t.Errorf("Json Marshal: %+v", err)
				}
				if _, err := in.Write(append(respBytes, "\n"...)); err != nil {
					t.Errorf("in Write: %+v", err)
				}
			}()

			response, err := tt.f(client, tt.request)
			if err != nil {
				t.Errorf("func exectue: %+v", err)
			}

			if !reflect.DeepEqual(response, tt.expectedResponse) {
				t.Errorf("response not as expected.\ngot  = %v\nwant = %v", response, tt.expectedResponse)
			}
		})
	}
}
