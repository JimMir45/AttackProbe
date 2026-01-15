# LLM Security

> 大语言模型安全测试工具 | LLM Security Testing Tool

## Overview | 概述

LLM Security 是一个自动化的大语言模型安全测试平台，帮助你验证 LLM 应用对抗各种攻击的防护能力。

## Features | 特性

- **58 个内置攻击用例**
  - 提示词注入 (20)
  - 越狱攻击 (18)
  - 敏感信息泄露 (15)
  - 其他攻击 (5)
- **多目标支持**: OpenAI 兼容 API (Ollama, vLLM, etc.)
- **Web UI**: 可视化管理界面
- **自动判定**: 基于关键词的结果判定引擎
- **并发执行**: Worker Pool 并发任务执行

## Quick Start | 快速开始

```bash
# 构建
go build -o llm-security ./cmd/server/

# 运行
./llm-security -port 8080

# 访问 http://localhost:8080
```

## Architecture | 架构

```
┌─────────────────────────────────────┐
│      Vue3 + Element Plus Frontend   │
├─────────────────────────────────────┤
│      Gin HTTP API Gateway           │
├─────────────────────────────────────┤
│  Business Logic (Service Layer)     │
│  • Target Mgmt  • Test Cases        │
│  • Task Execution  • Report Gen     │
├─────────────────────────────────────┤
│  Data Access (GORM)                 │
├─────────────────────────────────────┤
│      SQLite Database                │
└─────────────────────────────────────┘
```

## API Endpoints | API 接口

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/target/add | 添加目标 |
| POST | /api/v1/target/test | 测试连通性 |
| POST | /api/v1/task/add | 创建任务 |
| POST | /api/v1/task/start | 启动任务 |
| POST | /api/v1/task/results | 获取结果 |

## License

MIT License - see [LICENSE](../../LICENSE)
