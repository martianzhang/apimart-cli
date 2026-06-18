# 其他命令

## 查询模型列表

无需 API Key，免认证查询市场可用模型：

```bash
# 所有模型
apimart-cli models

# 按类型筛选
apimart-cli models image
apimart-cli models video
apimart-cli models chat

# 或使用 --type 参数
apimart-cli models --type image

# 查看特定模型的详细定价
apimart-cli models pricing gpt-image-2-official
apimart-cli models pricing doubao-seedance-2.0
```

## 查询任务状态

```bash
apimart-cli task task_01KV4KD9FBH3AZ4DE18A7Y17S3
```

返回完整的任务信息（状态、进度、耗时、费用、结果 URL 等）。图片任务完成后自动下载图片到 `--output` 目录。

## 查询余额

```bash
# 查询当前 API Key（Token）的余额
apimart-cli balance

# 查询用户账号的总余额
apimart-cli balance user
```

## Dry-run 调试

打印即将提交的 curl 命令，不实际调用 API：

```bash
# 图片 dry-run
apimart-cli image --prompt "test" --size "16:9" --dry-run

# 视频 dry-run
apimart-cli video --prompt "test" --duration 4 --dry-run
```

## 查看版本

```bash
apimart-cli version
# 或
apimart-cli --version
```

## API 参考

| 端点 | 用途 | 状态 |
|---|---|---|
| `POST /v1/chat/completions` | AI 对话 | ✅ |
| `POST /v1/images/generations` | 文生图 | ✅ |
| `POST /v1/videos/generations` | 文生视频 | ✅ |
| `POST /v1/uploads/images` | 上传图片 | ✅ |
| `GET /v1/tasks/{task_id}` | 查询任务状态 | ✅ |
| `GET /v1/balance` | Token 余额查询 | ✅ |
| `GET /v1/user/balance` | 用户余额查询 | ✅ |
| `GET /api/marketplace/models` | 模型列表（免认证） | ✅ |
| `GET /api/pricing/model` | 模型定价详情（免认证） | ✅ |

完整文档见 [docs.apimart.ai](https://docs.apimart.ai/en)。
