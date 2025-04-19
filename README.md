# bangumi-cli

一套管理季度动画收藏和媒体流相关任务的 CLI 工具。

## 本地运行 & 构建

请先在项目根目录下创建 `.env` 文件：

```sh
cd bangumi-cli
touch .env
```

```.env
LOCAL_SERVER_PORT=8765
BANGUMI_CLIENT_ID=<bangumi APP ID, 可在开发者平台获取>
BANGUMI_CLIENT_SECRET=<bangumi APP Secret, 可在开发者平台获取>
QBITTORRENT_SERVER=http://localhost:8080
QBITTORRENT_USERNAME=admin
QBITTORRENT_PASSWORD=<可在 qbittorrent web ui 中查看/更改>
MIKAN_IDENTITY_COOKIE=<蜜柑计划 identity cookie>
```

根据你的系统使用对应的命令来编译：

```sh
# windows
go build -o bangumi.exe

# linux, macos
go build -o bangumi
```

## 安装

```sh
# windows
mv .\bangumi.exe C:\Users\<usrname>\go\bin -Force

# macos
mv bangumi /usr/local/bin
```

## 使用

```sh
# 登录 bangumi.tv 并获取 API 访问令牌（access_token）
bangumi login

# 订阅 Mikan 最新季度动画（解析、生成元数据、预下载准备）
bangumi subscribe

# 查询 Mikan 订阅的 RSS 是否有新的种子
bangumi update

# 批量处理 RSS，并将下载任务发送到 qBitTorrent 队列
bangumi download

# 递归格式化文件，使其符合 Plex / Jellyfin / Emby 媒体流格式
bangumi format

# 收集已整理的动画数据 (批量收藏bangumi.tv条目、 同步mikan订阅列表)
bangumi collect

# 取消订阅动画
bangumi unsubscribe
```
