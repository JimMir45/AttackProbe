# AI Brain

> AI 驱动的知识管理工具集 | AI-Powered Knowledge Management Tools

## Overview | 概述

AI Brain 是一套轻量级的知识管理命令行工具，帮助团队构建可搜索、可问答、可追溯的知识库。

## Tools | 工具

### brain-search

全文搜索工具，搜索知识库中的 Markdown 文档。

```bash
cd brain-search && go build -o brain-search .
./brain-search "关键词"
./brain-search -max 5 "多个 关键词"
```

**特性:**
- 中英文关键词搜索
- 多关键词 AND 逻辑
- 上下文显示
- 按相关度排序

### brain-ask

基于 RAG 的智能问答工具。

```bash
cd brain-ask && go build -o brain-ask .
./brain-ask index              # 构建索引
./brain-ask "你的问题"          # 提问
```

**特性:**
- 向量嵌入 (Ollama)
- 余弦相似度检索
- 来源引用
- 本地 JSON 索引

### brain-archive

对话自动归档工具，将对话保存为结构化 Markdown。

```bash
cd brain-archive && go build -o brain-archive .
echo "对话内容" | ./brain-archive -project myproject -topic "讨论主题"
```

**特性:**
- YAML Frontmatter 元数据
- 时间戳命名
- 参与者/决策记录

### brain-decision

ADR (Architecture Decision Records) 管理工具。

```bash
cd brain-decision && go build -o brain-decision .
./brain-decision new "决策标题"
./brain-decision list
```

**特性:**
- 标准 ADR 格式
- 状态追踪 (proposed → approved → deprecated)
- 按阶段组织

## Requirements | 依赖

- Go 1.22+
- Ollama (用于 brain-ask 的嵌入和问答)

## License

MIT License - see [LICENSE](../../LICENSE)
