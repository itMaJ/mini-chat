# mini-chat

一个基于 Go、Gin 和 WebSocket 的简易聊天室，前端使用 Vue 3、TypeScript 和 Vite。

目前已支持填写昵称加入公共聊天室、实时收发 JSON 消息，并展示昵称、发送时间和连接状态。所有连接共用一个公共聊天室，尚未接入用户账户、JWT 或数据库。

前后端代码分别组织，当前构建后的前端页面由 Gin 提供，暂缓独立部署。开发时使用 Vite 提供页面并代理 WebSocket 请求。

## 第一阶段功能进度

| 功能 | 具体描述 | 优先级 | 当前状态 |
| --- | --- | --- | --- |
| 用户注册与登录 | 注册、登录及 JWT 身份认证 | 必须 | 未实现，目前仅使用昵称 |
| 聊天室管理 | 创建、加入、退出独立聊天室 | 必须 | 未实现，目前所有连接共用一个 Hub |
| 实时消息 | 通过 WebSocket 实时发送和接收消息 | 必须 | 已实现 |
| 消息持久化 | 使用 MySQL 保存历史聊天记录 | 必须 | 未实现 |
| 历史消息查询 | 分页加载聊天记录 | 必须 | 未实现，刷新页面后消息列表清空 |
| 在线状态 | 显示其他用户在线、离线状态 | 必须 | 未实现，目前仅展示自身连接状态 |
| 私聊功能 | 用户之间一对一聊天 | 推荐 | 未实现 |
| 群聊功能 | 多用户群组聊天 | 推荐 | 部分实现，支持公共广播，尚无独立群组 |

## 已有能力

- Gin 路由：健康检查、WebSocket 接口及前端静态文件服务。
- Hub 统一管理连接注册、注销和消息广播，每条连接分别处理读写。
- JSON 消息包含消息类型、正文、昵称和服务端生成的 UTC 时间。
- 页面提供昵称输入、消息列表、自动滚动、错误提示和连接状态展示；断开后可手动重新加入。
- 页面校验昵称长度为 1～20 个 Unicode 字符；昵称目前只用于展示，不代表经过认证的用户身份。
- 服务端过滤空正文和不支持的消息类型，限制单条入站 WebSocket 消息为 4KB，设置写超时并清理发送队列已满的慢连接。
- 前端按组件、聊天状态、通信服务和消息类型拆分，支持 Vue 类型检查与生产构建。

## 技术选型

| 技术 | 用途 | 状态 |
| --- | --- | --- |
| Go 1.27.3 | 后端语言，版本见 `go.mod` | 已使用 |
| Gin | HTTP 路由、日志及异常恢复 | 已使用 |
| Gorilla WebSocket | 实时双向通信 | 已使用 |
| Vue 3 + TypeScript | 前端界面与类型定义 | 已使用 |
| Vite | 前端开发服务器、WebSocket 代理和生产构建 | 已使用 |
| MySQL | 用户、聊天室及消息存储 | 计划接入 |
| JWT | 用户身份认证 | 计划接入 |

## 项目结构

```text
mini-chat/
├── cmd/server/main.go          # 创建 Hub、初始化路由并启动服务
├── internal/
│   ├── server/router.go       # Gin 路由和静态文件服务
│   └── chat/
│       ├── cilent.go          # 单连接管理、消息读取与写入
│       ├── hub.go             # 连接集合与消息广播
│       └── message.go         # JSON 消息结构
├── web/
│   ├── index.html            # Vue 挂载入口
│   ├── src/
│   │   ├── main.ts           # 启动 Vue 应用
│   │   ├── App.vue           # 组合聊天界面
│   │   ├── components/       # 昵称、状态、消息列表与发送表单
│   │   ├── composables/      # 聊天状态和操作
│   │   ├── services/         # WebSocket 通信与消息校验
│   │   ├── types/            # 消息和连接状态类型
│   │   └── styles/           # 页面样式
│   ├── package.json          # 前端依赖和命令
│   ├── package-lock.json     # 前端依赖锁文件
│   ├── tsconfig.json         # TypeScript 与 Vue 类型检查配置
│   └── vite.config.ts        # Vue 插件、开发端口和代理配置
├── .node-version             # Node.js 24
├── .gitignore                # 忽略依赖、构建产物和本地配置
├── go.mod
├── go.sum
└── readme.md
```

## 环境准备

- Go 1.27.3。
- Node.js 24，版本声明见 `.node-version`，以及 npm。
- 当前运行无需配置 MySQL 或 JWT 密钥。

以下命令均以项目根目录为起点。

安装依赖：

```bash
go mod download
npm --prefix web ci
```

## 开发运行

终端一，在项目根目录启动后端：

```bash
go run ./cmd/server
```

终端二，在项目根目录启动前端：

```bash
npm --prefix web run dev
```

浏览器访问 `http://localhost:5173/`。使用 VS Code 远程开发时，转发 5173 端口并打开对应地址。

开发页面由 Vite 提供，浏览器发出的 `/ws` 请求由 Vite 转发到 Go 服务的 8765 端口。当前 Vite 配置仅代理 `/ws`；后续接入注册登录等 HTTP API 时，再增加 `/api` 代理。

## 构建并由 Gin 提供页面

在项目根目录执行：

```bash
npm --prefix web run build
go run ./cmd/server
```

浏览器访问 `http://localhost:8765/`。Gin 会提供 `web/dist/index.html` 和 `web/dist/assets` 中的构建产物，运行时无需启动 Vite。

Go 服务需要从项目根目录启动，因为静态文件路径为相对路径。修改前端源码后，需要重新执行构建，8765 上的页面才会更新。开发源码入口 `/src/main.ts` 应通过 Vite 访问。

`web/node_modules`、`web/dist` 和本地环境配置已加入 Git 忽略规则；`package-lock.json` 和 `.node-version` 应提交到版本管理。

## 当前接口

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| GET | `/api/health` | 返回 `{"status":"ok"}`，用于检查 HTTP 服务 |
| GET | `/ws?username=Alice` | WebSocket 握手，昵称通过查询参数传入 |
| GET | `/` | 提供构建后的首页 |
| GET | `/assets/*filepath` | 提供构建后的 JavaScript、CSS 等资源 |

健康检查：

```bash
curl http://localhost:8765/api/health
```

客户端发送的聊天消息：

```json
{
  "type": "chat",
  "content": "大家好"
}
```

服务端广播的消息示例：

```json
{
  "type": "chat",
  "username": "Alice",
  "content": "大家好",
  "sent_at": "2026-11-01T08:00:00Z"
}
```

消息会广播给所有连接，包括发送者。4KB 入站限制针对整条 JSON 消息的字节数。

## 检查与验证

在项目根目录执行：

```bash
go test ./cmd/... ./internal/...
npm --prefix web run typecheck
npm --prefix web run build
```

目前尚无 Go 自动化测试用例，Go 命令主要检查项目包能否编译。前端 `build` 已包含 Vue 类型检查。ESLint 和 Prettier 已安装，尚未配置对应规则和运行脚本。

手动验证可以打开两个标签页，使用不同昵称加入并互发消息；停止 Go 服务后检查断开提示，重启后手动重新加入。

## 后续开发顺序

1. 接入 MySQL，建立用户表及数据库访问层。
2. 实现用户注册和登录，保存密码哈希。
3. 接入 JWT，提供当前用户接口，并在 WebSocket 升级前校验身份；聊天身份改为来自认证结果。
4. 实现独立聊天室的创建、加入、退出及成员权限校验。
5. 实现消息持久化和历史消息分页查询。
6. 实现用户在线列表、上线和离线通知。
7. 按需补充私聊、独立群组以及前后端独立部署。

随着功能增加，补充 WebSocket 心跳保活、连接到期处理、代码规范检查和自动化测试。
