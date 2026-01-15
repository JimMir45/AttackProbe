# 02-开发阶段

> 本目录存放开发阶段相关文档

## 目录结构

```
02-开发/
├── 开发规范.md          # 编码规范、Git规范
├── 功能模块/            # 各功能模块设计
│   ├── 目标管理.md
│   ├── 用例管理.md
│   ├── 任务中心.md
│   └── 报告中心.md
├── 开发日志/            # 开发过程记录
└── 问题追踪.md          # 开发中遇到的问题
```

## 复用资源

### llm-security-bas 可复用模块

| 模块 | 路径 | 说明 |
|------|------|------|
| LLM客户端 | `pkg/llm/` | OpenAI兼容API客户端 |
| 判定引擎 | `pkg/judge/` | 关键词判定逻辑 |
| 执行器 | `internal/service/executor.go` | Worker Pool实现 |
| 内置用例 | `internal/model/builtin_cases.go` | 58个LLM攻击用例 |
| 前端框架 | `web/` | Vue3 + Element Plus |

## 待补充

- [ ] 编码规范
- [ ] Git分支策略
- [ ] CI/CD配置
