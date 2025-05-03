# Salemind Backend Tiny

这是一个基于 Golang 和 Gin 框架开发的图像和视频生成服务后端。

## 功能特点

- 图像生成：根据关键词生成高质量图像
- 视频生成：将生成的图像转换为动态视频
- RESTful API：提供简单易用的 HTTP 接口
- 异步任务处理：支持长时间运行的任务状态查询

## 技术栈

- Golang 1.21+
- Gin Web 框架
- HMAC-SHA1 签名认证
- 阿里云通义千问 API
- Liblibai API

## 快速开始

1. 克隆项目
```bash
git clone https://github.com/yourusername/salemind_backend_tiny.git
cd salemind_backend_tiny
```

2. 安装依赖
```bash
go mod download
```

3. 配置
复制 `config/config.yaml.example` 到 `config/config.yaml` 并填写相关配置：
```yaml
liblibai:
  api_url: "https://api.liblibai.com"
  query_url: "https://api.liblibai.com"
  access_key: "your_access_key"
  secret_key: "your_secret_key"

qwen:
  base_url: "https://dashscope.aliyuncs.com/api/v1"
  api_key: "your_api_key"

video:
  api_url: "https://dashscope.aliyuncs.com/api/v1/services/aigc/video-generation/generation"
  task_url: "https://dashscope.aliyuncs.com/api/v1/tasks"
```

4. 运行服务
```bash
go run cmd/server/main.go
```

## API 使用说明

### 1. 创建图像生成任务

```bash
curl -X POST http://43.134.215.212:8081/api/image_task/create -H "Content-Type: application/json" -d '{"keyword":{"pose":"站立","location":"海边","time_of_day":"傍晚","hair_color":"浅绿色","hairstyle":"偏分直短发","top_wear":"白衬衫","bottom_wear":"黑色小裙子","leg_wear":"黑色丝袜"}}'
```

响应示例：
```json
{
  "image_task_id": "a478a7b6926f44179d0f3ed56313cac2"
}
```

### 2. 查询图像任务状态

```bash
curl -X POST http://43.134.215.212:8081/api/image_task/status \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": "a478a7b6926f44179d0f3ed56313cac2"
  }'
```

响应示例：
```json
{
  "status": "completed",
  "image_url": "https://liblibai-tmp-image.liblib.cloud/img/2404cd2fbdbc477180edc62e1ca84946/c83c877f605a5352ae48f807253728d9d5701eee46a8514071a074226108b1da.png"
}
```

### 3. 创建视频生成任务

```bash
curl -X POST http://43.134.215.212:8081/api/video_task/create \
  -H "Content-Type: application/json" \
  -d '{
    "img_url": "https://liblibai-tmp-image.liblib.cloud/img/2404cd2fbdbc477180edc62e1ca84946/c83c877f605a5352ae48f807253728d9d5701eee46a8514071a074226108b1da.png"
  }'
```

响应示例：
```json
{
  "video_task_id": "34ebbf70-aba3-4473-84c2-e3568a59a6e0"
}
```

### 4. 查询视频任务状态

```bash
curl -X POST http://43.134.215.212:8081/api/video_task/status \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": "34ebbf70-aba3-4473-84c2-e3568a59a6e0"
  }'
```

响应示例：
```json
{
  "status": "completed",
  "video_url": "https://dashscope-result-sh.oss-cn-shanghai.aliyuncs.com/1d/26/20250504/8ce2943a/34ebbf70-aba3-4473-84c2-e3568a59a6e0.mp4?Expires=1746393615\u0026OSSAccessKeyId=LTAI5tKPD3TMqf2Lna1fASuh\u0026Signature=ROjPTa7HTXU3FLkaHOmNISMpwtM%3D"
}
```

## 项目结构

```
.
├── cmd/
│   ├── server/        # 服务器入口
│   └── video_generation/  # 视频生成命令行工具
├── config/            # 配置文件
├── model/            # 数据模型
├── pkg/
│   ├── api/          # API 处理器
│   ├── config/       # 配置管理
│   └── services/     # 业务逻辑
└── README.md
```

## 开发说明

1. 代码规范
- 使用 `gofmt` 格式化代码
- 遵循 Go 标准项目布局
- 使用有意义的变量和函数命名

2. 错误处理
- 所有错误都需要被正确处理和记录
- 返回给客户端的错误信息要友好且有意义

3. 测试
- 编写单元测试
- 测试覆盖率要求 > 80%

## 许可证

MIT License
