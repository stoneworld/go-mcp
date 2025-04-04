# Go-MCP



<p align="center">
  <a href="https://pkg.go.dev/github.com/ThinkInAIXYZ/go-mcp"><img src="https://pkg.go.dev/badge/github.com/ThinkInAIXYZ/go-mcp.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/ThinkInAIXYZ/go-mcp"><img src="https://goreportcard.com/badge/github.com/ThinkInAIXYZ/go-mcp" alt="Go Report Card"></a>
  <a href="https://github.com/ThinkInAIXYZ/go-mcp/blob/main/LICENSE"><img src="https://img.shields.io/github/license/ThinkInAIXYZ/go-mcp" alt="License"></a>
  <a href="https://github.com/ThinkInAIXYZ/go-mcp/releases"><img src="https://img.shields.io/github/v/release/ThinkInAIXYZ/go-mcp" alt="Release"></a>
</p>

<p align="center">
  <a href="README_CN.md">中文文档</a>
</p>

## 📖 Overview

Go-MCP is a powerful and easy-to-use Go client library designed for interacting with the Model Context Protocol (MCP). This SDK provides comprehensive API coverage, including core features such as resource management, configuration, monitoring, and automated operations.

MCP (Model Context Protocol) is a standardized protocol for AI model interaction that enables seamless communication between applications and AI models. In today's rapidly evolving AI landscape, standardized communication protocols have become increasingly important, ensuring interoperability between different systems and components, reducing integration costs, and improving development efficiency. Go-MCP brings this capability to Go applications through clean, idiomatic APIs, enabling developers to easily integrate AI functionality into their Go projects.

Go language is known for its excellent performance, concise syntax, and powerful concurrency support, making it particularly suitable for building high-performance network services and system tools. Through Go-MCP, developers can fully leverage these advantages of Go while enjoying the standardization and interoperability benefits brought by the MCP protocol. Whether building edge AI applications, microservice architectures, or enterprise-level systems, Go-MCP provides reliable and efficient solutions.

### Core Features

- **Core Protocol Implementation**: Go-MCP fully supports the MCP specification, ensuring seamless interaction with all compatible MCP services and clients. The SDK implements all core methods and notification mechanisms defined in the protocol, including initialization, tool invocation, resource management, prompt handling, and more.

- **Multiple Transport Methods**: Supports SSE (Server-Sent Events) and stdio transport, adapting to different application scenarios and deployment environments. SSE transport is suitable for web-based applications, providing real-time server push capabilities; while stdio transport is suitable for inter-process communication and command-line tools, making MCP functionality easy to integrate into various systems.

- **Rich Notification System**: Comprehensive event handling mechanism supporting real-time updates and status change notifications. Through registering custom notification handlers, applications can respond in real-time to tool list changes, resource updates, prompt list changes, and other events, achieving dynamic and interactive user experiences.

- **Flexible Architecture**: Easy to extend, supporting custom implementations and customization needs. Go-MCP's modular design allows developers to extend or replace specific components according to their needs while maintaining compatibility with the core protocol.

- **Production Ready**: Thoroughly tested and performance-optimized, suitable for high-demand production environments. The SDK adopts Go language best practices, including concurrency control, error handling, and resource management, ensuring stability and reliability under high load conditions.

- **Comprehensive Documentation and Examples**: Provides comprehensive documentation and rich example code to help developers get started quickly and understand in depth. Whether beginners or experienced Go developers can easily master the SDK usage through documentation and examples.

### Why Choose Go-MCP and Its Future Prospects

Go-MCP SDK has significant advantages in the current technical environment and has broad development prospects. Here are the core advantages:

1. **Local Deployment Advantages**: Go-MCP supports local AI application deployment, providing faster response times, better cost-effectiveness, stronger data control capabilities, and more flexible customization options. Go's static compilation features make deployment extremely simple, without managing complex dependencies.

2. **Edge Computing Support**: Particularly suitable for running AI models on resource-constrained edge devices, supporting real-time processing, low-bandwidth environments, offline operations, and data privacy protection. Go's high performance and low memory footprint make it an ideal choice for edge computing.

3. **Microservice Architecture Adaptation**: Perfectly fits modern microservice and serverless architectures, supporting AI microservice encapsulation, event-driven processing, distributed AI systems, and hybrid cloud deployment. Go's lightweight runtime and concurrency model are particularly suitable for handling large numbers of concurrent requests.

4. **Strong Ecosystem**: Benefits from Go's active community and enterprise support, providing rich libraries and frameworks, as well as excellent development toolchains. As the community grows, the SDK's functionality and performance will continue to improve.

5. **Data Security Protection**: Supports local data processing, reducing data transmission requirements and lowering data leakage risks. Provides secure communication methods that can integrate with encryption and authentication mechanisms, meeting data sovereignty and regulatory compliance requirements.

6. **Cross-Platform Compatibility**: Supports all major operating systems and processor architectures, providing consistent behavior and simple deployment methods. Through static linking and cross-compilation support, ensures a unified experience across different platforms.

## 🚀 Installation

Installing the Go-MCP SDK is very simple, just use Go's standard package management tool `go get` command:

```bash
go get github.com/ThinkInAIXYZ/go-mcp
```

This will download and install the SDK and all its dependencies. Go-MCP requires Go 1.18 or higher to ensure support for the latest language features and standard library.


## 🔍 Quick Start

### Client Implementation

Here's a basic client implementation example showing how to create an MCP client, connect to a server, and perform basic operations:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ThinkInAIXYZ/go-mcp/client"
    "github.com/ThinkInAIXYZ/go-mcp/protocol"
    "github.com/ThinkInAIXYZ/go-mcp/transport"
)

func main() {
	// Create transport client (using SSE in this example)
	transportClient, err := transport.NewSSEClientTransport(context.Background(), "http://127.0.0.1:8080/sse")
	if err != nil {
		log.Fatalf("Failed to create transport client: %v", err)
	}

	// Create MCP client using transport
	mcpClient, err := client.NewClient(transportClient,
		// Optional: Set custom notification handler
		client.WithToolsListChangedNotifyHandler(func(ctx context.Context, notification *protocol.ToolListChangedNotification) error {
			fmt.Printf("Tool list updated: %+v\n", notification)
			return nil
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer mcpClient.Close()

	// Get server capabilities
	capabilities := mcpClient.GetServerCapabilities()
	fmt.Printf("Server capabilities: %+v\n", capabilities)

	// List available tools
	toolsResult, err := mcpClient.ListTools(context.Background())
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	fmt.Printf("Available tools: %+v\n", toolsResult.Tools)

	// Call tool
	callResult, err := mcpClient.CallTool(context.Background(), protocol.NewCallToolRequest("example_tool", map[string]interface{}{
		"param1": "value1",
		"param2": 42,
	}))
	if err != nil {
		log.Fatalf("Failed to call tool: %v", err)
	}

	fmt.Printf("Tool call result: %+v\n", callResult)
}
```

This example shows how to:
1. Create an SSE client transport
2. Initialize an MCP client and configure notification handlers
3. Get server capability information
4. List available tools
5. Call a specific tool and handle results

### Server Implementation

Here's a basic server implementation example showing how to create an MCP server, register tool handlers, and handle client requests:

```go
package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

func main() {
	// Create transport server (using SSE in this example)
	transportServer, err := transport.NewSSEServerTransport("127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Failed to create transport server: %v", err)
	}

	// Create MCP server using transport
	mcpServer, err := server.NewServer(transportServer,
		// Set server implementation information
		server.WithServerInfo(protocol.Implementation{
			Name:    "Example MCP Server",
			Version: "1.0.0",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// Register tool handler
	mcpServer.RegisterTool(&protocol.Tool{
		Name:        "current time",
		Description: "Get current time with timezone, Asia/Shanghai is default",
		InputSchema: protocol.InputSchema{
			Type: protocol.Object,
			Properties: map[string]interface{}{
				"timezone": map[string]string{
					"type":        "string",
					"description": "current time timezone",
				},
			},
			Required: []string{"timezone"},
		},
	}, func(req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
		timezone, ok := req.Arguments["timezone"].(string)
		if !ok {
			return nil, errors.New("timezone must be a string")
		}

		loc, err := time.LoadLocation(timezone)
		if err != nil {
			return nil, fmt.Errorf("parse timezone with error: %v", err)
		}
		text := fmt.Sprintf(`current time is %s`, time.Now().In(loc))

		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: text,
				},
			},
		}, nil
	})

	if err = mcpServer.Run(); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
		return
	}
}

```

This example shows how to:
1. Create an SSE transport server
2. Initialize an MCP server and configure server information and capabilities
3. Register tool handlers
4. Start HTTP server to handle client connections

## 🏗️ Architecture Design

The Go-MCP SDK adopts a clear layered architecture design, ensuring code modularity, extensibility, and maintainability. Through a deep understanding of this architecture, developers can better utilize all SDK features and even customize and extend according to their needs.

### Three-Layer Architecture

Go-MCP's architecture can be abstracted into three main layers:

![img.png](docs/images/img.png)

1. **Transport Layer**: Handles underlying communication, supporting different transport protocols
2. **Protocol Layer**: Implements all MCP protocol functionality and data structures
3. **User Layer**: Includes server and client implementations, providing user-facing APIs

This layered design decouples layers from each other, allowing independent evolution and replacement while maintaining overall functionality consistency.

### Transport Layer

The transport layer handles underlying communication details, currently supporting two main transport methods:

![img_1.png](docs/images/img_1.png)

- **HTTP SSE/POST**: HTTP-based server-sent events and POST requests, suitable for network communication
- **Stdio**: Standard input/output stream-based communication, suitable for inter-process communication

The transport layer abstracts through a unified interface, so upper-layer code doesn't need to care about specific transport implementation details. This design allows easy addition of new transport methods, such as WebSocket, gRPC, etc., without affecting upper-layer code.

### Protocol Layer

The protocol layer is the core of Go-MCP, containing all MCP protocol-related definitions, including:

- Data structure definitions
- Request construction
- Response parsing
- Notification handling

The protocol layer implements all functionality defined in the MCP specification, including but not limited to:

- Initialization (Initialize)
- Heartbeat detection (Ping)
- Cancellation operations (Cancellation)
- Progress notifications (Progress)
- Root resource management (Roots)
- Sampling control (Sampling)
- Prompt management (Prompts)
- Resource management (Resources)
- Tool invocation (Tools)
- Completion requests (Completion)
- Logging (Logging)
- Pagination handling (Pagination)

The protocol layer is decoupled from the transport layer through the transport interface, allowing protocol implementation to be independent of specific transport methods.

### User Layer

The user layer includes server and client implementations, providing developer-friendly APIs:

- **Server Implementation**: Handles requests from clients, provides resources and tools, sends notifications
- **Client Implementation**: Connects to server, sends requests, handles responses and notifications

The user layer's design philosophy is to provide synchronous request-response patterns, even if the underlying implementation may be asynchronous. This design makes the API more intuitive and easier to use while maintaining efficient asynchronous processing capabilities.

### Message Processing Flow

In Go-MCP, messages can be abstracted into three types:

1. **Request**: Request messages sent from client to server
2. **Response**: Response messages returned from server to client
3. **Notification**: One-way notification messages that can be sent by server or client

Both server and client have sending and receiving capabilities:

- **Sending Capability**: Includes sending messages (requests, responses, notifications) and matching requests with responses
- **Receiving Capability**: Includes routing messages (requests, responses, notifications) and asynchronous/synchronous processing

### Project Structure

Go-MCP's project structure clearly reflects its architectural design:

```
go-mcp/
├── transport/                 # Transport layer implementation
│   ├── sse_client.go          # SSE client implementation
│   ├── sse_server.go          # SSE server implementation
│   ├── stdio_client.go        # Stdio client implementation
│   ├── stdio_server.go        # Stdio server implementation
│   └── transport.go           # Transport interface definition
├── protocol/                  # Protocol layer implementation
│   ├── initialize.go          # Initialization related
│   ├── ping.go                # Heartbeat detection related
│   ├── cancellation.go        # Cancellation operations related
│   ├── progress.go            # Progress notifications related
│   ├── roots.go               # Root resources related
│   ├── sampling.go            # Sampling control related
│   ├── prompts.go             # Prompt management related
│   ├── resources.go           # Resource management related
│   ├── tools.go               # Tool invocation related
│   ├── completion.go          # Completion requests related
│   ├── logging.go             # Logging related
│   ├── pagination.go          # Pagination handling related
│   └── jsonrpc.go             # JSON-RPC related
├── server/                    # Server implementation
│   ├── server.go              # Server core implementation
│   ├── call.go                # Send messages to client
│   ├── handle.go              # Handle messages from client
│   ├── send.go                # Send message implementation
│   └── receive.go             # Receive message implementation
├── client/                    # Client implementation
│   ├── client.go              # Client core implementation
│   ├── call.go                # Send messages to server
│   ├── handle.go              # Handle messages from server
│   ├── send.go                # Send message implementation
│   └── receive.go             # Receive message implementation
└── pkg/                       # Common utility packages
    ├── errors.go              # Error definitions
    └── log.go                 # Log interface definitions
```

This structure makes code organization clear and easy for developers to understand and extend.

### Design Principles

Go-MCP follows these core design principles:

1. **Modularity**: Components are decoupled through clear interfaces, allowing independent development and testing
2. **Extensibility**: Architecture design allows easy addition of new transport methods, protocol features, and user layer APIs
3. **Usability**: Despite potentially complex internal implementation, provides clean, intuitive APIs externally
4. **High Performance**: Leverages Go's concurrency features to ensure efficient message processing
5. **Reliability**: Comprehensive error handling and recovery mechanisms ensure system stability in various scenarios

Through this carefully designed architecture, Go-MCP provides developers with a powerful and flexible tool, enabling them to easily integrate the MCP protocol into their applications, whether simple command-line tools or complex distributed systems.

## 🤝 Contributing

We welcome contributions to the Go-MCP project! Whether reporting issues, suggesting features, or submitting code improvements, your participation will help us make the SDK better.

### Contribution Guidelines

Please check [CONTRIBUTING.md](CONTRIBUTING.md) for detailed contribution process and guidelines. Here are the basic steps for contributing:

1. **Submit Issues**: If you find a bug or have feature suggestions, create an issue on GitHub
2. **Discuss**: Discuss issues or suggestions with maintainers and community to determine solutions
3. **Develop**: Implement your changes, ensuring code meets project style and quality standards
4. **Test**: Add tests to verify your changes and ensure all existing tests pass
5. **Submit PR**: Create a Pull Request describing your changes and issues addressed
6. **Review**: Participate in code review process, making necessary adjustments based on feedback
7. **Merge**: Once your PR is approved, it will be merged into the main branch

### Code Style

We follow standard Go code style and best practices:

- Use `gofmt` or `goimports` to format code
- Follow guidelines in [Effective Go](https://golang.org/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Add documentation comments for all exported functions, types, and variables
- Write clear, concise code, avoiding unnecessary complexity
- Use meaningful variable and function names

## 📄 License

This project is licensed under the [MIT License](LICENSE). The MIT License is a permissive license that allows you to freely use, modify, distribute, and privatize the code, as long as you retain the original license and copyright notice.

## 📞 Contact

For questions, suggestions, or issues, please contact us through:

- **GitHub Issues**: Create an issue on the [project repository](https://github.com/ThinkInAIXYZ/go-mcp/issues)

We welcome any form of feedback and contribution and are committed to building a friendly, inclusive community.

## ✨ Contributors

Thank you for your contribution to Go-MCP!

<a href="https://github.com/ThinkInAIXYZ/go-mcp/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=ThinkInAIXYZ/go-mcp" alt="Contributors" />
</a>
