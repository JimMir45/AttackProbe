# SecOpsHub

> Open-source Security Operations Platform | 开源安全运营平台

SecOpsHub 是一个面向安全运营团队的开源工具平台，包含三大核心模块：

- **BAS Platform** - 攻击模拟有效性验证
- **LLM Security** - 大模型安全测试
- **SecOps Brain** - 安全团队知识管理

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

### 2. SecOps Brain (ai-brain)

安全团队知识管理工具集，帮助安全团队沉淀经验、追溯决策。

**Tools | 工具:**

| Tool | Description |
|------|-------------|
| `brain-search` | 全文搜索 - 搜索安全知识库文档 |
| `brain-ask` | RAG 问答 - 基于向量检索的智能问答 |
| `brain-archive` | 对话归档 - 自动归档安全讨论为结构化文档 |
| `brain-decision` | 决策记录 - 安全架构决策记录(ADR)管理 |

**Quick Start | 快速开始:**
```bash
cd packages/ai-brain/brain-search
go build -o brain-search .
./brain-search "漏洞修复"
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
SecOpsHub/
├── packages/
│   ├── llm-security/       # LLM 安全测试模块
│   │   ├── cmd/            # 程序入口
│   │   ├── internal/       # 内部实现
│   │   ├── pkg/            # 公共包
│   │   └── web/            # 前端源码
│   └── ai-brain/           # 安全团队知识管理
│       ├── brain-search/   # 全文搜索
│       ├── brain-ask/      # RAG 问答
│       ├── brain-archive/  # 对话归档
│       └── brain-decision/ # 决策记录
├── docs/                   # 项目文档
│   ├── 00-立项/
│   ├── 01-设计/
│   └── ...
└── LICENSE
```

---

## Why SecOpsHub? | 为什么选择 SecOpsHub

| Feature | Description |
|---------|-------------|
| **安全运营聚焦** | 专为安全团队设计的工具集 |
| **中文原生** | 针对中文安全场景和中文 LLM 优化 |
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

# 构建 SecOps Brain 工具
cd packages/ai-brain/brain-search
go build -o brain-search .
```

---

## Contributing | 贡献

欢迎贡献代码、文档、Issue 和建议！

- **安全运营工程师**: 分享安全运营最佳实践
- **安全研究者**: 贡献攻击用例、检测规则
- **LLM 爱好者**: 改进提示词攻击库、支持更多模型

详见 [CONTRIBUTING.md](./CONTRIBUTING.md) (即将添加)

---

## Roadmap | 路线图

- [x] LLM Security MVP (58 攻击用例)
- [x] SecOps Brain 工具集 (search/ask/archive/decision)
- [ ] 统一 BAS 平台
- [ ] 传统 BAS 攻击模块
- [ ] ATT&CK 映射
- [ ] 安全运营 Dashboard

---

## License

[MIT License](./LICENSE)

---

## Links | 链接

- **Issues**: [GitHub Issues](https://github.com/JimMir45/SecOpsHub/issues)
- **Discussions**: [GitHub Discussions](https://github.com/JimMir45/SecOpsHub/discussions)
