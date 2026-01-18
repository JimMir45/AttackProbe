# BAS (Breach and Attack Simulation) 行业调研

> 版本: v1.0 | 调研日期: 2026-01-18 | 状态: 完成

## 1. BAS 概述

### 1.1 什么是 BAS

BAS (Breach and Attack Simulation，入侵与攻击模拟) 是一种**自动化安全验证技术**，通过安全地模拟真实世界的网络攻击来持续评估组织的安全防御能力。

与传统渗透测试不同，BAS 具有以下特点：
- **持续性**: 可以自动化、持续运行，而非一次性评估
- **安全性**: 模拟攻击不会造成实际破坏
- **全面性**: 覆盖完整攻击链，从初始访问到数据外泄
- **可量化**: 提供防护有效性评分和差距分析

### 1.2 BAS 核心价值

| 价值点 | 说明 |
|-------|------|
| 持续验证 | 不是一次性扫描，而是 7x24 持续的安全态势验证 |
| 攻击视角 | 从攻击者角度验证防护有效性，发现真实差距 |
| 量化指标 | 提供可量化的安全防护成功率，便于汇报 |
| 主动防御 | 从被动响应转向主动发现和修复 |

---

## 2. 市场格局

### 2.1 主要厂商

| 厂商 | 产品 | 特点 | 评分 |
|------|------|------|------|
| Picus Security | Picus Platform | 攻击库丰富，检测验证强 | 9.0/10 |
| Cymulate | Cymulate Platform | 全面的攻击向量覆盖 | 8.8/10 |
| SafeBreach | SafeBreach Platform | "数字孪生"方法，Hacker's Playbook | 8.5/10 |
| AttackIQ | AttackIQ Platform | AI/ML 组件测试，Anatomic Engine | 8.3/10 |
| Pentera | Pentera Platform | 自动化渗透测试 | 8.2/10 |
| Horizon3.ai | NodeZero | 自主渗透测试 | 8.0/10 |

> 数据来源: PeerSpot 2025年12月排名

### 2.2 开源方案

| 项目 | 用途 | 链接 |
|------|------|------|
| MITRE Caldera | 自动化对手模拟 | https://caldera.mitre.org |
| Atomic Red Team | ATT&CK 技术测试库 | https://atomicredteam.io |
| Infection Monkey | 自动化渗透测试 | https://www.akamai.com/infectionmonkey |
| Garak (NVIDIA) | LLM 安全测试 | https://github.com/leondz/garak |

---

## 3. 技术架构

### 3.1 三层架构模型

根据学术研究和行业实践，现代 BAS 平台通常采用三层架构：

```
┌─────────────────────────────────────────────────────────────────┐
│                    SCE Orchestrator Layer                        │
│                      (编排层)                                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │ 威胁情报库   │  │混沌实验设计器│  │    攻击树生成器         │  │
│  │ (TTP库)     │  │             │  │  (Attack Tree Gen)      │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
└───────────────────────────┬─────────────────────────────────────┘
                            │ API
┌───────────────────────────▼─────────────────────────────────────┐
│                      Connector Layer                             │
│                       (连接层)                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │ 状态获取器   │  │ 操作管理器   │  │    结果收集器           │  │
│  │State Fetcher│  │Op Manager   │  │  Result Retriever       │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
└───────────────────────────┬─────────────────────────────────────┘
                            │ API
┌───────────────────────────▼─────────────────────────────────────┐
│                        BAS Layer                                 │
│                       (执行层)                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │   Agents    │  │  Abilities  │  │     Operations          │  │
│  │  (执行代理) │  │  (攻击能力) │  │   (攻击操作)            │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

### 3.2 核心组件

| 组件 | 职责 | 实现方式 |
|------|------|---------|
| **威胁情报库** | 存储 TTP (战术/技术/过程) | MITRE ATT&CK 映射 |
| **攻击树生成器** | 构建攻击路径 | 基于对手画像生成 |
| **Agent** | 在目标环境执行攻击 | 轻量级可执行程序 |
| **Abilities** | 具体攻击动作/技术 | 脚本/Payload |
| **判定引擎** | 评估攻击是否被防护 | 规则/关键词/AI |

---

## 4. MITRE ATT&CK 映射

### 4.1 Enterprise 矩阵 (14 个战术)

完整 BAS 平台应覆盖 MITRE ATT&CK Enterprise 的 14 个战术、203 个技术、453 个子技术：

| ID | 战术 (Tactic) | 中文 | 攻击目的 |
|---|--------------|------|---------|
| TA0043 | Reconnaissance | 侦察 | 收集可用于规划的信息 |
| TA0042 | Resource Development | 资源开发 | 建立攻击所需的资源 |
| TA0001 | Initial Access | 初始访问 | 进入目标网络 |
| TA0002 | Execution | 执行 | 运行恶意代码 |
| TA0003 | Persistence | 持久化 | 保持立足点 |
| TA0004 | Privilege Escalation | 权限提升 | 获取更高权限 |
| TA0005 | Defense Evasion | 防御规避 | 避免被检测 |
| TA0006 | Credential Access | 凭据访问 | 窃取账号密码 |
| TA0007 | Discovery | 发现 | 探索目标环境 |
| TA0008 | Lateral Movement | 横向移动 | 在网络中移动 |
| TA0009 | Collection | 收集 | 收集目标数据 |
| TA0011 | Command and Control | 命令控制 | 与被控系统通信 |
| TA0010 | Exfiltration | 数据外泄 | 窃取数据 |
| TA0040 | Impact | 影响 | 破坏/干扰系统 |

### 4.2 攻击覆盖建议

根据组织类型优先覆盖：

| 组织类型 | 优先战术 | 原因 |
|---------|---------|------|
| 金融机构 | TA0006, TA0010, TA0040 | 凭据窃取、数据泄露、勒索 |
| 制造业 | TA0002, TA0003, TA0040 | 工控系统攻击、持久化 |
| 互联网企业 | TA0001, TA0007, TA0009 | 初始访问、信息收集 |
| 政府机构 | TA0005, TA0008, TA0011 | APT、横向移动、C2 |

---

## 5. 功能模块

### 5.1 完整 BAS 平台模块

```
完整 BAS 平台
│
├── 传统安全 BAS
│   ├── 网络攻击模拟 (Network Attack)
│   │   ├── 网络渗透/端口扫描
│   │   ├── 防火墙规则绕过
│   │   ├── IDS/IPS 检测测试
│   │   └── 网络分段验证
│   │
│   ├── 端点攻击模拟 (Endpoint Attack)
│   │   ├── 恶意软件模拟 (Ransomware, RAT, Rootkit)
│   │   ├── EDR/AV 绕过测试
│   │   ├── 权限提升
│   │   └── 进程注入
│   │
│   ├── Web 攻击模拟 (Web Attack)
│   │   ├── OWASP Top 10 (XSS, SQLi, SSRF...)
│   │   ├── WAF 绕过测试
│   │   ├── API 安全测试
│   │   └── 认证/授权绕过
│   │
│   ├── 邮件攻击模拟 (Email Attack)
│   │   ├── 钓鱼邮件模拟
│   │   ├── 恶意附件测试
│   │   ├── 邮件网关绕过
│   │   └── BEC 攻击模拟
│   │
│   └── 数据泄露模拟 (Data Exfiltration)
│       ├── DLP 测试
│       ├── 隐蔽通道 (DNS/HTTPS隧道)
│       └── 云存储外泄
│
├── LLM 安全 BAS (新兴领域)
│   ├── Prompt Injection (提示词注入)
│   ├── Jailbreak (越狱攻击)
│   ├── Info Disclosure (敏感信息泄露)
│   ├── Training Data Extraction (训练数据提取)
│   └── Model Denial of Service (模型拒绝服务)
│
├── 安全控制验证
│   ├── NGFW/防火墙有效性
│   ├── WAF/Web应用防火墙
│   ├── IDS/IPS 检测率
│   ├── EDR/XDR 响应
│   ├── SIEM 告警覆盖
│   ├── DLP 数据防泄漏
│   └── Email Gateway 邮件网关
│
└── 报告与集成
    ├── ATT&CK 热力图
    ├── 防护有效性评分
    ├── 差距分析报告
    ├── 修复建议
    └── SIEM/SOAR 集成
```

### 5.2 功能优先级建议

| 优先级 | 功能 | 原因 |
|-------|------|------|
| P0 | LLM 安全模块 | 已有，差异化竞争力 |
| P0 | Web 攻击模拟 | 需求普遍，实现相对简单 |
| P1 | 网络攻击模拟 | 基础能力 |
| P1 | ATT&CK 映射 | 行业标准 |
| P2 | 端点攻击模拟 | 需要 Agent，复杂度高 |
| P2 | 邮件攻击模拟 | 需要邮件基础设施 |
| P3 | SIEM/SOAR 集成 | 企业级需求 |

---

## 6. 技术趋势 (2025-2026)

### 6.1 AI/ML 集成

- **AI 驱动的攻击模拟**: 使用 AI 生成更智能的攻击变体
- **预测性分析**: 预测潜在攻击路径 (如 Foreseeti)
- **AI 组件测试**: 测试企业 AI/ML 系统的安全性 (如 AttackIQ)

### 6.2 云原生 BAS

- **多云支持**: AWS/Azure/GCP 统一测试
- **容器安全**: Kubernetes 安全验证
- **云工作负载**: 运行时安全测试

### 6.3 持续验证

- **CI/CD 集成**: 安全左移，开发阶段验证
- **自动化修复**: 与 SOAR 联动自动修复
- **实时告警**: 与 SIEM 集成实时响应

---

## 7. 对 SecOpsHub 的建议

### 7.1 差异化定位

| 竞品 | 我们的差异 |
|------|-----------|
| 商业 BAS (Picus, Cymulate) | 开源/私有化部署，价格优势 |
| 开源 BAS (Caldera) | 更易用的 UI，LLM 安全独有 |
| LLM 安全工具 (Garak) | 更完整的 BAS 能力 |

**建议定位**: **国内首个整合传统安全 + LLM 安全的开源 BAS 平台**

### 7.2 路线图建议

```
Phase 1 (已完成)
├── LLM-Security MVP (58 个攻击用例)
└── 基础判定引擎

Phase 2 (建议)
├── Web 攻击模拟引擎
├── ATT&CK 基础映射
└── 统一 AttackProbe 平台

Phase 3 (未来)
├── 网络攻击模拟
├── Agent 架构
└── 完整 ATT&CK 覆盖

Phase 4 (远期)
├── 端点攻击模拟
├── SIEM/SOAR 集成
└── 企业级功能
```

### 7.3 技术选型建议

| 组件 | 建议方案 | 原因 |
|------|---------|------|
| 攻击库格式 | YAML + MITRE ATT&CK ID | 标准化、可扩展 |
| Agent 架构 | Go 编译单文件 | 轻量、跨平台 |
| 报告格式 | Markdown + PDF | 灵活、专业 |
| API 标准 | OpenAPI 3.0 | 便于集成 |

---

## 参考资料

- [MITRE ATT&CK Framework](https://attack.mitre.org/)
- [MITRE Caldera](https://caldera.mitre.org/)
- [Picus Security](https://www.picussecurity.com/)
- [Cymulate](https://cymulate.com/)
- [SafeBreach](https://www.safebreach.com/)
- [BAS Architecture Research (arXiv)](https://arxiv.org/html/2508.03882)
- [PeerSpot BAS Rankings 2025](https://www.peerspot.com/categories/breach-and-attack-simulation-bas)

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-01-18 | v1.0 | 初稿创建 | Claude |
