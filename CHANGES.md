# Changelog

## v1.1.0 · 集成版本（历史优点整合）

**发布日期：2026-10-18**

本次版本将 TSG 网关项目历史 v3.x–v4.x 的成熟优点整合到 TarsPortkey MCP 桥接器中：

- **配置文件支持**（源自 TSG 治理层配置管理）
  - `config.Load` 现支持 JSON 配置文件路径，不再仅限于环境变量
  - 命令行新增 `-config` 参数指定配置文件
  - 配置结构：`{"endpoint": "...", "api_key": "...", "timeout": 30}`

- **指数退避重试机制**（源自 TSG v3.2.2 连接池/断路器可靠性设计）
  - `InvokeRaw` 新增最多 3 次自动重试
  - 退避间隔：100 ms → 200 ms → 400 ms
  - 仅对 5xx 服务端错误和网络超时/连接失败触发重试，4xx 客户端错误不重试

- **安全响应头加固**（源自 TSG WAF 安全层）
  - 保留并强化 `X-Content-Type-Options: nosniff`、`X-Frame-Options: DENY`
  - CORS 预检支持完整

- **测试覆盖扩展**
  - 新增配置文件加载 3 个测试用例（正常加载、缺失 endpoint、默认超时）
  - 新增重试机制 3 个测试用例（失败恢复、重试耗尽、4xx 不重试）
  - 全量测试：`go test -race ./...` 通过

## v1.0.0 · 初始版本

- MCP stdio / HTTP 双模式支持
- 7 个 Portkey 工具：list_models、chat、embeddings、feedback、logs、health、status
- 环境变量配置：TARS_MCP_ENDPOINT、TARS_MCP_API_KEY、TARS_MCP_TIMEOUT
- 基础健康检查与模型列表接口
