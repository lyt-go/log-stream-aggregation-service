# 日志聚合系统 (Log Aggregation)

纯 Go 标准库实现的日志聚合后端服务，零第三方依赖，可编译、可测试、可运行。

## 运行说明

```bash
cd origin/
go run ./cmd/server
```

默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 修改。

```bash
PORT=9090 go run ./cmd/server
```

## 项目结构

```
origin/
├── cmd/server/main.go
├── internal/
│   ├── app/app.go
│   ├── config/config.go
│   ├── model/
│   │   ├── errors.go
│   │   ├── stream.go
│   │   ├── collector.go
│   │   ├── log_entry.go
│   │   ├── search_record.go
│   │   └── alert_rule.go
│   ├── store/
│   │   ├── store.go
│   │   ├── memory.go
│   │   ├── stream_store.go
│   │   ├── collector_store.go
│   │   ├── log_entry_store.go
│   │   ├── search_record_store.go
│   │   └── alert_rule_store.go
│   ├── service/
│   │   ├── service.go
│   │   ├── stream_service.go
│   │   ├── collector_service.go
│   │   ├── log_entry_service.go
│   │   ├── search_record_service.go
│   │   ├── alert_rule_service.go
│   │   └── stats_service.go
│   └── handler/
│       ├── server.go
│       ├── stream_handler.go
│       ├── collector_handler.go
│       ├── log_entry_handler.go
│       ├── search_record_handler.go
│       ├── alert_rule_handler.go
│       └── stats_handler.go
└── pkg/
    ├── httpx/httpx.go
    ├── idgen/idgen.go
    └── logger/logger.go
```

## API 列表

### Stream 日志流

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/streams | 创建日志流 |
| GET | /api/streams | 列表（支持 status/format/keyword 筛选，分页） |
| GET | /api/streams/{id} | 获取详情 |
| PUT | /api/streams/{id} | 更新 |
| DELETE | /api/streams/{id} | 删除 |
| POST | /api/streams/{id}/pause | 暂停（active->paused） |
| POST | /api/streams/{id}/resume | 恢复（paused->active） |

### Collector 采集器

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/collectors | 创建采集器（校验 Stream 存在） |
| GET | /api/collectors | 列表（支持 status/stream_id/keyword 筛选，分页） |
| GET | /api/collectors/{id} | 获取详情 |
| PUT | /api/collectors/{id} | 更新 |
| DELETE | /api/collectors/{id} | 删除 |
| POST | /api/collectors/{id}/stop | 停止（active->stopped） |
| POST | /api/collectors/{id}/start | 启动（stopped->active） |

### LogEntry 日志条目

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/log-entries | 创建日志条目（校验 Stream 存在） |
| GET | /api/log-entries | 列表（支持 stream_id/level/keyword/tag 筛选，分页） |
| GET | /api/log-entries/{id} | 获取详情 |
| PUT | /api/log-entries/{id} | 更新 |
| DELETE | /api/log-entries/{id} | 删除 |
| POST | /api/log-entries/batch-delete | 批量删除 |

### SearchRecord 检索记录

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/search-records | 创建检索记录 |
| GET | /api/search-records | 列表（支持 keyword 筛选，分页） |
| GET | /api/search-records/{id} | 获取详情 |
| PUT | /api/search-records/{id} | 更新 |
| DELETE | /api/search-records/{id} | 删除 |

### AlertRule 告警规则

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/alert-rules | 创建告警规则（校验 Stream 存在） |
| GET | /api/alert-rules | 列表（支持 status/stream_id/level/keyword 筛选，分页） |
| GET | /api/alert-rules/{id} | 获取详情 |
| PUT | /api/alert-rules/{id} | 更新 |
| DELETE | /api/alert-rules/{id} | 删除 |
| POST | /api/alert-rules/{id}/disable | 禁用（active->disabled） |
| POST | /api/alert-rules/{id}/enable | 启用（disabled->active） |
| POST | /api/alert-rules/batch-status | 批量更新状态 |

### Stats 统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/stats/by-stream | 按流统计日志条数与各级别数量 |
| GET | /api/stats/by-level | 按级别统计日志数量 |
| GET | /api/stats/top-streams?n=5 | TOP N 高频日志流 |
| GET | /api/stats/overview | 综合概览（各实体总数与活跃数） |
| GET | /api/stats/search?q=xxx&start=RFC3339&end=RFC3339 | 全文检索（消息+标签） |

## 统一响应格式

```json
{"code":0,"message":"ok","data":...}
```

错误映射：
- 400：参数校验失败
- 404：记录不存在
- 409：记录已存在或状态冲突
- 500：服务器内部错误
