// MySQL MCP Server - MySQL database access tools
//
// Usage:
//
//	mysql_mcp_server -dsn "username:password@tcp(host:port)/database_name"
//
// Supported tools:
//   - mysql_query: Execute MySQL queries (read-only, SELECT statements)
//   - mysql_execute: Execute MySQL update operations (non-query statements)
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

// dsn defines MySQL database connection string
// Format: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
// Default value: root:password@tcp(127.0.0.1:3306)/test
var dsn = flag.String("dsn", "root:password@tcp(127.0.0.1:3306)/test", "MySQL connection string")

// Database connection
var db *sql.DB

type mysqlQueryReq struct {
	SQL string `json:"sql" description:"SQL query statement to execute"`
}

type mysqlExecuteReq struct {
	SQL string `json:"sql" description:"SQL statement to execute"`
}

func main() {
	flag.Parse()

	// Initialize database
	if err := initDB(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	// Create MCP server
	srv, err := server.NewServer(
		transport.NewStdioServerTransport(),
		server.WithServerInfo(protocol.Implementation{
			Name:    "mysql-mcp-server",
			Version: "1.0.0",
		}),
	)
	if err != nil {
		log.Fatalf("Server creation failed: %v", err)
	}

	queryTool, err := protocol.NewTool(
		"mysql_query",
		"Execute MySQL queries (read-only, SELECT statements)",
		mysqlQueryReq{})
	if err != nil {
		log.Fatalf("Failed to create query tool: %v", err)
		return
	}
	// Register query tool
	srv.RegisterTool(queryTool, handleQuery)

	executeTool, err := protocol.NewTool(
		"mysql_execute",
		"Execute MySQL update operations (INSERT/UPDATE/DELETE and other non-query statements)",
		mysqlExecuteReq{})
	if err != nil {
		log.Fatalf("Failed to create query tool: %v", err)
		return
	}
	// Register execute tool
	srv.RegisterTool(executeTool, handleExecute)

	// Start server
	log.Println("Starting MySQL MCP Server with stdio transport mode")
	if err = srv.Run(); err != nil {
		log.Fatalf("Service runtime error: %v", err)
	}
}

// Initialize database connection
func initDB() error {
	var err error
	db, err = sql.Open("mysql", *dsn)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(60) // 1 minute

	return db.Ping()
}

// Handle MySQL query requests
func handleQuery(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	req := new(mysqlQueryReq)
	if err := json.Unmarshal(request.RawArguments, &req); err != nil {
		return nil, err
	}

	// Execute query
	rows, err := db.Query(req.SQL)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %v", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column names: %v", err)
	}

	// Process results
	var results []map[string]interface{}
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		if err = rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("failed to read row data: %v", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through results: %v", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, fmt.Errorf("JSON serialization failed: %v", err)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Type: "text",
				Text: string(jsonData),
			},
		},
	}, nil
}

// Handle MySQL execute requests
func handleExecute(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	req := new(mysqlExecuteReq)
	if err := json.Unmarshal(request.RawArguments, &req); err != nil {
		return nil, err
	}

	// Execute SQL
	result, err := db.Exec(req.SQL)
	if err != nil {
		return nil, fmt.Errorf("SQL execution error: %v", err)
	}

	// Get results
	lastInsertID, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()

	response := map[string]interface{}{
		"lastInsertId": lastInsertID,
		"rowsAffected": rowsAffected,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("JSON serialization failed: %v", err)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Type: "text",
				Text: string(jsonData),
			},
		},
	}, nil
}
