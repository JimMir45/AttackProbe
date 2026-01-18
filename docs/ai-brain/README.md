# AI-Brain 模块文档

> 安全团队知识管理工具集

## 概述

AI-Brain 是一套为安全团队设计的知识管理工具，帮助团队沉淀知识、快速检索、辅助决策。

## 工具集

| 工具 | 功能 | 使用场景 |
|------|------|---------|
| **brain-search** | 全文搜索 | 快速查找知识库文档 |
| **brain-ask** | RAG 问答 | 基于知识库的智能问答 |
| **brain-archive** | 对话归档 | 结构化保存重要讨论 |
| **brain-decision** | ADR 管理 | 架构决策记录管理 |

## 快速开始

### brain-search - 全文搜索

```bash
cd packages/ai-brain/brain-search
go build -o brain-search .
./brain-search "关键词"
```

### brain-ask - RAG 问答

```bash
cd packages/ai-brain/brain-ask
go build -o brain-ask .
./brain-ask "你的问题"
```

需要先配置 Ollama 用于向量嵌入。

### brain-archive - 对话归档

```bash
cd packages/ai-brain/brain-archive
go build -o brain-archive .
./brain-archive
```

### brain-decision - ADR 管理

```bash
cd packages/ai-brain/brain-decision
go build -o brain-decision .
./brain-decision list
./brain-decision create "决策标题"
```

## 技术栈

- **语言**: Go 1.22+
- **向量嵌入**: Ollama (用于 brain-ask)
- **存储**: 本地 Markdown 文件 + JSON 索引

## 相关链接

- [源代码](../../packages/ai-brain/)
- [SecOpsHub 整体架构](../secops-hub/01-设计/整体架构设计.md)
