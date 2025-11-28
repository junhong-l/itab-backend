# 远程同步接口文档

## 认证方式

所有同步接口使用 Header 认证，需要在请求头中携带：

| Header | 说明 |
|--------|------|
| `x-access-key` | 访问密钥 Access Key |
| `x-secret-key` | 密钥 Secret Key |

---

## 1. 获取备份列表

获取当前用户的所有备份数据列表。

### 请求

```
GET /api/sync/list
```

### 请求头

```
x-access-key: AKxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
x-secret-key: SKxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### 响应

#### 成功 (200)

```json
{
    "data": [
        {
            "id": 1,
            "name": "我的导航备份",
            "size": 15360,
            "sync_count": 5,
            "created_at": "2025-11-28T10:00:00Z",
            "updated_at": "2025-11-28T15:30:00Z"
        },
        {
            "id": 2,
            "name": "工作导航",
            "size": 8192,
            "sync_count": 3,
            "created_at": "2025-11-27T09:00:00Z",
            "updated_at": "2025-11-28T12:00:00Z"
        }
    ]
}
```

#### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| id | number | 备份ID |
| name | string | 备份名称（唯一） |
| size | number | 备份大小（字节） |
| sync_count | number | 同步次数 |
| created_at | string | 创建时间 (ISO 8601) |
| updated_at | string | 最后更新时间 (ISO 8601) |

#### 错误响应

```json
// 401 未授权
{ "error": "未提供访问密钥" }

// 401 密钥无效
{ "error": "invalid access key" }

// 401 密钥已过期
{ "error": "access key has expired" }
```

---

## 2. 下载备份数据

根据备份ID下载完整的备份JSON数据。

### 请求

```
GET /api/sync/download/{id}
```

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| id | number | 备份ID |

### 请求头

```
x-access-key: AKxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
x-secret-key: SKxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### 响应

#### 成功 (200)

返回 `Content-Type: application/json`，响应体为备份的完整JSON数据：

```json
{
    "partitions": [
        {
            "id": 1,
            "name": "默认工作区",
            "order": 0,
            "isPrivate": false
        }
    ],
    "folders": [
        {
            "id": 1,
            "name": "常用网站",
            "collapsed": false,
            "partitionId": 1,
            "isPrivate": false
        }
    ],
    "shortcuts": [
        {
            "id": 1,
            "name": "GitHub",
            "url": "https://github.com",
            "icon": "",
            "folderId": 1,
            "partitionId": 1,
            "isPrivate": false,
            "isPinned": true
        }
    ],
    "searchEngines": [
        {
            "id": 1,
            "name": "Google",
            "url": "https://www.google.com/search?q=%s",
            "icon": ""
        }
    ],
    "settings": {
        "bgType": "gradient",
        "gradientColor1": "#667eea",
        "gradientColor2": "#764ba2",
        "gradientAngle": 135,
        "solidColor": "#1a1a2e",
        "bgImage": "",
        "iconSize": 57,
        "folderSize": 65,
        "iconGap": 15,
        "folderGap": 20,
        "iconRadius": 11,
        "searchRadius": 16,
        "btnRadius": 12
    }
}
```

#### 错误响应

```json
// 400 参数错误
{ "error": "无效的备份ID" }

// 401 未授权
{ "error": "未提供访问密钥" }

// 404 不存在
{ "error": "备份不存在" }
```

---

## 3. 上传备份数据

上传或更新备份数据。如果同名备份已存在则更新，否则创建新备份。

### 请求

```
POST /api/sync/upload
```

### 请求头

```
Content-Type: application/json
x-access-key: AKxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
x-secret-key: SKxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### 请求体

```json
{
    "name": "我的导航备份",
    "data": {
        "partitions": [...],
        "folders": [...],
        "shortcuts": [...],
        "searchEngines": [...],
        "settings": {...}
    }
}
```

#### 请求参数说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 备份名称，用于标识备份（同用户下唯一） |
| data | object | 是 | 备份数据对象 |

#### data 对象结构

| 字段 | 类型 | 说明 |
|------|------|------|
| partitions | array | 工作区/分区列表 |
| folders | array | 文件夹列表 |
| shortcuts | array | 书签列表 |
| searchEngines | array | 搜索引擎列表 |
| settings | object | 外观设置 |

### 响应

#### 创建成功 (200)

```json
{
    "message": "备份创建成功",
    "backup_id": 1
}
```

#### 更新成功 (200)

```json
{
    "message": "备份更新成功",
    "backup_id": 1
}
```

#### 错误响应

```json
// 400 参数错误
{ "error": "参数错误: Key: 'SyncUploadRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag" }

// 401 未授权
{ "error": "未提供访问密钥" }

// 500 服务器错误
{ "error": "创建备份失败" }
```

---

## 完整示例

### cURL 示例

#### 获取备份列表

```bash
curl -X GET "http://localhost:8445/api/sync/list" \
  -H "x-access-key: AK1234567890abcdef1234567890ab" \
  -H "x-secret-key: SK1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd"
```

#### 下载备份

```bash
curl -X GET "http://localhost:8445/api/sync/download/1" \
  -H "x-access-key: AK1234567890abcdef1234567890ab" \
  -H "x-secret-key: SK1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd"
```

#### 上传备份

```bash
curl -X POST "http://localhost:8445/api/sync/upload" \
  -H "Content-Type: application/json" \
  -H "x-access-key: AK1234567890abcdef1234567890ab" \
  -H "x-secret-key: SK1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd" \
  -d '{
    "name": "我的导航备份",
    "data": {
      "partitions": [{"id": 1, "name": "默认", "order": 0, "isPrivate": false}],
      "folders": [],
      "shortcuts": [{"id": 1, "name": "Google", "url": "https://google.com", "icon": "", "folderId": null, "partitionId": 1, "isPrivate": false, "isPinned": false}],
      "searchEngines": [{"id": 1, "name": "Google", "url": "https://www.google.com/search?q=%s", "icon": ""}],
      "settings": {"bgType": "gradient", "gradientColor1": "#667eea", "gradientColor2": "#764ba2", "gradientAngle": 135, "solidColor": "#1a1a2e", "bgImage": "", "iconSize": 57, "folderSize": 65, "iconGap": 15, "folderGap": 20, "iconRadius": 11, "searchRadius": 16, "btnRadius": 12}
    }
  }'
```

### JavaScript 示例

```javascript
const API_BASE = 'http://localhost:8445';
const ACCESS_KEY = 'AK1234567890abcdef1234567890ab';
const SECRET_KEY = 'SK1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd';

const headers = {
  'Content-Type': 'application/json',
  'x-access-key': ACCESS_KEY,
  'x-secret-key': SECRET_KEY
};

// 获取备份列表
async function getBackupList() {
  const response = await fetch(`${API_BASE}/api/sync/list`, { headers });
  return response.json();
}

// 下载备份
async function downloadBackup(id) {
  const response = await fetch(`${API_BASE}/api/sync/download/${id}`, { headers });
  return response.json();
}

// 上传备份
async function uploadBackup(name, data) {
  const response = await fetch(`${API_BASE}/api/sync/upload`, {
    method: 'POST',
    headers,
    body: JSON.stringify({ name, data })
  });
  return response.json();
}
```

---

## 错误码说明

| HTTP 状态码 | 说明 |
|-------------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 认证失败（密钥无效或已过期） |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
