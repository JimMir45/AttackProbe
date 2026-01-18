# SecOpsHub 文档中心

> 安全运营中枢 - 文档导航

## 文档结构

```
docs/
├── secops-hub/          # SecOpsHub 整体项目
├── bas-platform/        # AttackProbe BAS 平台
├── llm-security/        # LLM 安全模块
└── ai-brain/            # AI-Brain 知识管理
```

---

## 项目文档

### [SecOpsHub 整体项目](./secops-hub/)

SecOpsHub 整体产品愿景和架构设计。

| 文档 | 说明 |
|------|------|
| [产品愿景](./secops-hub/00-立项/产品愿景.md) | 整体产品定位和发展路径 |
| [整体架构设计](./secops-hub/01-设计/整体架构设计.md) | 多模块整合架构 |
| [ADR-001 项目整合策略](./secops-hub/_decisions/ADR-001-项目整合策略.md) | 架构决策记录 |

---

## 模块文档

### [AttackProbe BAS 平台](./bas-platform/)

入侵与攻击模拟 (BAS) 平台，整合传统安全和 LLM 安全测试能力。

| 阶段 | 文档 |
|------|------|
| **立项** | [产品需求说明书](./bas-platform/00-立项/产品需求说明书.md) |
| | [六个月工作计划](./bas-platform/00-立项/六个月工作计划.md) |
| | [里程碑清单](./bas-platform/00-立项/里程碑清单.md) |
| **设计** | [系统架构设计](./bas-platform/01-设计/系统架构设计.md) |
| | [技术栈评估](./bas-platform/01-设计/技术栈评估.md) |
| **调研** | [BAS 行业调研](./bas-platform/research/BAS行业调研.md) |

### [LLM-Security 模块](./llm-security/)

LLM 大模型安全自动化测试工具。

| 文档 | 说明 |
|------|------|
| [模块说明](./llm-security/README.md) | 功能概述和快速开始 |

**源代码**: [packages/llm-security/](../packages/llm-security/)

### [AI-Brain 模块](./ai-brain/)

安全团队知识管理工具集。

| 文档 | 说明 |
|------|------|
| [模块说明](./ai-brain/README.md) | 工具集介绍和使用方法 |

**源代码**: [packages/ai-brain/](../packages/ai-brain/)

---

## 快速链接

| 需求 | 去哪里 |
|------|--------|
| 了解整体产品 | [SecOpsHub 产品愿景](./secops-hub/00-立项/产品愿景.md) |
| 了解 BAS 平台 | [AttackProbe 产品需求](./bas-platform/00-立项/产品需求说明书.md) |
| BAS 行业知识 | [BAS 行业调研](./bas-platform/research/BAS行业调研.md) |
| 使用 LLM 安全工具 | [LLM-Security 说明](./llm-security/README.md) |
| 使用知识管理工具 | [AI-Brain 说明](./ai-brain/README.md) |

---

## 文档规范

- 按生命周期阶段组织 (00-立项 → 01-设计 → 02-开发 → ...)
- ADR (架构决策记录) 存放在 `_decisions/` 目录
- 调研资料存放在 `research/` 目录
