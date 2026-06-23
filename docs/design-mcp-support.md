# MCP 协议支持 — 设计文档

## 1. 背景

目前 APIMart 的 Image/Video API 是异步任务模式，而主流第三方工具（ComfyUI、OpenAI SDK 等）只支持 OpenAI 的同步格式。已有人通过 **genai-mcp**（MCP Server）将 APIMart 图片能力接入 AI Agent。我们的 CLI 已经有完整的 API 封装，加上 MCP 支持的成本很低，能为 Claude Desktop、Cursor 等 AI 工具提供原生的 APIMart 能力。

## 2. 目标

- 现有 CLI **追加** `mcp` 子命令，启动 MCP Server
- AI Agent 通过 MCP 协议直接调用 APIMart 的图片/视频/对话能力
- 不对现有命令做任何改动，不引入新依赖（除 MCP SDK）

## 3. 非目标

- ❌ 不替换现有的 `image`、`video`、`chat` 命令
- ❌ 不更改现有输出格式
- ❌ 不重新实现 MCP 协议（使用现成 SDK）
- ❌ 不做多 Provider 抽象

## 4. 架构

```
┌─────────────┐   stdio/HTTP    ┌────────────────────┐
│ Claude      │ ◄─────────────► │ apimart-cli mcp     │
│ Desktop     │                 │                     │
│ Cursor      │                 │ ├─ generate_image   │
│ VS Code     │                 │ ├─ generate_video   │
│ OpenClaw    │                 │ ├─ chat             │
│ ...         │                 │ ├─ list_models      │
└─────────────┘                 │ └─ get_balance      │
                                │                     │
                                │ ┌─────────────────┐ │
                                │ │ internal/client  │ │
                                │ │ (复用现有)        │ │
                                │ └─────────────────┘ │
                                └────────────────────┘
```

### 4.1 传输层

| 传输方式 | 第一阶段 | 后续 |
|----------|----------|------|
| stdio | ✅ 首选 | 永久支持 |
| Streamable HTTP | ❌ 不做 | 可按需追加 |

**选择理由**：stdio 是 MCP 最成熟的使用方式，MCP 客户端（Claude Desktop 等）以子进程方式启动 CLI，零网络配置，CLI 也无须引入 HTTP server。

### 4.2 MCP SDK

使用 `github.com/mark3labs/mcp-go`。它是目前 Go 生态最成熟的 MCP 实现，正在被讨论成为官方 SDK，有 2k+ stars，生产可用。

### 4.3 新增文件

```
cmd/mcp.go                  # cobra 子命令注册，约 30 行
internal/mcp/
  server.go                 # MCP Server 初始化、工具注册，约 100 行
  tools_generate.go         # image/video 工具实现，约 150 行
  tools_query.go            # models/balance/task 等查询工具，约 100 行
```

总计约 **380 行新增代码**，核心逻辑复用 `internal/client`。

## 5. 配置默认值传递策略（关键设计）

### 问题

CLI 配置文件（`config.yaml`）中有大量默认值：

```yaml
defaults:
  image:
    model: "gpt-image-2-official"
    size: "3:1"
    resolution: "1k"
    quality: "low"
    output_format: "png"
  video:
    model: "grok-imagine-1.5-video-apimart"
    duration: 5
    resolution: "480p"
```

这些默认值在 MCP 模式中需要：
1. **AI 知道它们是什么** — 否则 AI 以为全都是默认值，实际上用户定制过
2. **AI 不需要每次都填它们** — 否则浪费 token，且多余
3. **AI 在用户明确提及时可以覆盖** — 比如用户说「给我 4k 的」

### 方案：动态描述注入 + 精简 Schema

MCP 中，AI 读取工具定义的流程是：**`tools/list` 获取描述 + Schema → 缓存 → 按需调用**。描述是一次性读取的，token 开销仅一次。

核心做法：

```
┌─────────────────────────────────────────────┐
│ MCP Server 启动                              │
│                                              │
│  1. 加载 config.yaml                         │
│  2. 读取 defaults.image / defaults.video     │
│  3. 动态构建工具描述文本                       │
│  4. 注册工具（描述已含配置信息）                │
│                                              │
│    ↓                                         │
│                                              │
│ AI 收到工具定义：                              │
│  generate_image                              │
│  描述: "配置: model=gpt-image-2-official,     │
│          size=3:1, resolution=1k, quality=low │
│          除非用户明确要求，否则不要填这些参数"     │
│  Schema: {prompt, model?, size?, ...}        │
└─────────────────────────────────────────────┘
```

**工具描述动态构建示例**（启动时根据 config 生成）：

```
Generate images via APIMart.

当前配置:
  model = gpt-image-2-official | size = 3:1 | resolution = 1k
  quality = low | output_format = png | output_compression = (未设置)

策略: 参数已设好默认值，不要主动填写。只有用户在提示词中明确指定了某个参数时（如 "用 4k 分辨率"），才传入对应参数覆盖。
```

AI 读完这个描述就知道：
- 当前配置是什么
- 哪些是用户偏好，不需要 AI 决策
- 什么时候需要覆盖

### Schema 精简原则

保留全部参数（让 AI 有能力覆盖），但遵循以下规则：

| 参数 | 处理方式 |
|---|---|
| `prompt` | **required**，AI 必须填 |
| `model` | 保留，AI 很少需要改 |
| `size`, `resolution`, `quality`, `output_format` | 保留但描述注明「默认值见配置，仅在用户明确要求时覆盖」 |
| `image_urls`, `mask_url` | 保留，AI 根据场景自主决定 |
| `n` | 保留，AI 很少需要改 |

## 6. 工具定义

### 6.1 `generate_image`

```json
{
  "name": "generate_image",
  "description": "Generate images via APIMart.\n\n当前配置:\n  model = gpt-image-2-official | size = 3:1 | resolution = 1k\n  quality = low | output_format = png\n\n不要主动填写参数，除非用户在提示词中明确要求覆盖默认值。",
  "inputSchema": {
    "type": "object",
    "properties": {
      "prompt": { "type": "string", "description": "图片描述（必需）" },
      "model": { "type": "string", "description": "覆盖配置的模型名" },
      "size": { "type": "string", "description": "覆盖配置的尺寸" },
      "resolution": { "type": "string", "enum": ["1k", "2k", "4k"], "description": "覆盖配置的分辨率" },
      "quality": { "type": "string", "enum": ["auto", "low", "medium", "high"], "description": "覆盖配置的质量" },
      "image_urls": { "type": "array", "items": { "type": "string" }, "description": "参考图片 URL，用于图生图" },
      "mask_url": { "type": "string", "description": "蒙版 URL，用于 inpainting" },
      "output_format": { "type": "string", "enum": ["png", "jpeg", "webp"], "description": "覆盖配置的输出格式" },
      "n": { "type": "integer", "description": "生成数量 1-4，覆盖配置" }
    },
    "required": ["prompt"]
  }
}
```

**异步处理策略**：

| 工具 | 策略 |
|---|---|
| `generate_image` | **内部阻塞轮询**（复用 `client.PollTask`），通常 10-30s 完成。AI 一次调用拿结果。 |
| `generate_video` | **立即返回 task 信息**，不等待。AI 通过 `get_task` 工具稍后轮询结果。视频生成可能 60s+，阻塞不现实。 |

**图片返回格式**：图片自动下载到 `--output` 目录（默认当前目录），返回本地路径：

```json
{
  "content": [
    {
      "type": "text",
      "text": "图片已保存: image_task_xxx_0_0.png (12s, $0.00144)"
    }
  ]
}
```

> 为什么不返回 URL？因为图片有有效期（72h），下载到本地更可靠。AI 拿到本地路径后可以直接读取。

**视频返回格式**：返回 task_id，AI 稍后用 `get_task` 查询：

```json
{
  "content": [
    {
      "type": "text",
      "text": "视频任务已提交，task_id: task_xxx\n预计等待时间: 约 60s\n稍后使用 get_task 工具查询结果。"
    }
  ]
}
```

### 6.2 `generate_video`

```json
{
  "name": "generate_video",
  "description": "Generate videos via APIMart.\n\n当前配置:\n  model = grok-imagine-1.5-video-apimart | duration = 5s\n  resolution = 480p\n\n不要主动填写参数，除非用户在提示词中明确要求覆盖默认值。",
  "inputSchema": {
    "type": "object",
    "properties": {
      "prompt": { "type": "string", "description": "视频描述（必需）" },
      "model": { "type": "string", "description": "覆盖配置的模型名" },
      "duration": { "type": "integer", "description": "覆盖配置的时长（秒）" },
      "size": { "type": "string", "description": "覆盖配置的宽高比" },
      "resolution": { "type": "string", "enum": ["480p", "720p", "1080p"], "description": "覆盖配置的分辨率" },
      "image_urls": { "type": "array", "items": { "type": "string" }, "description": "参考图片 URL" },
      "video_urls": { "type": "array", "items": { "type": "string" }, "description": "参考视频 URL" },
      "generate_audio": { "type": "boolean", "description": "是否生成 AI 音频" }
    },
    "required": ["prompt"]
  }
}
```

同样阻塞轮询，返回视频 URL。

### 6.3 查询工具

| 工具名 | 描述 |
|---|---|
| `list_models` | 列出可用模型及类型（无需 API Key） |
| `get_model_pricing` | 查询特定模型的定价详情 |
| `get_balance` | 查询当前 API Key 余额 |
| `get_user_balance` | 查询用户账号的总余额 |
| `get_task` | 查询任务状态和结果，返回进度、费用、下载链接 |

### 6.4 查询工具定义

```json
{
  "name": "list_models",
  "description": "列出 APIMart 市场所有可用模型及其类型（chat/image/video）。无需 API Key。",
  "inputSchema": {
    "type": "object",
    "properties": {
      "type": { "type": "string", "enum": ["image", "video", "chat"], "description": "按类型筛选（可选）" }
    }
  }
}
```

```json
{
  "name": "get_model_pricing",
  "description": "查询指定模型的详细定价信息。无需 API Key。",
  "inputSchema": {
    "type": "object",
    "properties": {
      "model": { "type": "string", "description": "模型名称，如 gpt-image-2-official" }
    },
    "required": ["model"]
  }
}
```

```json
{
  "name": "get_balance",
  "description": "查询余额和用量。同时返回当前 API Key 的余额和用户账号的总余额。",
  "inputSchema": { "type": "object", "properties": {} }
}
```

```json
{
  "name": "get_task",
  "description": "查询异步任务（视频生成）的状态和结果。视频提交通常返回 task_id，用此工具轮询直到 status 为 completed。",
  "inputSchema": {
    "type": "object",
    "properties": {
      "task_id": { "type": "string", "description": "任务 ID，如 task_xxx" }
    },
    "required": ["task_id"]
  }
}
```

### 6.5 不包含的工具

| 工具 | 原因 |
|---|---|
| `chat` | 大模型本身具备对话能力，不需要嵌套调用 |

## 7. 启动时配置注入流程

```
MCP Server 启动
  │
  ├─ 加载 config.yaml / 环境变量
  │
  ├─ 读取 defaults.image / defaults.video / defaults.chat
  │
  ├─ 构建动态工具描述
  │   描述文本中嵌入当前配置值
  │
  └─ 注册工具 → 等待 tools/list 请求
                       │
                       ▼
             AI 收到含配置信息的工具描述
```

**实现要点**：描述文本在 `server.go` 初始化时拼装一次，后续 `tools/list` 返回的是同一个静态结果，没有运行时开销。

```go
func buildToolDescriptions(cfg *config.Config) {
    imageDesc := fmt.Sprintf(`Generate images via APIMart.

当前配置:
  model = %s | size = %s | resolution = %s
  quality = %s | output_format = %s

不要主动填写参数，除非用户在提示词中明确要求覆盖默认值。`,
        cfg.Defaults.Image.Model, cfg.Defaults.Image.Size,
        cfg.Defaults.Image.Resolution, cfg.Defaults.Image.Quality,
        cfg.Defaults.Image.OutputFormat)
    // ...
}
```

## 8. 启动方式

### stdio 模式（默认）

```bash
apimart-cli mcp
```

MCP Host 在配置中添加：

```json
{
  "mcpServers": {
    "apimart": {
      "command": "apimart-cli",
      "args": ["mcp"]
    }
  }
}
```

> 就是这么简单。没有端口、没有网络配置、没有 Docker。

## 9. 错误处理

| 情况 | 行为 |
|---|---|
| API Key 未配置 | 返回 MCP 错误，提示用户设置 |
| 图片/视频生成超时（180s） | 返回 MCP 错误，提示超时 |
| 余额不足 | 返回 MCP 错误，包含余额信息 |
| 网络错误 | 返回 MCP 错误，透传错误原因 |

## 10. 实施计划

| 阶段 | 内容 | 预估 |
|---|---|---|
| 1 | 安装 MCP SDK，创建 `cmd/mcp.go` 子命令 | 1h |
| 2 | 实现 `internal/mcp/server.go` — Server 初始化 + 工具注册 | 1h |
| 3 | 实现 `internal/mcp/tools_generate.go` — image + video 工具 | 3h |
| 4 | 实现 `internal/mcp/tools_chat.go` + `tools_query.go` | 2h |
| 5 | 本地测试：用 Claude Desktop / inspector 验证 | 1h |
| 6 | 文档：更新 README，新增 MCP 集成指南 | 1h |
| **合计** | | **~9h** |

## 11. 待确认的问题

1. **generate_image `n > 1` 时返回策略？** 多张图是分别列出路径还是合并描述？
2. **图片下载后的文件名格式？** 当前使用 `image_{task_id}_{i}_{j}.{ext}` 格式（v0.1.1 起），足够清晰。

---

**已确认的设计决策**：

| 问题 | 决定 |
|---|---|
| 是否包含 chat 工具 | ❌ 不包含 |
| 配置默认值如何传递给 AI | 启动时读取 config，动态注入工具描述文本 |
| AI 是否需要填所有参数 | 不需要。描述中告知 AI「不要主动填写，除非用户明确要求」 |
| AI 能否覆盖默认值 | 可以。Schema 保留全部参数，AI 按需传入 |
| 需要实现的工具 | `generate_image`, `generate_video`, `list_models`, `get_model_pricing`, `get_balance`（含 token+user）, `get_task` |
| 图片处理策略 | 生成后自动下载到 `--output` 目录，返回本地路径 |
| 视频处理策略 | 异步提交 → 返回 task_id → AI 通过 `get_task` 轮询 |
| 图片下载 | 复用现有逻辑，自动下载到 `--output` 目录（默认当前目录） |

---
