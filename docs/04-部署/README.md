# 04-部署阶段

> 本目录存放部署与运维相关文档

## 目录结构

```
04-部署/
├── 部署手册.md          # 安装部署指南
├── 维护手册.md          # 日常运维指南
├── 配置说明.md          # 配置项说明
└── 故障排查.md          # 常见问题处理
```

## 部署方式

### 单文件部署 (推荐)

```bash
# 下载
wget https://github.com/xxx/AttackProbe/releases/download/v1.0.0/AttackProbe-linux-amd64

# 运行
chmod +x AttackProbe-linux-amd64
./AttackProbe-linux-amd64 -port 8080
```

### Docker部署

```bash
docker run -d -p 8080:8080 -v ./data:/app/data AttackProbe:1.0.0
```

## 系统要求

| 项目 | 最低配置 | 推荐配置 |
|------|---------|---------|
| CPU | 2核 | 4核 |
| 内存 | 4GB | 8GB |
| 磁盘 | 20GB | 50GB |
| 系统 | Linux/Windows/macOS | Linux |

## 待补充

- [ ] 详细部署步骤
- [ ] 配置项文档
- [ ] 监控配置
- [ ] 备份恢复
