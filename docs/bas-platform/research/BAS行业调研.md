# BAS (Breach and Attack Simulation) 行业调研

> 版本: v1.1 | 调研日期: 2026-01-18 | 状态: 完成

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

## 2. 目标行业分析

### 2.1 金融行业

**行业特点:**
- 监管要求严格 (等保2.0、银保监会要求)
- 数据敏感度高 (账户信息、交易数据)
- APT 攻击重点目标
- 内部威胁风险高

**高优先级威胁场景:**

| 威胁类型 | ATT&CK 战术 | 典型攻击 |
|---------|------------|---------|
| 凭据窃取 | TA0006 | 钓鱼获取员工凭据、暴力破解 |
| 数据泄露 | TA0010 | 客户信息外泄、交易数据窃取 |
| 勒索攻击 | TA0040 | 加密核心系统、勒索赎金 |
| 内部威胁 | TA0007/TA0009 | 越权访问、敏感数据收集 |
| 供应链攻击 | TA0001 | 第三方接口攻击、开源组件漏洞 |

**金融行业 BAS 重点:**
- 网银/手机银行 Web 安全测试
- 核心系统 API 安全验证
- 员工安全意识 (钓鱼模拟)
- DLP 有效性验证

### 2.2 运营商行业

**行业特点:**
- 关键基础设施，影响面广
- 网络架构复杂 (IT/CT 融合)
- 用户数据量大 (号码、位置、通话记录)
- 5G/物联网新场景多

**高优先级威胁场景:**

| 威胁类型 | ATT&CK 战术 | 典型攻击 |
|---------|------------|---------|
| 网络渗透 | TA0001/TA0008 | 边界突破、核心网渗透 |
| 用户数据泄露 | TA0009/TA0010 | 用户信息批量窃取 |
| 服务中断 | TA0040 | DDoS、核心系统破坏 |
| 横向移动 | TA0008 | 从 IT 网渗透到 CT 网 |
| 供应链攻击 | TA0001 | 设备供应商后门 |

**运营商行业 BAS 重点:**
- 网络边界防护验证
- BSS/OSS 系统安全测试
- 5G 核心网安全验证
- 物联网设备安全

---

## 3. 市场格局

### 3.1 主要厂商

| 厂商 | 产品 | 特点 | 评分 |
|------|------|------|------|
| Picus Security | Picus Platform | 攻击库丰富，检测验证强 | 9.0/10 |
| Cymulate | Cymulate Platform | 全面的攻击向量覆盖 | 8.8/10 |
| SafeBreach | SafeBreach Platform | "数字孪生"方法，Hacker's Playbook | 8.5/10 |
| AttackIQ | AttackIQ Platform | AI/ML 组件测试，Anatomic Engine | 8.3/10 |
| Pentera | Pentera Platform | 自动化渗透测试 | 8.2/10 |
| Horizon3.ai | NodeZero | 自主渗透测试 | 8.0/10 |

> 数据来源: PeerSpot 2025年12月排名

### 3.2 开源方案

| 项目 | 用途 | 链接 |
|------|------|------|
| MITRE Caldera | 自动化对手模拟 | https://caldera.mitre.org |
| Atomic Red Team | ATT&CK 技术测试库 | https://atomicredteam.io |
| Infection Monkey | 自动化渗透测试 | https://www.akamai.com/infectionmonkey |
| Garak (NVIDIA) | LLM 安全测试 | https://github.com/leondz/garak |

---

## 4. MITRE ATT&CK 覆盖策略

### 4.1 覆盖率说明

MITRE ATT&CK Enterprise 包含 **14 个战术、203 个技术、453 个子技术**。

**我们的策略: 务实覆盖 50%，其余提供技术原理文档**

| 覆盖类型 | 技术数量 | 说明 |
|---------|---------|------|
| **直接覆盖** | ~100 个 (50%) | 自动化模拟执行 |
| **原理文档** | ~100 个 (50%) | 技术原理 + 手动验证指南 |
| **不适用** | ~3 个 | 物理攻击等无法模拟 |

### 4.2 直接覆盖的技术 (50%)

基于金融和运营商行业需求，优先覆盖以下技术：

#### 第一优先级 - 核心 50 技术

| 战术 | 技术 ID | 技术名称 | 金融 | 运营商 |
|------|--------|---------|------|-------|
| **Initial Access** | T1566 | Phishing | ★★★ | ★★ |
| | T1190 | Exploit Public-Facing App | ★★★ | ★★★ |
| | T1078 | Valid Accounts | ★★★ | ★★ |
| **Execution** | T1059 | Command and Scripting | ★★ | ★★ |
| | T1204 | User Execution | ★★★ | ★★ |
| **Persistence** | T1547 | Boot/Logon Autostart | ★★ | ★★ |
| | T1053 | Scheduled Task/Job | ★★ | ★★ |
| | T1136 | Create Account | ★★ | ★★ |
| **Privilege Escalation** | T1548 | Abuse Elevation Control | ★★ | ★★ |
| | T1068 | Exploitation for Privilege | ★★ | ★★ |
| **Defense Evasion** | T1027 | Obfuscated Files | ★★ | ★★ |
| | T1070 | Indicator Removal | ★★ | ★★ |
| | T1055 | Process Injection | ★★ | ★ |
| **Credential Access** | T1003 | OS Credential Dumping | ★★★ | ★★ |
| | T1110 | Brute Force | ★★★ | ★★ |
| | T1552 | Unsecured Credentials | ★★★ | ★★ |
| **Discovery** | T1082 | System Information | ★★ | ★★ |
| | T1083 | File and Directory | ★★ | ★★ |
| | T1046 | Network Service Scan | ★★ | ★★★ |
| **Lateral Movement** | T1021 | Remote Services | ★★ | ★★★ |
| | T1570 | Lateral Tool Transfer | ★★ | ★★★ |
| **Collection** | T1005 | Data from Local System | ★★★ | ★★ |
| | T1560 | Archive Collected Data | ★★ | ★★ |
| **C2** | T1071 | Application Layer Protocol | ★★ | ★★ |
| | T1105 | Ingress Tool Transfer | ★★ | ★★ |
| **Exfiltration** | T1041 | Exfiltration Over C2 | ★★★ | ★★ |
| | T1048 | Exfiltration Over Alternative | ★★★ | ★★ |
| **Impact** | T1486 | Data Encrypted (Ransomware) | ★★★ | ★★ |
| | T1489 | Service Stop | ★★ | ★★★ |

#### 第二优先级 - 扩展 50 技术

包括但不限于：
- Web 攻击相关技术 (OWASP Top 10 映射)
- 云环境攻击技术
- 容器/K8s 相关技术
- API 安全相关技术

### 4.3 未直接覆盖的技术 (50%)

**未覆盖原因分类:**

| 原因 | 技术示例 | 处理方式 |
|------|---------|---------|
| **环境依赖** | T1558 Kerberoasting (需AD) | 提供技术原理 + 环境搭建指南 |
| **破坏性高** | T1485 Data Destruction | 仅标记模拟，提供验证方法 |
| **硬件依赖** | T1200 Hardware Additions | 提供检测建议，无法模拟 |
| **复杂度高** | T1557 MITM | 规划中，定制开发可支持 |
| **场景特殊** | T1195 Supply Chain | 提供检测方法，无法完全模拟 |

**技术原理文档包含:**
1. 攻击技术详解
2. 真实案例分析
3. 手动验证步骤
4. 防护建议
5. 检测方法

**定制开发说明:**
对于"规划中"和"定制开发可支持"的技术，企业客户可通过定制开发服务获得支持。

---

## 5. 技术架构

### 5.1 三层架构模型

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

### 5.2 核心组件

| 组件 | 职责 | 实现方式 |
|------|------|---------|
| **威胁情报库** | 存储 TTP (战术/技术/过程) | MITRE ATT&CK 映射 |
| **攻击树生成器** | 构建攻击路径 | 基于对手画像生成 |
| **Agent** | 在目标环境执行攻击 | 轻量级可执行程序 |
| **Abilities** | 具体攻击动作/技术 | 脚本/Payload |
| **判定引擎** | 评估攻击是否被防护 | 规则/关键词/AI |

---

## 6. 功能模块

### 6.1 完整 BAS 平台模块

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
├── LLM 安全 BAS (差异化能力)
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

### 6.2 功能优先级 (金融/运营商)

| 优先级 | 功能 | 金融 | 运营商 | 说明 |
|-------|------|------|-------|------|
| P0 | LLM 安全模块 | ★★ | ★★ | 已有，差异化竞争力 |
| P0 | Web 攻击模拟 | ★★★ | ★★ | 网银/门户安全 |
| P0 | 钓鱼模拟 | ★★★ | ★★ | 员工安全意识 |
| P1 | 网络攻击模拟 | ★★ | ★★★ | 边界防护 |
| P1 | ATT&CK 映射 | ★★★ | ★★★ | 合规要求 |
| P1 | 数据泄露模拟 | ★★★ | ★★★ | DLP 验证 |
| P2 | 端点攻击模拟 | ★★ | ★★ | 需要 Agent |
| P2 | 邮件攻击模拟 | ★★ | ★ | 需要邮件基础设施 |
| P3 | SIEM/SOAR 集成 | ★★★ | ★★★ | 企业级需求 |

---

## 7. 产品定位

### 7.1 差异化定位

| 竞品 | 我们的差异 |
|------|-----------|
| 商业 BAS (Picus, Cymulate) | 开源/私有化部署，价格优势，国产化 |
| 开源 BAS (Caldera) | 更易用的 UI，LLM 安全独有，中文支持 |
| LLM 安全工具 (Garak) | 更完整的 BAS 能力，企业级功能 |

**定位**: **面向金融和运营商行业的开源 BAS 平台，整合传统安全 + LLM 安全验证能力**

### 7.2 技术选型建议

| 组件 | 建议方案 | 原因 |
|------|---------|------|
| 攻击库格式 | YAML + MITRE ATT&CK ID | 标准化、可扩展 |
| Agent 架构 | Go 编译单文件 | 轻量、跨平台 |
| 报告格式 | Markdown + PDF | 灵活、专业 |
| API 标准 | OpenAPI 3.0 | 便于集成 |

---

## 8. 参考资料

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
| 2026-01-18 | v1.0 | 初稿创建 | SecOps Team |
| 2026-01-18 | v1.1 | 增加目标行业分析，明确50%覆盖策略 | SecOps Team |
