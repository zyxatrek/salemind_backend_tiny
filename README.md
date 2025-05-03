# 图像和视频生成服务

这是一个基于Golang的图像和视频生成服务，支持以下功能：

1. 根据用户输入的参数生成图像
2. 将生成的图像转换为视频

## 功能特点

- 使用Qwen API生成自定义提示词
- 使用Liblibai API生成高质量图像
- 使用DashScope API将图像转换为视频
- RESTful API接口
- 异步任务处理
- 配置化管理

## 安装

1. 克隆仓库：
```bash
git clone https://github.com/yourusername/salemind_backend_tiny.git
cd salemind_backend_tiny
```

2. 安装依赖：
```bash
go mod tidy
```

3. 配置：
复制`config/config.yaml.example`为`config/config.yaml`并填写相关配置：
```yaml
qwen:
  api_key: "your-qwen-api-key"
  base_url: "https://dashscope.aliyuncs.com/compatible-mode/v1"

liblibai:
  access_key: "your-liblibai-access-key"
  secret_key: "your-liblibai-secret-key"
  api_url: "https://openapi.liblibai.cloud/api/generate/webui/text2img"
  query_url: "https://openapi.liblibai.cloud/api/generate/webui/status"
```

## 运行

```bash
go run cmd/server/main.go
```

服务将在`http://localhost:8080`启动。

## API接口

### 创建图像生成任务

```http
POST /api/image-task
Content-Type: application/json

{
    "pose": "站立",
    "location": "海边",
    "time_of_day": "傍晚",
    "hair_color": "浅绿色",
    "hairstyle": "偏分直短发",
    "top_wear": "白衬衫",
    "bottom_wear": "黑色小裙子",
    "leg_wear": "黑色丝袜"
}
```

使用curl命令：
```bash
curl -X POST http://localhost:8080/api/image-task \
  -H "Content-Type: application/json" \
  -d '{
    "pose": "站立",
    "location": "海边",
    "time_of_day": "傍晚",
    "hair_color": "浅绿色",
    "hairstyle": "偏分直短发",
    "top_wear": "白衬衫",
    "bottom_wear": "黑色小裙子",
    "leg_wear": "黑色丝袜"
  }'
```

响应：
```json
{
    "imageTaskID": "xxx-uuid"
}
```

### 查询图像任务状态

```http
GET /api/image-task/{imageTaskID}/status
```

使用curl命令：
```bash
curl -X GET http://localhost:8080/api/image-task/{imageTaskID}/status
```

响应：
```json
{
    "status": "GENERATING",
    "imgUrl": "https://xxx/xxx.jpg"
}
```

### 创建视频生成任务

```http
POST /api/video-task
Content-Type: application/json

{
    "imgUrl": "https://xxx/xxx.jpg",
    "videoPrompt": "美女轻轻摇摆着，表情魅惑而动人。她微微左右扭动着胯部，看起来很惬意。"
}
```

使用curl命令：
```bash
curl -X POST http://localhost:8080/api/video-task \
  -H "Content-Type: application/json" \
  -d '{
    "imgUrl": "https://xxx/xxx.jpg",
    "videoPrompt": "美女轻轻摇摆着，表情魅惑而动人。她微微左右扭动着胯部，看起来很惬意。"
  }'
```

响应：
```json
{
    "videoTaskID": "xxx-uuid"
}
```

### 查询视频任务状态

```http
GET /api/video-task/{videoTaskID}/status
```

使用curl命令：
```bash
curl -X GET http://localhost:8080/api/video-task/{videoTaskID}/status
```

响应：
```json
{
    "status": "GENERATING",
    "videoUrl": "https://xxx/xxx.mp4"
}
```

## 任务状态说明

- `GENERATING`: 任务正在生成中
- `SUCCEEDED`: 任务成功完成
- `FAILED`: 任务失败
- `TIMEOUT`: 任务超时

## 错误处理

服务会返回适当的HTTP状态码和错误信息：

- 400: 请求参数错误
- 500: 服务器内部错误

## 许可证

MIT
