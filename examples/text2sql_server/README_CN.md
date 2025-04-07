# Text2SQL 服务器

<div align="center">
  <a href="./README.md">English</a> / 中文
</div>

这是一个基于 go-mcp 框架的自然语言转 SQL 查询服务，支持将自然语言问题转换为可执行的 SQL 语句并返回查询结果。

## 功能特点

- 自然语言转 SQL：自动将用户问题转换为相应的 SQL 查询语句
- 查询执行：自动执行生成的 SQL 语句并返回查询结果
- 多次尝试：支持多次尝试以获取最佳查询结果

## 环境要求

- Go 1.23+
- 数据库连接（支持标准 SQL 数据库）
- OpenAI API 密钥（或者兼容openAI协议的供应商的密钥）

## 配置说明

### 环境变量

- `link`：数据库连接字符串（必需）
- OpenAI 配置在代码中设置：
  - API 密钥
  - 基础 URL
  - 模型 名字

### 配置参数

```go
text2sql.Config{
    DbLink:    link,     // 数据库连接字符串
    ShouldRun: true,     // 是否执行 SQL
    Times:     3,        // 获取 SQL 的数量
    Try:       3         // 失败尝试次数
}
```

## 使用方法

1. 设置环境变量：
```bash
export link="username:password@tcp(host:port)/database_name"
export OPENAI_API_KEY="sk-******"
export OPENAI_BASE_URL="https://api.openai.com/v1"
export OPENAI_MODEL_NAME="gpt-4o-mini"
```

2. 启动服务：
```bash
go run .
```


## 响应格式

服务将返回生成的 SQL 语句和执行结果：
```
sql: <生成的 SQL 语句>, result: <查询结果>
```

## 注意事项

1. 确保数据库连接字符串配置正确
2. 配置有效的 OpenAI API 密钥
3. 根据实际需求调整重试和尝试次数
4. 建议在生产环境中使用环境变量管理敏感配置信息