# apimart-cli

APIMart 图片生成 CLI 工具。支持 GPT-Image-2 等模型的参数配置，通过命令行参数或 JSON 输入快速生成图片。

## 安装

```bash
go install github.com/martianzhang/apimart-cli@latest
```

或从源码构建：

```bash
git clone https://github.com/martianzhang/apimart-cli.git
cd apimart-cli
make build
```

## 配置

### 1. API Key

设置 API Key 有三种方式（优先级从高到低）：

```bash
# 方式一：命令行参数
apimart-cli generate --prompt "..." --api-key "sk-xxx"

# 方式二：环境变量
export APIMART_API_KEY="sk-xxx"

# 方式三：配置文件 ~/.config/apimart/config.yaml
api_key: "sk-xxx"
```

### 2. 配置文件

配置文件位于 `~/.config/apimart/config.yaml`，可设置默认参数：

```yaml
api_key: "sk-xxx"

# HTTP 代理（支持 http/https/socks5）
# 也可通过 APIMART_HTTP_PROXY 环境变量或 --http-proxy 参数设置
http_proxy: "http://127.0.0.1:7890"

defaults:
  model: "gpt-image-2-official"
  size: "3:1"           # 3:1 + 1k → 1536x512
  resolution: "1k"      # 最低分辨率
  quality: "low"        # 最经济质量档
  output_format: "png"
```

完整示例见 [config.example.yaml](config.example.yaml)。

## 使用

### 基本用法

```bash
# 直接传提示词
apimart-cli generate --prompt "一只猫在星空下"

# 从文件读取提示词（自动识别）
apimart-cli generate --prompt prompt.txt

# 从 stdin 读取
echo "赛博朋克城市夜景" | apimart-cli generate --prompt -

# 等待结果并下载
apimart-cli generate --prompt "..." --wait --output ./images
```

### 详细参数

```bash
apimart-cli generate --prompt "..." \
  --size "16:9" \
  --resolution "2k" \
  --quality "high" \
  --output-format "jpeg" \
  --output-compression 90 \
  --n 2
```

### JSON 输入

```bash
# JSON 文件
apimart-cli generate --json request.json

# JSON 字符串
apimart-cli generate --json '{"prompt":"a red fox","n":4}'

# 从 stdin
cat request.json | apimart-cli generate --json -
```

### 代理

```bash
# 命令行指定
apimart-cli generate --prompt "..." --http-proxy "http://127.0.0.1:7890"

# 或设置环境变量（支持 HTTP_PROXY / HTTPS_PROXY / ALL_PROXY / NO_PROXY）
export HTTP_PROXY="http://127.0.0.1:7890"

# 也支持 SOCKS5
apimart-cli generate --prompt "..." --http-proxy "socks5://127.0.0.1:1080"
```

### 最经济配置

参考 [APIMart 定价](https://apimart.ai/pricing)，`gpt-image-2-official` 最低 **$0.00144/张**：

```bash
apimart-cli generate --prompt "..." \
  --size "3:1" \
  --resolution "1k" \
  --quality "low"
```

或写入 config.yaml 作为默认值。

## API 参考

本工具封装了 APIMart 的 GPT-Image-2 接口：

| 端点 | 用途 |
|---|---|
| `POST /v1/images/generations` | 提交图片生成任务 |
| `GET /v1/tasks/{task_id}` | 查询任务状态 |
| `GET /v1/balance` | 查询余额 |

完整 API 文档见 [docs.apimart.ai](https://docs.apimart.ai/en/api-reference/images/gpt-image-2/official)。

## Makefile

```bash
make          # 编译
make run ARGS="generate --help"   # 运行
make clean    # 清理产物
make lint     # 静态检查
```

## 优先级规则

**CLI 参数 > JSON 输入 > YAML 配置 > 代码默认值**

代理的优先级：
**`--http-proxy` 参数 > `APIMART_HTTP_PROXY` 环境变量 > `HTTP_PROXY` 标准环境变量**
