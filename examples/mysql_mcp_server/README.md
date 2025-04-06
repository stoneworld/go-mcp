# MySQL MCP Server Example

This is a MySQL tool server implemented based on the Go-MCP framework, providing the ability to interact with MySQL databases through the MCP protocol.

## Features

The server provides the following tools:

1. `mysql_query` - Execute MySQL queries (read-only SELECT statements)
2. `mysql_execute` - Execute MySQL update operations (non-query statements like INSERT/UPDATE/DELETE)

## Installation and Running

### Installation Steps

Use the `go install` command to install to the `$GOPATH/bin` directory:

```bash
# Clone the repository from GitHub
git clone https://github.com/ThinkInAIXYZ/go-mcp.git
cd go-mcp/examples/mysql_mcp_server

# Install MySQL MCP server
go install
```

After installation, you can run it directly:

```bash
# Use default configuration
mysql_mcp_server

# Specify MySQL connection
mysql_mcp_server -dsn "username:password@tcp(host:port)/database_name"
```

## Configuration and Usage in Deepchat
### Configuration Steps
1. Go to Settings->MCP Settings->Add Server
2. Enter the following configuration
```json
{
    "mcpServers": {
      "go_mysql_mcp": {
        "command": "mysql_mcp_server",
        "args": [
            "-dsn",
            "root:password@tcp(127.0.0.1:3306)/test"
        ]
      }
    }
}
```

## Configuration and Usage in Cursor

### Configuration Steps

1. Edit Cursor's MCP configuration file (usually located at `~/.cursor/mcp.json`)
2. Add the following configuration:

```json
"go_mysql_mcp": {
  "command": "mysql_mcp_server",
  "args": [
    "-dsn",
    "root:password@tcp(127.0.0.1:3306)/test"
  ]
}
```

### Using in Cursor

Use directly in Cursor dialog:

```
/mysql_query SELECT * FROM users LIMIT 10
```

```
/mysql_execute INSERT INTO users (name, email) VALUES ('John Doe', 'johndoe@example.com')
```