# AI 对话

支持流式输出（默认），完全兼容 OpenAI 格式，支持 OpenRouter，可使用 GPT、Claude、Gemini、DeepSeek 等模型。

## 基本用法

```bash
# 基本对话（流式输出）
apimart-cli chat --message "你好，请介绍一下自己"

# 系统提示词
apimart-cli chat --system "你是一位诗人" --message "写一首关于AI的诗"

# 多轮对话
apimart-cli chat \
  --message "什么是机器学习？" \
  --message "能举个例子吗？"

# 非流式输出
echo "Explain Go in 3 words" | apimart-cli chat --no-stream

# 指定模型
apimart-cli chat --model gpt-4o --message "Hello"
```

### 使用 OpenRouter

```bash
export OPENAI_API_KEY="sk-or-xxx"
export OPENAI_BASE_URL="https://openrouter.ai/api/v1"

apimart-cli chat --model "openai/gpt-4o" --message "Hello"
# OpenRouter 会额外显示费用: ... | Cost: $0.000014
```

### 使用任意 OpenAI 兼容中转

```bash
export OPENAI_API_KEY="sk-xxx"
export OPENAI_BASE_URL="https://your-relay.com/v1"

apimart-cli chat --message "Hello"
```

## 参数

| 参数 | 说明 |
|---|---|
| `--message` | 用户消息（可重复，实现多轮对话） |
| `--system` | 系统提示词，设定 AI 角色 |
| `--model` | 模型名，默认读取配置文件或 `deepseek-v4-flash` |
| `--temperature` | 采样温度 0-2，默认 1.0 |
| `--max-tokens` | 最大生成 token 数 |
| `--no-stream` | 关闭流式输出，等待完整响应 |
| `--json` | JSON 输入 |
