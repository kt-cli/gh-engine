# 上游同步策略

上游 `https://github.com/cli/cli`，MIT；origin 为 `kt-cli/gh-engine`，upstream 只读。
KT 内核基线固定官方 `v2.102.0`，保留全部上游历史，消费方以不可变提交锁定。
`main` 为稳定内核分支；`trunk` 保留 Fork 创建时的上游镜像，不供生产消费。
每个正式 release 评估同步，安全修复即时评估，季度至少一次。在 `upstream-sync`
合并指定 release tag，PR 到 main。安全修复以上游为准，上游已有接缝后移除 KT 定制。

定制白名单只有 `internal/config/transport.go`、`pkg/ktengine/`、本策略、KT 专用 CI。
另有 API 注入所需的三个加法接缝：`api/http_client.go` HTTPClientOptions.Transport、
`pkg/cmdutil/factory.go` Factory.APITransport、`pkg/cmd/api/api.go` 传递该选项。
原值 nil 时保持上游行为，保留 api 的逐命令日志、缓存和请求头。
不修改命令业务、REST/GraphQL 实现、文件语义或签名验证。此 Fork 无凭据、无授权规则、无生产部署；由 kt-gh 链接，发布走中央组件工作流。

验证：配置不读环境/文件/keyring、不返回 token；所有三类 HTTP factory 注入；上游命令树生成能力目录；官方真实命令的文件/输出行为及 Broker 边界；`go fix`、全包测试、lint、kt-gh 接缝测试。升级逐项复查新命令与网络客户端，能力增加不自动授予机器权限。
