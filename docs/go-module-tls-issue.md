# Go Module 下载 TLS 证书问题

## 现象

在公司网络环境下执行 `go get` 时报错：

```
go: module github.com/openai/openai-go/v3: Get "https://proxy.golang.org/github.com/openai/openai-go/v3/@v/list":
  tls: failed to verify certificate: x509: "*.xxx.com" certificate is not trusted
```

## 原因

公司网络（xxx.com）部署了 **SSL 中间人检测（TLS Inspection）**，防火墙会拦截 HTTPS 请求并替换 TLS 证书。

Go 的默认模块代理 `proxy.golang.org` 在通过公司网络访问时，返回的证书被防火墙替换为 `*.xxx.com`，Go 不信任该证书，因此拒绝连接。

此外，当使用 `GOPROXY=direct` 直连源仓库时，`golang.org/x/*` 等包托管在 Google 服务器上，公司网络 IPv6 连接 Google 会超时。

## 解决方案

设置以下 Go 环境变量，使用国内代理绕过防火墙：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GONOSUMCHECK=*
go env -w GONOSUMDB=off
```

### 各变量说明

| 变量 | 值 | 说明 |
|------|-----|------|
| `GOPROXY` | `https://goproxy.cn,direct` | 使用国内 Go 代理（七牛云），失败时回退到直连 |
| `GONOSUMCHECK` | `*` | 跳过所有模块的校验和验证 |
| `GONOSUMDB` | `off` | 关闭 sum.golang.org 校验数据库查询，避免 IPv6 连接超时 |

### 其他可用的国内 Go 代理

| 代理 | 地址 |
|------|------|
| 七牛云 | `https://goproxy.cn` |
| 阿里云 | `https://mirrors.aliyun.com/goproxy/` |
| 字节跳动 | `https://goproxy.io` |
| 官方（需科学上网） | `https://proxy.golang.org` |

### 恢复默认设置

如需恢复默认配置：

```bash
go env -u GOPROXY
go env -u GONOSUMCHECK
go env -u GONOSUMDB
```

## 相关文件

- Go 环境配置：`~/.config/go/env`（或通过 `go env` 查看）
- 系统代理：系统偏好设置 → 网络 → 高级 → 代理
