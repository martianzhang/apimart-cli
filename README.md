# apimart-cli

[![CI](https://github.com/martianzhang/apimart-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/martianzhang/apimart-cli/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/martianzhang/apimart-cli)](https://go.dev/)
[![License](https://img.shields.io/github/license/martianzhang/apimart-cli)](LICENSE)
[![Release](https://img.shields.io/github/v/release/martianzhang/apimart-cli)](https://github.com/martianzhang/apimart-cli/releases)

APIMart API 的统一命令行工具。支持 **图片生成**、**视频生成** 和 **AI 对话**。

## 快速开始

```bash
# 安装
go install github.com/martianzhang/apimart-cli@latest

# 设置 API Key
export APIMART_API_KEY="sk-xxx"

# 生成一张图
apimart-cli image --prompt "一只猫在星空下"

# 看更多用法
apimart-cli image --help
```

## 命令

```
apimart-cli
├── image      图片生成 →  docs/guide-image.md
├── video      视频生成 →  docs/guide-video.md
├── chat       AI 对话  →  docs/guide-chat.md
├── models     模型列表及定价
├── task       查询任务状态
├── balance    查询余额
└── mcp        MCP Server 🧪 →  docs/mcp.md
```

## 文档

| 文档 | 内容 |
|---|---|
| [安装与配置](docs/installation.md) | 安装、API Key、配置文件、代理 |
| [图片生成](docs/guide-image.md) | 全部参数、图生图、Inpainting、JSON 输入 |
| [视频生成](docs/guide-video.md) | 全部参数、首尾帧、参考视频 |
| [AI 对话](docs/guide-chat.md) | 流式输出、多轮对话、参数说明 |
| [其他命令](docs/guide-commands.md) | models、task、balance、dry-run、API 参考 |
| [MCP 集成](docs/mcp.md) | AI 代理（Claude/Cursor）集成指南 |

## MCP 集成 🧪（测试中）

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

详见 [docs/mcp.md](docs/mcp.md)。

## 优先级规则

**CLI 参数 > JSON 输入 > YAML 配置 > 代码默认值**

代理优先级：
**`--http-proxy` 参数 > `APIMART_HTTP_PROXY` 环境变量 > `HTTP_PROXY` 标准环境变量**

## License

[MIT](LICENSE)

<div align="center">

<img src="wechat.jpg" width="400" alt="没有那多" />

**欢迎关注微信公众号“没有那多”，分享更多好用、好玩的工具**

</div>
