# MCP 集成 🧪

> **当前状态**：测试中，功能和接口可能变化，欢迎试用反馈。

`apimart-cli` 支持 [MCP 协议](https://modelcontextprotocol.io/)（Model Context Protocol），允许 AI 代理（Claude Desktop、Cursor、VS Code 等）直接调用 API 能力。

## 快速配置

在 MCP 客户端配置中添加：

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

## 可用工具

| 工具名 | 描述 | 是否需要 API Key |
|---|---|---|
| `generate_image` | 图片生成（文生图、图生图、Inpainting） | ✅ |
| `generate_video` | 视频生成（提交后返回 task_id，异步查询） | ✅ |
| `list_models` | 列出市场可用模型，支持按类型筛选 | ❌ |
| `get_model_pricing` | 查询指定模型的定价详情 | ❌ |
| `get_balance` | 查询 Token 和用户余额 | ✅ |
| `get_task` | 查询异步任务状态和结果 | ✅ |

## 配置

MCP 模式复用现有配置体系，支持三种方式：

```bash
# 方式一：配置文件
# ~/.config/openai/config.yaml 或 ~/.config/apimart/config.yaml

# 方式二：环境变量
OPENAI_API_KEY=sk-xxx apimart-cli mcp

# 方式三：CLI 参数
apimart-cli mcp --api-key sk-xxx --output ./downloads
```

## 工具描述动态注入

启动时自动读取 `config.yaml` 中的默认值（model、size、resolution、quality 等），并注入到工具描述中。AI 代理因而知道当前配置，只会在用户明确要求时才覆盖参数。
