# API接口参考手册

> **学习目标**：完成本手册学习后，你将能够通过API集成火宝短剧平台的全部核心功能，实现剧本创建、角色管理、分镜生成、视频生成等功能的程序化调用。
> 
> **前置知识**：建议具备基本的RESTful API调用经验，熟悉HTTP协议和JSON数据格式。建议先完成《用户操作手册》的学习，理解平台的功能和数据模型。
> 
> **难度等级**：⭐⭐⭐（进阶分析）
> 
> **预计学习时间**：4-8小时

---

## 第一章：API概述与设计原则

### 1.1 API架构设计

火宝短剧平台采用**RESTful API**架构设计，所有接口遵循HTTP协议的语义规范，使用JSON作为数据交换格式。这种设计选择基于以下考量：RESTful架构是当前Web服务的主流标准，具有良好的可理解性、可扩展性和工具支持；JSON格式被所有主流编程语言原生支持，便于快速集成开发。理解这个基础设计有助于更好地使用API——当你遇到问题时，可以参考RESTful的设计原则进行推理和排查。

API的**基础URL**配置取决于部署方式，默认配置如下：

| 部署环境 | 基础URL |
|----------|---------|
| 本地开发 | `http://localhost:5678/api/v1` |
| Docker部署 | `http://localhost:5678/api/v1` |
| 生产环境 | `https://your-domain.com/api/v1` |

**版本控制**：API采用URL路径版本控制，当前版本为`v1`。所有接口路径都以`/api/v1/`开头。这种版本控制方式的优势是清晰直观，不同版本的API可以同时运行，便于渐进式迁移。版本升级时会尽量保持向后兼容，但重大变更会发布新版本。

**认证机制**：API使用**Bearer Token**认证方式。客户端在请求头中需要携带`Authorization: Bearer <token>`字段。Token通过登录接口获取，有效期为24小时（可配置）。Token丢失或过期会导致401未授权错误，需要重新获取。所有涉及项目数据的接口都需要认证，仅公开资源（如静态文件访问）可以匿名访问。

### 1.2 请求与响应规范

**请求格式规范**：

所有POST和PUT请求的Body必须为有效的JSON对象，并设置正确的Content-Type头：`Content-Type: application/json`。请求体中的字段命名采用**snake_case**风格（全部小写，以下划线分隔），与数据库字段命名保持一致。日期时间字段统一使用**ISO 8601格式**（如`2026-02-05T12:00:00Z`），时区统一使用UTC。

```json
// 请求示例
{
  "name": "我的短剧项目",
  "description": "这是一个测试项目",
  "created_at": "2026-02-05T12:00:00Z",
  "settings": {
    "resolution": "1080p",
    "frame_rate": 24
  }
}
```

**响应格式规范**：

API的响应统一封装在标准的响应结构中，区分成功和错误两种情况。成功的响应包含`data`字段存放业务数据；错误的响应包含`error`字段存放错误详情。

```json
// 成功响应
{
  "code": 0,
  "message": "success",
  "data": {
    // 业务数据
  }
}

// 错误响应
{
  "code": 1001,
  "message": "validation error",
  "error": {
    "field": "name",
    "reason": "field is required"
  }
}
```

**分页查询规范**：

返回列表的接口统一支持分页参数，分页信息包含在响应的`pagination`字段中。分页参数通过Query String传递：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | integer | 1 | 当前页码，从1开始 |
| per_page | integer | 20 | 每页数量，最大100 |
| order_by | string | created_at | 排序字段 |
| order | string | desc | 排序方向，asc或desc |

```json
// 分页响应结构
{
  "code": 0,
  "message": "success",
  "data": [
    // 数据数组
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### 1.3 错误处理与状态码

HTTP状态码和业务错误码共同描述请求的处理结果。正确理解和处理这些错误码是构建健壮应用的基础。

**HTTP状态码**：

| 状态码 | 含义 | 处理建议 |
|--------|------|----------|
| 200 | 请求成功 | 解析响应数据 |
| 201 | 资源创建成功 | 获取新资源的ID或URI |
| 400 | 请求参数错误 | 检查请求参数格式和内容 |
| 401 | 未授权 | 检查认证Token是否有效 |
| 403 | 禁止访问 | 检查权限配置 |
| 404 | 资源不存在 | 检查资源ID是否正确 |
| 422 | 业务验证失败 | 查看具体错误信息 |
| 429 | 请求过于频繁 | 实现退避重试机制 |
| 500 | 服务器内部错误 | 记录错误日志，联系支持 |
| 503 | 服务不可用 | 检查服务状态 |

**业务错误码**：

| 错误码 | 含义 | 说明 |
|--------|------|------|
| 1000 | 通用错误 | 未分类的错误 |
| 1001 | 参数验证错误 | 请求参数不合法 |
| 1002 | 资源不存在 | 指定的资源ID不存在 |
| 1003 | 资源已存在 | 创建的资源已存在 |
| 1004 | 权限不足 | 无权访问该资源 |
| 1005 | 操作频率限制 | 超过API调用频率限制 |
| 2001 | AI服务错误 | AI服务返回错误 |
| 2002 | AI服务超时 | AI服务响应超时 |
| 2003 | AI配额不足 | AI服务调用次数超限 |

---

## 第二章：剧集管理API

> **重要说明**：火宝短剧使用「剧集(Drama)」作为项目的核心概念，与传统「项目(Project)」术语不同。Drama包含剧本、角色、分镜等所有创作内容。

### 2.1 剧集CRUD操作

```
POST /api/v1/dramas
```

创建剧集时需要提供剧集的基本信息。剧集名称在同一用户下必须唯一，重复创建会返回409冲突错误。建议在创建剧集时同时设置合适的描述，便于后续管理和识别。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 剧集名称，1-100字符 |
| description | string | 否 | 剧集描述 |
| settings | object | 否 | 剧集配置 |

**请求示例**：

```bash
curl -X POST "http://localhost:5678/api/v1/dramas" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "我的第一个短剧",
    "description": "AI辅助创作的青春校园短剧",
    "settings": {
      "resolution": "1080p",
      "frame_rate": 24
    }
  }'
```

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "proj_abc123",
    "name": "我的第一个短剧",
    "description": "AI辅助创作的青春校园短剧",
    "settings": {
      "resolution": "1080p",
      "frame_rate": 24
    },
    "created_at": "2026-02-05T12:00:00Z",
    "updated_at": "2026-02-05T12:00:00Z"
  }
}
```

**获取剧集列表**：

```
GET /api/v1/dramas
```

获取当前用户的剧集列表，支持分页和排序。列表按照更新时间倒序排列，最新更新的剧集排在最前面。

**查询参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| page | integer | 页码，默认1 |
| per_page | integer | 每页数量，默认20 |
| keyword | string | 搜索关键词，匹配剧集名称 |

**获取单个剧集**：

```
GET /api/v1/dramas/{drama_id}
```

获取指定剧集的详细信息，包括剧集配置、创建时间、更新时间等。

**更新剧集**：

```
PUT /api/v1/dramas/{drama_id}
```

更新剧集的名称、描述或配置。更新时使用**全量更新**策略，即请求中需要包含所有需要保留的字段，未包含的字段会被清空（除系统字段外）。如果需要部分更新，请使用PATCH方法。

**删除剧集**：

```
DELETE /api/v1/dramas/{drama_id}
```

删除剧集及其关联的所有资源（剧本、角色、分镜、视频）。此操作不可逆，删除前请确认。

### 2.2 剧本管理接口

剧本是剧集的核心内容，承载了短剧的文字信息。剧本API提供了剧本的创建、读取、更新、删除以及AI生成等功能。

**创建剧本**：

```
POST /api/v1/dramas/{drama_id}/scripts
```

剧本可以手动输入，也可以通过AI生成。创建剧本时可以同时设置剧本的结构化信息，如场景列表、角色列表等。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 剧本标题 |
| content | string | 是 | 剧本正文内容 |
| scenes | array | 否 | 场景列表 |
| metadata | object | 否 | 元数据 |

**响应字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 剧本唯一标识符 |
| drama_id | string | 所属剧集ID |
| title | string | 剧本标题 |
| content | string | 剧本正文 |
| scene_count | integer | 场景数量 |
| word_count | integer | 字数统计 |
| status | string | 剧本状态 |

**AI辅助生成剧本**：

```
POST /api/v1/dramas/{drama_id}/scripts/generate
```

这是平台的核心功能之一，通过AI自动生成剧本。请求时需要提供创作意图描述，AI会根据描述生成完整的剧本结构。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| prompt | string | 是 | 创作意图描述 |
| style | string | 否 | 生成风格 |
| length | string | 否 | 生成长度 |
| language | string | 否 | 语言 |

**请求示例**：

```bash
curl -X POST "http://localhost:5678/api/v1/dramas/proj_abc123/scripts/generate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "prompt": "生成一个关于青春成长的短剧，讲述高中生小明通过努力实现梦想的故事",
    "style": "温馨励志",
    "length": "medium"
  }'
```

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "script_xyz789",
    "drama_id": "proj_abc123",
    "title": "青春的梦想",
    "content": "场景一：学校操场（内，日）...",
    "scenes": [
      {
        "id": "scene_001",
        "location": "学校操场",
        "time": "日",
        "description": "小明站在跑道上..."
      }
    ],
    "status": "completed"
  }
}
```

**获取剧本详情**：

```
GET /api/v1/dramas/{drama_id}/scripts/{script_id}
```

获取剧本的详细信息，包括剧本内容、场景列表、统计数据等。

**更新剧本**：

```
PUT /api/v1/dramas/{drama_id}/scripts/{script_id}
```

更新剧本内容或元数据。支持全量更新和部分更新两种模式。

### 2.3 剧本版本管理

剧本支持版本历史追踪，便于回溯修改和比较差异。

**获取版本列表**：

```
GET /api/v1/dramas/{drama_id}/scripts/{script_id}/versions
```

列出剧本的所有历史版本，每次保存会自动创建一个新版本。

**响应字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 版本ID |
| script_id | string | 剧本ID |
| content | string | 版本内容快照 |
| created_at | string | 创建时间 |
| message | string | 版本说明 |

**恢复版本**：

```
POST /api/v1/dramas/{drama_id}/scripts/{script_id}/versions/{version_id}/restore
```

将剧本恢复到指定版本的内容。恢复操作会创建一个新的版本快照，不会删除历史版本。

---

## 第三章：角色管理API

### 3.1 角色CRUD操作

角色API提供了角色的创建、读取、更新、删除等基本操作，以及AI形象生成和图片上传等特色功能。

**创建角色**：

```
POST /api/v1/dramas/{drama_id}/characters
```

创建角色时需要提供角色的基本信息和形象描述。角色名称在项目内必须唯一。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 角色名称 |
| description | string | 是 | 角色描述，用于AI生成形象 |
| gender | string | 否 | 性别 |
| age_range | string | 否 | 年龄范围 |
| image_prompt | string | 否 | 图像生成提示词 |

**请求示例**：

```bash
curl -X POST "http://localhost:5678/api/v1/dramas/proj_abc123/characters" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "小明",
    "description": "17岁高中生，性格内向但善良，学习刻苦",
    "gender": "male",
    "age_range": "15-20",
    "image_prompt": "17岁男生，短发，戴眼镜，穿着整洁的校服"
  }'
```

**AI生成角色形象**：

```
POST /api/v1/dramas/{drama_id}/characters/{character_id}/generate-image
```

使用AI根据角色的描述生成形象图片。生成是异步过程，返回任务ID供后续查询状态。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| prompt | string | 否 | 图像生成提示词，不提供则使用角色描述 |
| style | string | 否 | 生成风格 |
| count | integer | 否 | 生成数量，默认1 |

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "task_img_001",
    "status": "pending",
    "estimated_time": 60
  }
}
```

**查询任务状态**：

```
GET /api/v1/tasks/{task_id}
```

查询异步任务的处理状态和结果。

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "task_img_001",
    "status": "completed",
    "result": {
      "images": [
        {
          "id": "img_001",
          "url": "/static/images/characters/img_001.jpg",
          "thumbnail_url": "/static/images/characters/thumb_img_001.jpg"
        }
      ]
    }
  }
}
```

**上传角色图片**：

```
POST /api/v1/dramas/{drama_id}/characters/{character_id}/upload-image
```

上传本地图片作为角色形象。请求需要使用`multipart/form-data`格式。

**请求参数**（表单字段）：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| image | file | 是 | 图片文件 |
| make_default | boolean | 否 | 是否设为默认形象 |

### 3.2 角色批量操作

当项目中有大量角色时，批量操作可以显著提高效率。

**批量获取角色**：

```
GET /api/v1/dramas/{drama_id}/characters/batch
```

通过角色ID列表批量获取角色详情。

**请求参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| ids | string | 角色ID列表，逗号分隔 |

**批量删除角色**：

```
DELETE /api/v1/dramas/{drama_id}/characters/batch
```

批量删除指定角色，被角色关联的分镜不会自动删除。

---

## 第四章：分镜制作API

### 4.1 分镜生成接口

分镜是连接剧本和视频的桥梁。分镜API提供了从剧本自动生成分镜、手动创建分镜、分镜编辑调整等功能。

**AI自动生成分镜**：

```
POST /api/v1/dramas/{drama_id}/storyboards/generate
```

这是核心功能之一，自动解析剧本内容，生成对应的分镜序列。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| script_id | string | 是 | 剧本ID |
| options | object | 否 | 生成选项 |

**options字段说明**：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
|镜头类型|string|auto|自动选择镜头|
|图像风格|string|realistic|图像生成风格|
|count_per_scene|integer|1|每个场景生成的分镜数量|

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "task_story_001",
    "status": "pending",
    "estimated_time": 120,
    "scene_count": 10
  }
}
```

**查询分镜生成任务状态**：

```
GET /api/v1/tasks/{task_id}
```

查询分镜生成的处理状态。完成后返回生成的分镜列表。

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "task_story_001",
    "status": "completed",
    "result": {
      "storyboards": [
        {
          "id": "sb_001",
          "scene_id": "scene_001",
          "scene_number": 1,
          "shot_type": "medium_shot",
          "description": "小明走进教室，阳光透过窗户...",
          "image_url": "/static/images/storyboards/sb_001.jpg",
          "status": "completed"
        }
      ],
      "total_count": 10,
      "success_count": 9,
      "failed_count": 1
    }
  }
}
```

### 4.2 分镜编辑接口

生成的分镜可以进一步编辑调整，以达到理想的视觉效果。

**更新分镜**：

```
PUT /api/v1/dramas/{drama_id}/storyboards/{storyboard_id}
```

更新分镜的描述、镜头类型、构图等信息。

**请求参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| description | string | 画面描述 |
| shot_type | string | 镜头类型 |
| composition | object | 构图信息 |
| character_ids | array | 关联的角色ID列表 |

**镜头类型枚举值**：

| 值 | 说明 |
|------|------|
| close_up | 特写 |
| medium_close_up | 近景 |
| medium | 中景 |
| wide | 全景 |
| extreme_wide | 大远景 |
| high_angle | 俯拍 |
| low_angle | 仰拍 |

**重新生成分镜图像**：

```
POST /api/v1/dramas/{drama_id}/storyboards/{storyboard_id}/regenerate-image
```

使用新的参数重新生成分镜图像。适用于对当前图像不满意的场景。

**请求参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| prompt | string | 图像生成提示词 |
| image_style | string | 图像风格 |
| fix_face | boolean | 是否修复面部 |

### 4.3 分镜排序与关联

分镜的顺序决定了最终视频的叙事逻辑。正确的排序和关联是高质量产出的基础。

**调整分镜顺序**：

```
PUT /api/v1/dramas/{drama_id}/storyboards/reorder
```

批量调整分镜的顺序。请求中需要提供完整的分镜ID列表及其新位置。

**请求示例**：

```json
{
  "order": [
    {"id": "sb_003", "position": 1},
    {"id": "sb_001", "position": 2},
    {"id": "sb_002", "position": 3}
  ]
}
```

**关联剧本场景**：

```
PUT /api/v1/dramas/{drama_id}/storyboards/{storyboard_id}/link-scene
```

将分镜与剧本中的特定场景关联。关联后，分镜会继承场景的信息（如出现角色、位置、时间等）。

---

## 第五章：视频生成API

### 5.1 视频生成任务

视频生成是计算密集型操作，采用异步任务模式处理。理解任务的生命周期是使用视频API的关键。

**创建视频生成任务**：

```
POST /api/v1/dramas/{drama_id}/videos/generate
```

创建视频生成任务，可以选择单个分镜或批量生成分镜视频。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| storyboard_ids | array | 是 | 要生成视频的分镜ID列表 |
| settings | object | 是 | 视频生成设置 |

**settings字段说明**：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| resolution | string | 1080p | 视频分辨率 |
| frame_rate | integer | 24 | 帧率 |
| duration | integer | 5 | 单个视频时长（秒） |
| camera_movement | string | auto | 运镜类型 |
| transition | string | fade | 转场效果 |

**请求示例**：

```bash
curl -X POST "http://localhost:5678/api/v1/dramas/proj_abc123/videos/generate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "storyboard_ids": ["sb_001", "sb_002", "sb_003"],
    "settings": {
      "resolution": "1080p",
      "frame_rate": 24,
      "duration": 5,
      "camera_movement": "push_in",
      "transition": "fade"
    }
  }'
```

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "task_video_001",
    "status": "pending",
    "estimated_time": 300,
    "storyboard_count": 3
  }
}
```

**查询视频生成进度**：

```
GET /api/v1/tasks/{task_id}
```

查询视频生成任务的实时进度。

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "task_video_001",
    "status": "processing",
    "progress": {
      "total": 3,
      "completed": 1,
      "processing": 1,
      "pending": 1,
      "percentage": 33
    },
    "current_item": {
      "storyboard_id": "sb_002",
      "status": "processing"
    },
    "estimated_remaining": 200
  }
}
```

**任务完成状态**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "task_video_001",
    "status": "completed",
    "result": {
      "videos": [
        {
          "id": "video_001",
          "storyboard_id": "sb_001",
          "url": "/static/videos/video_001.mp4",
          "thumbnail_url": "/static/videos/thumb_video_001.jpg",
          "duration": 5.2,
          "file_size": 10485760
        }
      ],
      "total_count": 3,
      "success_count": 3,
      "failed_count": 0
    }
  }
}
```

### 5.2 视频导出与下载

视频生成完成后，需要导出到本地或进行下一步处理。

**获取视频文件信息**：

```
GET /api/v1/dramas/{drama_id}/videos/{video_id}
```

获取单个视频的详细信息，包括文件URL、时长、分辨率等。

**导出视频**：

```
GET /api/v1/dramas/{drama_id}/videos/{video_id}/export
```

导出视频文件到本地。返回文件的下载URL或直接返回文件流。

**批量导出**：

```
POST /api/v1/dramas/{drama_id}/videos/export-batch
```

批量导出多个视频，支持打包成ZIP文件。

### 5.3 视频拼接与合成

支持将多个分镜视频合成为完整的长视频。

**创建拼接任务**：

```
POST /api/v1/dramas/{drama_id}/videos/concatenate
```

将多个视频片段拼接成一个完整视频。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| video_ids | array | 是 | 要拼接的视频ID列表 |
| output_settings | object | 否 | 输出设置 |

**output_settings字段说明**：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| format | string | mp4 | 输出格式 |
| resolution | string | same_as_input | 输出分辨率 |
| transition | string | none | 片段间转场 |

---

## 第六章：资源管理API

### 6.1 资产管理接口

资源是项目中的各类文件资产，包括图片、视频、音频等。

**获取资源列表**：

```
GET /api/v1/dramas/{drama_id}/assets
```

获取项目的资源列表，支持分类筛选和分页。

**查询参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| type | string | 资源类型筛选 |
| tags | string | 标签筛选 |
| keyword | string | 关键词搜索 |

**资源类型枚举**：

| 值 | 说明 |
|------|------|
| image | 图片资源 |
| video | 视频资源 |
| audio | 音频资源 |
| document | 文档资源 |

### 6.2 文件上传与下载

文件上传使用multipart/form-data格式，支持大文件分片上传。

**上传文件**：

```
POST /api/v1/dramas/{drama_id}/assets/upload
```

上传文件到项目资源库。

**请求参数**（表单字段）：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 要上传的文件 |
| type | string | 否 | 资源类型 |
| tags | string | 否 | 标签列表，逗号分隔 |

---

## 第七章：认证与用户API

### 7.1 认证接口

用户认证采用JWT（JSON Web Token）机制。

**用户登录**：

```
POST /api/v1/auth/login
```

使用用户名和密码获取访问令牌。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "token_type": "Bearer"
  }
}
```

**刷新令牌**：

```
POST /api/v1/auth/refresh
```

使用刷新令牌获取新的访问令牌。

### 7.2 用户信息接口

**获取当前用户信息**：

```
GET /api/v1/users/me
```

获取当前登录用户的详细信息。

---

## 第八章：错误处理与最佳实践

### 8.1 错误处理模式

构建健壮的API客户端需要正确处理各种错误情况。以下是推荐的错误处理模式。

**重试机制**：

对于临时性错误（如网络超时、服务暂时不可用），实现指数退避重试机制：

```python
import time
import requests

def retry_with_backoff(func, max_retries=3, base_delay=1):
    """带指数退避的重试装饰器"""
    for attempt in range(max_retries):
        try:
            return func()
        except (requests.Timeout, requests.ConnectionError) as e:
            if attempt == max_retries - 1:
                raise
            delay = base_delay * (2 ** attempt)
            time.sleep(delay)
```

**错误分类处理**：

```python
def handle_api_error(response):
    """根据错误类型采取相应处理"""
    if response.status_code == 401:
        # Token过期，尝试刷新或重新登录
        refresh_token_and_retry()
    elif response.status_code == 429:
        # 超过频率限制，等待后重试
        wait_and_retry()
    elif response.status_code >= 500:
        # 服务器错误，记录并告警
        log_and_alert()
    else:
        # 客户端错误，检查请求参数
        handle_client_error()
```

### 8.2 性能优化建议

**并发请求**：

对于独立的API调用（如批量获取角色信息），使用并发请求可以显著提高效率：

```python
import concurrent.futures

def batch_fetch_characters(drama_id, character_ids, max_workers=5):
    """并发获取角色信息"""
    def fetch_one(character_id):
        url = f"/api/v1/dramas/{drama_id}/characters/{character_id}"
        return requests.get(url, headers=headers)
    
    with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = [executor.submit(fetch_one, cid) for cid in character_ids]
        results = [f.result() for f in concurrent.futures.as_completed(futures)]
    
    return results
```

**分页遍历**：

对于大型列表，使用分页参数进行高效遍历：

```python
def fetch_all_items(base_url, params=None):
    """获取所有分页数据"""
    all_items = []
    page = 1
    per_page = 100
    
    while True:
        response = requests.get(
            base_url,
            params={**(params or {}), "page": page, "per_page": per_page}
        )
        data = response.json()["data"]
        all_items.extend(data["items"])
        
        if page >= data["pagination"]["total_pages"]:
            break
        page += 1
    
    return all_items
```

### 8.3 安全最佳实践

**敏感信息处理**：

- API密钥和密码不应硬编码在代码中，应使用环境变量或配置管理服务
- 传输过程中始终使用HTTPS
- 日志中不应记录敏感信息

**访问控制**：

```python
# 验证响应数据的完整性
def validate_response(response, schema):
    """使用JSON Schema验证响应结构"""
    try:
        validate(instance=response.json(), schema=schema)
        return True
    except ValidationError as e:
        log_error(f"Response validation failed: {e}")
        return False
```

---

## 第九章：SDK与集成示例

### 9.1 Python SDK使用示例

以下是一个完整的Python SDK封装示例：

```python
import requests
import json
from typing import Optional, Dict, List
from dataclasses import dataclass

@dataclass
class HuobaoClient:
    """火宝短剧API客户端"""
    base_url: str
    access_token: str
    
    def __init__(self, base_url: str, access_token: str):
        self.base_url = base_url.rstrip('/')
        self.access_token = access_token
        self.session = requests.Session()
        self.session.headers.update({
            "Authorization": f"Bearer {access_token}",
            "Content-Type": "application/json"
        })
    
    def _request(self, method: str, endpoint: str, **kwargs) -> dict:
        """统一的请求处理方法"""
        url = f"{self.base_url}/api/v1{endpoint}"
        response = self.session.request(method, url, **kwargs)
        response.raise_for_status()
        return response.json()
    
    def get_dramas(self, page: int = 1, per_page: int = 20) -> dict:
        """获取剧集列表"""
        return self._request("GET", "/dramas", params={
            "page": page, "per_page": per_page
        })
    
    def create_drama(self, name: str, description: str = "") -> dict:
        """创建剧集"""
        return self._request("POST", "/dramas", json={
            "name": name,
            "description": description
        })
    
    def generate_script(self, drama_id: str, prompt: str) -> dict:
        """AI生成剧本"""
        return self._request("POST", f"/dramas/{drama_id}/scripts/generate", json={
            "prompt": prompt
        })
    
    def get_task_status(self, task_id: str) -> dict:
        """查询任务状态"""
        return self._request("GET", f"/tasks/{task_id}")
```

### 9.2 完整工作流示例

以下示例展示了一个完整的项目创建到视频生成的工作流程：

```python
def create_drama_workflow(client: HuobaoClient, prompt: str) -> dict:
    """完整的短剧创建流程"""
    
    # 1. 创建剧集
    drama = client.create_drama("AI生成短剧项目")
    drama_id = drama["data"]["id"]
    print(f"剧集创建成功: {drama_id}")
    
    # 2. AI生成剧本
    script_task = client.generate_script(drama_id, prompt)
    task_id = script_task["data"]["task_id"]
    
    # 等待剧本生成完成
    script = wait_for_task(client, task_id, interval=5)
    script_id = script["data"]["result"]["id"]
    print(f"剧本生成成功: {script_id}")
    
    # 3. 生成分镜
    storyboard_task = client.generate_storyboards(drama_id, script_id)
    task_id = storyboard_task["data"]["task_id"]
    
    storyboards = wait_for_task(client, task_id, interval=10)
    storyboard_ids = [sb["id"] for sb in storyboards["data"]["result"]["storyboards"]]
    print(f"分镜生成成功，共{len(storyboard_ids)}个")
    
    # 4. 生成视频
    video_task = client.generate_videos(drama_id, storyboard_ids)
    task_id = video_task["data"]["task_id"]
    
    videos = wait_for_task(client, task_id, interval=30)
    print(f"视频生成成功，共{len(videos['data']['result']['videos'])}个")
    
    return {
        "drama_id": drama_id,
        "script_id": script_id,
        "storyboard_ids": storyboard_ids,
        "videos": videos["data"]["result"]["videos"]
    }

def wait_for_task(client: HuobaoClient, task_id: str, interval: int = 5) -> dict:
    """等待任务完成"""
    import time
    while True:
        result = client.get_task_status(task_id)
        status = result["data"]["status"]
        
        if status == "completed":
            return result
        elif status == "failed":
            raise Exception(f"任务失败: {result['data'].get('error')}")
        else:
            print(f"任务处理中: {status}...")
            time.sleep(interval)
```

---

## 附录A：API完整端点列表

| HTTP方法 | 端点 | 功能 |
|----------|------|------|
| POST | /api/v1/auth/login | 用户登录 |
| POST | /api/v1/auth/refresh | 刷新令牌 |
| GET | /api/v1/dramas | 获取项目列表 |
| POST | /api/v1/dramas | 创建项目 |
| GET | /api/v1/dramas/{id} | 获取项目详情 |
| PUT | /api/v1/dramas/{id} | 更新项目 |
| DELETE | /api/v1/dramas/{id} | 删除项目 |
| POST | /api/v1/dramas/{id}/scripts | 创建剧本 |
| POST | /api/v1/dramas/{id}/scripts/generate | AI生成剧本 |
| GET | /api/v1/dramas/{id}/scripts/{sid} | 获取剧本详情 |
| POST | /api/v1/dramas/{id}/characters | 创建角色 |
| POST | /api/v1/dramas/{id}/characters/{cid}/generate-image | AI生成角色形象 |
| POST | /api/v1/dramas/{id}/storyboards/generate | AI生成分镜 |
| POST | /api/v1/dramas/{id}/videos/generate | 视频生成 |
| GET | /api/v1/tasks/{id} | 查询任务状态 |

---

## 附录B：错误码完整列表

| 错误码 | 说明 |
|--------|------|
| 1000 | 通用错误 |
| 1001 | 参数验证错误 |
| 1002 | 资源不存在 |
| 1003 | 资源已存在 |
| 1004 | 权限不足 |
| 1005 | 操作频率限制 |
| 1006 | 认证令牌无效 |
| 1007 | 认证令牌过期 |
| 2001 | AI服务错误 |
| 2002 | AI服务超时 |
| 2003 | AI配额不足 |
| 2004 | AI内容过滤 |
| 3001 | 文件上传失败 |
| 3002 | 文件格式不支持 |
| 3003 | 文件大小超限 |
| 4001 | 数据库错误 |
| 4002 | 缓存错误 |

---

## 相关文档链接

| 文档名称 | 说明 | 难度 |
|----------|------|------|
| [快速入门指南](quick-start.md) | 15分钟快速上手 | ⭐ |
| [用户手册](user-manual.md) | 平台功能的完整操作指南 | ⭐⭐ |
| [架构设计文档](architecture.md) | 系统架构深度解析 | ⭐⭐⭐⭐ |
| [部署运维指南](deployment.md) | 生产环境部署和运维指南 | ⭐⭐ |
| [故障排查文档](troubleshooting.md) | 常见问题和解决方案 | ⭐⭐ |

---

## 版本信息

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0.0 | 2026-02-05 | 初始版本，包含完整API参考 |
