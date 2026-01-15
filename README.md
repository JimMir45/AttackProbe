# AttackProbe

> Open-source Security Validation Platform | 开源安全验证平台

AttackProbe 是一个综合性的开源安全项目，包含三大核心模块：

- **BAS Platform** - 攻击模拟有效性验证平台
- **LLM Security** - 大模型安全测试工具
- **AI Brain** - AI驱动的知识管理工具集

## Modules | 模块

### 1. LLM Security (llm-security)

针对大语言模型的自动化安全测试工具。

**Features | 特性:**
- 58 个内置攻击用例（提示词注入、越狱攻击、敏感信息泄露）
- 支持 OpenAI 兼容 API（Ollama、vLLM 等）
- Web UI 可视化界面
- 自动化结果判定

**Quick Start | 快速开始:**
```bash
cd packages/llm-security
go build -o llm-security ./cmd/server/
./llm-security -port 8080
```

**Tech Stack | 技术栈:** Go + Gin + Vue 3 + SQLite

---

### 2. AI Brain (ai-brain)

AI 驱动的团队知识管理工具集。

**Tools | 工具:**

| Tool | Description |
|------|-------------|
| `brain-search` | 全文搜索 - 搜索知识库中的 Markdown 文档 |
| `brain-ask` | RAG 问答 - 基于向量检索的智能问答 |
| `brain-archive` | 对话归档 - 自动归档对话为结构化文档 |
| `brain-decision` | 决策记录 - ADR 架构决策记录管理 |

**Quick Start | 快速开始:**
```bash
cd packages/ai-brain/brain-search
go build -o brain-search .
./brain-search "关键词"
```

**Tech Stack | 技术栈:** Go + Ollama (Embeddings)

---

### 3. BAS Platform (规划中)

综合性 BAS (Breach and Attack Simulation) 平台，整合传统安全验证与 LLM 安全验证。

**Planned Features | 规划功能:**
- 传统漏洞扫描模拟
- 网络攻击模拟
- ATT&CK 框架映射
- 统一任务调度与报告

---

## Project Structure | 项目结构

```
AttackProbe/
├── packages/
│   ├── llm-security/       # LLM 安全测试模块
│   │   ├── cmd/            # 程序入口
│   │   ├── internal/       # 内部实现
│   │   ├── pkg/            # 公共包
│   │   └── web/            # 前端源码
│   └── ai-brain/           # AI 知识管理工具
│       ├── brain-search/   # 全文搜索
│       ├── brain-ask/      # RAG 问答
│       ├── brain-archive/  # 对话归档
│       └── brain-decision/ # 决策记录
├── docs/                   # 项目文档
│   ├── 00-立项/            # 立项文档
│   ├── 01-设计/            # 设计文档
│   └── ...
└── LICENSE
```

---

## Why AttackProbe? | 为什么选择 AttackProbe

| Feature | Description |
|---------|-------------|
| **中文原生** | 专为中文安全场景和中文 LLM 优化 |
| **轻量部署** | 单文件部署，无外部依赖 |
| **开源开放** | MIT License，欢迎贡献 |
| **模块化** | 各模块可独立使用或组合 |

---

## Getting Started | 快速开始

### Prerequisites | 环境要求

- Go 1.22+
- Node.js 18+ (前端构建)
- Ollama (可选，用于 RAG 和 LLM 测试)

### Build | 构建

```bash
# 构建 LLM Security
cd packages/llm-security
go build -o llm-security ./cmd/server/

# 构建 AI Brain 工具
cd packages/ai-brain/brain-search
go build -o brain-search .
```

---

## Contributing | 贡献

欢迎贡献代码、文档、Issue 和建议！

- **安全研究者**: 贡献攻击用例、检测规则
- **LLM 爱好者**: 改进提示词攻击库、支持更多模型
- **知识管理爱好者**: 改进 AI Brain 工具、添加新功能

详见 [CONTRIBUTING.md](./CONTRIBUTING.md) (即将添加)

---

## Roadmap | 路线图

- [x] LLM Security MVP (58 攻击用例)
- [x] AI Brain 工具集 (search/ask/archive/decision)
- [ ] 统一 BAS 平台
- [ ] 传统 BAS 攻击模块
- [ ] ATT&CK 映射
- [ ] 多语言支持

---

## License

[MIT License](./LICENSE)

---

## Links | 链接

- **Issues**: [GitHub Issues](https://github.com/JimMir45/AttackProbe/issues)
- **Discussions**: [GitHub Discussions](https://github.com/JimMir45/AttackProbe/discussions)
