# Text2SQL Server

<div align="center">
  <a href="./README_CN.md">中文</a> / English
</div>

A natural language to SQL query service based on the go-mcp framework, supporting the conversion of natural language questions into executable SQL statements and returning query results.

## Features

- Natural Language to SQL: Automatically converts user questions into corresponding SQL query statements
- Query Execution: Automatically executes generated SQL statements and returns query results
- Multiple Attempts: Supports multiple attempts to obtain the best query results
- OpenAI Integration: Uses OpenAI's GPT model for natural language understanding and SQL generation

## Requirements

- Go 1.23+
- Database connection (supports standard SQL databases)
- OpenAI API key (or key from a vendor compatible with the openAI protocol)

## Configuration

### Environment Variables

- `link`: Database connection string (required)
- OpenAI configuration set in code:
  - API Key
  - Base URL
  - Model

### Configuration Parameters

```go
text2sql.Config{
    DbLink:    link,     // Database connection string
    ShouldRun: true,     // Whether to execute SQL
    Times:     3,        // Retry count
    Try:       3         // Attempt count
}
```

## Usage

1. Set environment variables:
```bash
export link="username:password@tcp(host:port)/database_name"
export OPENAI_API_KEY="sk-******"
export OPENAI_BASE_URL="https://api.openai.com/v1"
export OPENAI_MODEL_NAME="gpt-4o-mini"
```

2. Start the service:
```bash
go run .
```

## Response Format

The service will return the generated SQL statement and execution results:
```
sql: <generated SQL statement>, result: <query results>
```

## Notes

1. Ensure the database connection string is correctly configured
2. Configure a valid OpenAI API key
3. Adjust retry and attempt counts according to actual needs
4. Recommended to use environment variables for managing sensitive configuration information in production environments