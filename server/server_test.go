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
	reader1, writer1 := io.Pipe()
	reader2, writer2 := io.Pipe()

	var (
		in io.ReadWriter = struct {
			io.Reader
			io.Writer
		}{
			Reader: reader1,
			Writer: writer1,
		}

		out io.ReadWriter = struct {
			io.Reader
			io.Writer
		}{
			Reader: reader2,
			Writer: writer2,
		}

		outScan = bufio.NewScanner(out)
	)

	server, err := NewServer(
		transport.NewMockServerTransport(in, out),
		WithInfo(protocol.Implementation{
			Name:    "ExampleServer",
			Version: "1.0.0",
		}))

	// TODO: add mock session id auto
	server.sessionID2session.Store("mock", &session{})
	if err != nil {
		t.Fatalf("NewServer: %+v", err)
	}
	testTool := &protocol.Tool{
		Name:        "test_tool",
		Description: "test_tool",
		InputSchema: map[string]interface{}{
			"a": "int",
		},
	}
	server.AddTool(testTool)

	testPrompt := protocol.Prompt{
		Name:        "test_prompt",
		Description: "test_prompt_description",
		Arguments: []protocol.PromptArgument{
			{
				Name:        "params1",
				Description: "params1's description",
				Required:    true,
			},
		},
	}

	testPromtGetResponse := &protocol.GetPromptResult{
		Description: "test_prompt_description",
	}
	server.AddPrompt(testPrompt, func(protocol.GetPromptRequest) *protocol.GetPromptResult {
		return testPromtGetResponse
	})

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
		{
			name:    "test_initialize",
			method:  protocol.Initialize,
			request: protocol.InitializeRequest{},
			expectedResponse: protocol.InitializeResult{
				ProtocolVersion: protocol.PROTOCOL_VERSION,
				Capabilities:    server.capabilities,
				ServerInfo:      server.serverInfo,
			},
		},
		{
			name:             "test_ping",
			method:           protocol.Ping,
			request:          protocol.PingRequest{},
			expectedResponse: protocol.PingResult{},
		},
		{
			name:    "test_list_prompt",
			method:  protocol.PromptsList,
			request: protocol.ListPromptsRequest{},
			expectedResponse: protocol.ListPromptsResult{
				Prompts: []protocol.Prompt{testPrompt},
			},
		},
		{
			name:   "test_get_prompt",
			method: protocol.PromptsGet,
			request: protocol.GetPromptRequest{
				Name: testPrompt.Name,
			},
			expectedResponse: testPromtGetResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuid, _ := uuid.NewUUID()
			req := protocol.NewJSONRPCRequest(uuid, tt.method, tt.request)
			reqBytes, err := sonic.Marshal(req)
			if err != nil {
				t.Fatalf("json Marshal: %+v", err)
			}
			if _, err := in.Write(append(reqBytes, "\n"...)); err != nil {
				t.Fatalf("in Write: %+v", err)
			}

			var respBytes []byte
			if outScan.Scan() {
				respBytes = outScan.Bytes()
				if outScan.Err() != nil {
					t.Fatalf("outScan: %+v", err)
				}
			}

			var respMap map[string]interface{}
			if err := pkg.JsonUnmarshal(respBytes, &respMap); err != nil {
				t.Error(err)
			}

			expectedResp := protocol.NewJSONRPCSuccessResponse(uuid, tt.expectedResponse)
			expectedRespBytes, err := json.Marshal(expectedResp)
			if err != nil {
				t.Fatalf("json Marshal: %+v", err)
			}
			var expectedRespMap map[string]interface{}
			if err := pkg.JsonUnmarshal(expectedRespBytes, &expectedRespMap); err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(respMap, expectedRespMap) {
				t.Fatalf("response not as expected.\ngot  = %v\nwant = %v", respMap, expectedRespMap)
			}
		})
	}
}
