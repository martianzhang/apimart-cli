# 图片生成

支持文生图、图生图、Inpainting 三种模式。

## 基本用法（文生图）

```bash
# 直接传提示词
apimart-cli image --prompt "一只猫在星空下"

# 从文件读取（自动识别文件路径）
apimart-cli image --prompt prompt.txt

# 从 stdin 读取（--prompt 不传时默认读 stdin）
echo "赛博朋克城市夜景" | apimart-cli image
apimart-cli image < prompt.txt

# 自动轮询并下载图片到当前目录
apimart-cli image --prompt "..."
```

## 参数

| 参数 | 短参 | 说明 |
|---|---|---|
| `--prompt` | `-p` | 文本描述（自动识别文件/stdin） |
| `--model` | `-m` | 模型名，默认 `gpt-image-2-official` |
| `--size` | `-s` | 宽高比，如 `16:9`、`1:1`，或像素如 `1024x1024` |
| `--resolution` | `-r` | 分辨率档：`1k`、`2k`、`4k` |
| `--quality` | `-q` | 质量：`auto`、`low`、`medium`、`high` |
| `--output-format` | `-f` | 输出格式：`png`、`jpeg`、`webp` |
| `--output-compression` | | 压缩率 0-100（jpeg/webp） |
| `--n` | | 生成数量 1-4 |
| `--image-url` | | 参考图片 URL（可重复） |
| `--mask-url` | | 蒙版图片 URL（inpainting） |
| `--background` | | 背景：`auto`、`opaque`、`transparent` |
| `--moderation` | | 审核强度：`auto`、`low` |
| `--dry-run` | | 打印 curl 不调用 API |

```bash
apimart-cli image --prompt "..." \
  --size "16:9" \
  --resolution "2k" \
  --quality "high" \
  --output-format "jpeg" \
  --output-compression 90 \
  --n 2
```

## JSON 输入

构建完整的请求 JSON 传入：

```bash
# JSON 文件
apimart-cli image --json request.json

# JSON 字符串
apimart-cli image --json '{"prompt":"a red fox","n":4}'

# 从 stdin
cat request.json | apimart-cli image --json -
```

## 参考图生图 (image-to-image)

参考已有图片进行融合或编辑，支持本地文件（自动上传）和远程 URL：

```bash
# 本地文件（自动上传到 APIMart）
apimart-cli image \
  --prompt "把这张照片改成吉卜力风格" \
  --image-url ./my-photo.jpg

# 远程 URL
apimart-cli image \
  --prompt "融合两张参考图，保留主要轮廓" \
  --image-url "https://example.com/img1.png" \
  --image-url "https://example.com/img2.png"
```

## Inpainting（蒙版替换）

提供原图和蒙版，替换指定区域：

```bash
# 本地文件自动上传
apimart-cli image \
  --prompt "把背景换成沙漠日落" \
  --image-url ./photo.png \
  --mask-url ./mask.png

# 远程 URL
apimart-cli image \
  --prompt "Replace background with desert sunset" \
  --image-url "https://example.com/photo.png" \
  --mask-url "https://example.com/mask.png"
```

> `--image-url` 和 `--mask-url` 可接受本地文件路径或远程 URL。
> 本地文件会自动通过 `POST /v1/uploads/images` 上传，上传后的 URL 有效期 72 小时。

## 最经济配置

参考 [APIMart 定价](https://apimart.ai/pricing)，`gpt-image-2-official` 最低 **$0.00144/张**：

```bash
apimart-cli image --prompt "..." \
  --size "3:1" \
  --resolution "1k" \
  --quality "low"
```

或写入 `~/.config/apimart/config.yaml` 作为默认值。

## 输出格式建议

| 格式 | 适用场景 |
|---|---|
| PNG | 需要透明背景、后续编辑、对画质要求高 |
| JPEG | 日常使用、社交媒体分享、网页展示 |
| WebP | 网页使用、需要小文件体积、支持透明 |
