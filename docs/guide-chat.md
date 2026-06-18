# AI 对话

支持流式输出（默认），兼容 OpenAI 格式，可使用 GPT-5、Claude、Gemini、DeepSeek 等模型。

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
apimart-cli chat --model gpt-5 --message "Hello"
```

## 参数

| 参数 | 说明 |
|---|---|
| `--message` | 用户消息（可重复，实现多轮对话） |
| `--system` | 系统提示词，设定 AI 角色 |
| `--model` | 模型名，默认 `deepseek-v4-flash`（最便宜） |
| `--temperature` | 采样温度 0-2，默认 1.0 |
| `--max-tokens` | 最大生成 token 数 |
| `--no-stream` | 关闭流式输出，等待完整响应 |
| `--json` | JSON 输入 |
