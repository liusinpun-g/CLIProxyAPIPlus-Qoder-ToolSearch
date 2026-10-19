# CLIProxyAPI Plus — Qoder Tool Search

[English](README.md) | 中文

这是基于 [kaitranntt/CLIProxyAPIPlus](https://github.com/kaitranntt/CLIProxyAPIPlus) `v7.3.12-1` 的社区 Fork，为 Qoder 模型补充 Codex / OpenAI Responses 工具搜索适配。当前功能属于实验性支持。

本项目独立维护，并非 CLIProxyAPI、CLIProxyAPI Plus、Qoder 或 OpenAI 的官方版本。

## 解决什么问题

工具很多时，Codex 可以先保留精简的工具入口，需要时再搜索并加载工具定义，减少开场工具声明占用。适配逻辑合在 CPA Plus 程序内部，继续使用原来的端口和 Qoder 登录，不需要额外启动一层中转。

这项改动不会扩大模型的上下文窗口，也不会消除系统提示词等基础开销；实际节省量取决于客户端及工具目录。当前改动范围仅限 Qoder。

## 如何开启

使用 Go 1.26 或更新版本编译：

```sh
go build -o cli-proxy-api-plus ./cmd/server
```

从 `config.example.yaml` 创建配置，沿用 CPA Plus 的 Qoder 登录与配置流程。在已有的 `qoder` 节点下加入开关，保留原来的端点等设置：

```yaml
qoder:
  tool-search: true
```

启动服务：

```sh
./cli-proxy-api-plus -config config.yaml
```

默认关闭。Codex 客户端也需要支持工具搜索；使用静态模型目录时，仅对接入此版本、且已开启服务端开关的 Qoder 模型设置 `supports_search_tool: true`。重启客户端后新建会话测试。CC Switch 等目录管理工具重新生成文件时，可能覆盖手动设置，需要复核。

技术细节见 [集成说明](docs/qoder-tool-search.md)。其他提供商和原有配置功能见 [上游中文 README 快照](README_UPSTREAM_CN.md)。

## 验证范围

已使用 Qoder 的 Qwen Flash 验证普通与流式工具搜索、发现、调用及结果回传，以及真实 Codex CLI 对临时只读 MCP 工具的调用。相关包测试和专项竞态检查通过。

全量测试存在一个已在未修改上游版本复现的 WebRTC 失败：`TestPionMediaRelayBridgesAudioAndDataChannel`，错误为 `upstream DataChannel was not created`。因此不宣称全量测试通过。其他 Qoder 模型、大型真实插件目录、桌面端长会话及其他系统的二进制尚未完成此功能的实测。

## 来源与致谢

- 原始项目：[router-for-me/CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI)，提供代理架构与基础能力。
- 直接 Fork 来源：[kaitranntt/CLIProxyAPIPlus](https://github.com/kaitranntt/CLIProxyAPIPlus)，基于 `v7.3.12-1` / `04048d868e96a0450abae26728fef39c734c30d0`。
- 工具搜索协议核心：[hrygo/CLIProxyAPI](https://github.com/hrygo/CLIProxyAPI)，取自 `v1.3.2-upstream8.0.23` / `4edf1ea22abc7ef8e55f161319d3b208db66b034`。`internal/responsestools/` 中 52 个 Go 源码及测试文件保持原样，保留其许可证与来源记录。
- 本 Fork 的新增工作：Sinpun 与贡献者完成 Qoder 的配置开关、执行器内适配、能力声明、集成测试和发布文档；协议核心的作者归属仍属于原作者。

整体沿用宽松的 [MIT 许可证](LICENSE)，允许商用、修改和再分发，同时保留原作者版权及许可证声明。详见 [第三方来源说明](THIRD_PARTY_NOTICES.md)。Go 依赖继续遵守各自许可证。

## 发布

首版拟使用 `v7.3.12-qoder-tools.1`，标记为实验性预发布版。见 [发布说明](RELEASE_NOTES.md) 和 [GitHub 发布指南](PUBLISHING.md)。上游自动化保存在 `.github/upstream-workflows/`，不在本 Fork 中自动执行；本 Fork 不自动合并上游或发布 Docker 镜像。
