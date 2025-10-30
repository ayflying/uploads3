# uploads3

一个用于批量上传本地目录文件到 S3 兼容对象存储的命令行工具，基于 GoFrame 与 minio-go SDK。

## 功能特性
- 批量扫描本地目录并保持原有层级结构上传至 S3
- 并发上传，默认 `10`，可配置，最大 `50`
- 通过 `config/local.yaml` 管理多套 S3 配置（如 MinIO、七牛 S3）
- 统一路径映射：自动将 Windows `\` 转换为对象存储使用的 `/`

## 环境要求
- Go `1.23+`（仓库声明了 `toolchain go1.24.8`，本地使用 Go 1.23/1.24 均可）
- 可访问你的 S3 兼容对象存储（如 MinIO、七牛云 S3）

## 安装与构建
- 直接构建（Windows/Linux/macOS 均可）
  - Windows: `go build -o uploads3.exe .`
  - Linux/macOS: `go build -o uploads3 .`
- 运行
  - Windows: `./uploads3.exe -p <本地目录> -u <S3根路径> -w <并发数>`
  - Linux/macOS: `./uploads3 -p <本地目录> -u <S3根路径> -w <并发数>`

可选（需安装 GoFrame CLI 与 Make）：`make build` 会在 `bin/` 下生成各平台二进制。

## 配置说明
应用读取 `config/local.yaml` 的 `s3` 配置，示例：

```yaml
s3:
  type: "default"            # 当前使用的配置名称（default/qiniu/...）
  default:
    provider: "minio"
    accessKey: "<你的AccessKey>"
    secretKey: "<你的SecretKey>"
    address: "<你的S3地址>:<端口>"    # 如 minio: ay.cname.com:9000
    ssl: false                   # https 用 true
    url: "http://host/bucket/"   # 直链或 CDN 前缀，可选，用于拼接访问地址
    bucketName: "<桶名>"
  qiniu:
    provider: "qiniu"
    accessKey: "<你的AccessKey>"
    secretKey: "<你的SecretKey>"
    address: "s3.cn-south-1.qiniucs.com"
    ssl: true
    url: "https://attachment.example.com/"
    bucketName: "<桶名>"
```

- `s3.type` 决定使用哪套配置（上例为 `default`）。
- `address` 是 S3/MinIO 的服务地址；`ssl` 选择是否 https。
- `bucketName` 是目标桶名；上传时的对象 Key 将由本地路径与 `upload_path` 映射生成。

## 使用指南
命令行参数：
- `-p, --path` 本地文件夹路径（必须）
- `-u, --upload_path` S3 上传根路径（必须）
- `-w, --worker` 并发数，默认 `10`，最大 `50`

运行示例（Windows）：
```
./uploads3.exe -p D:\data\images -u cdn.yoyaworld.com/upload/images -w 20
```
运行示例（Linux/macOS）：
```
./uploads3 -p /data/images -u cdn.yoyaworld.com/upload/images -w 20
```

### 路径映射规则
- 工具会扫描 `path` 下所有文件，并以相对层级映射到 `upload_path`。
- 对象 Key 生成规则：用 `upload_path` 替换本地绝对路径中的 `path` 前缀，并统一改为 `/`。
  - 例如：
    - 本地文件：`D:\data\images\a\b.jpg`
    - 参数：`-p D:\data\images`，`-u cdn.yoyaworld.com/upload/images`
    - 对象 Key：`cdn.yoyaworld.com/upload/images/a/b.jpg`

### 并发与性能
- `-w` 控制同时上传的并行 worker 数，范围 `1~50`。
- 大量小文件场景建议适度提高并发；网络与目标 S3 的限流会影响实际吞吐。

## 日志输出
- 使用 GoFrame 日志，默认打印到控制台，包含当前上传进度 `(已传/总数)` 与对象 Key。
- 失败时会记录错误日志；可根据输出排查凭据与网络问题。

## 常见问题
- 连接失败：检查 `address` 与 `ssl` 是否匹配、端口是否可达。
- 认证失败：核对 `accessKey/secretKey`、桶权限与跨区 endpoint 设置。
- Key 不符合预期：确认命令行的 `path` 与 `upload_path` 是否书写正确，尤其是 Windows 路径需要使用转义或用双引号包裹。

## 开发者提示
- 入口：`main.go` 调用 `internal/cmd/cmd.go` 的命令。
- S3 封装：`internal/logic/s3/s3.go` 基于 `minio-go`，从 `config/local.yaml` 读取配置。
- 若需要扩展不同提供商，可在 `config/local.yaml` 增加新节，并通过 `s3.type` 切换。

## 许可证
- 使用 `MIT` 许可，详见 `LICENSE` 文件。