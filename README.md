# SwitchHosts-GO

```
    _____         _ _       _    _    _           _         _____  ____  
   / ____|       (_) |     | |  | |  | |         | |       / ____|/ __ \ 
  | (_____      ___| |_ ___| |__| |__| | ___  ___| |_ ___ | |  __| |  | |
   \___ \ \ /\ / / | __/ __|  __  |  __  |/ _ \/ __| __/ _ | | |_ | |  | |
   ____) \ V  V /| | | | (__| |  | |  | | (_) \__ \ |_\__ \| |__| | |__| |
  |_____/ \_/\_/ |_|\__\___|_|  |_|  |_|\___/|___/\__|___/ \_____|_____/ 
```

SwitchHosts-GO是一个用Go语言和Wails框架（Go+Web前端）开发的hosts文件管理工具，帮助您快速切换不同的hosts配置。无论是在开发、测试环境切换，还是在不同网络环境之间切换hosts配置，SwitchHosts-GO都能提供流畅的使用体验。
该项目是参照SwitchHost,使用cursor完成的 95%的代码来自于cursor(包括readme.md文件) 主要是尝鲜,第一次尝试使用cursor做一个App出来

## 环境要求

- **操作系统**：Windows 10/11, macOS 10.13+, 主流Linux发行版
- **权限要求**：需要管理员权限才能修改系统hosts文件
- **硬盘空间**：< 100MB
- **内存**：最低256MB

## 功能特点

- 简洁易用的用户界面
- 快速切换不同的hosts配置
- 支持追加模式和独立模式
- 实时生效，无需重启应用
- 数据存储位置可自定义，支持热加载
- 自动备份系统hosts文件
- 支持Windows、macOS和Linux系统

## 系统截图

### 主界面
主界面分为左右两栏，左侧是hosts配置列表，右侧是编辑区域。通过简洁的界面，您可以方便地管理多个hosts配置。

### 追加/独立模式切换
顶部工具栏提供模式切换功能，可以在多配置追加模式和单一配置独立模式之间轻松切换。

### 数据存储位置设置
通过设置按钮，您可以查看和更改数据存储位置。更改后数据会自动迁移并即时生效，无需重启应用。

## 安装方法

### Windows系统

1. 下载最新的`.exe`
2. 双击运行程序,即可开始使用  
> 在修改hosts文件时,可能会遇到权限问题,需要先设置系统hosts文件的读写权限

### macOS系统

暂未编译相应二进制文件

### Linux系统

暂未编译相应二进制文件

#### Debian/Ubuntu

暂未编译相应二进制文件



#### 通用方法


### 从源码构建

1. 确保您已安装Go,wails和Node.js环境
2. 安装Wails CLI：`go install github.com/wailsapp/wails/v2/cmd/wails@latest`
3. 克隆代码仓库：`git clone git@github.com:ldjx7/SwitchHosts-GO.git`
4. 进入项目目录：`cd SwitchHosts-GO`
5. 构建应用：`wails build`

编译后的应用将位于`build/bin`目录中。

## 使用说明

### 基本操作

1. 启动应用后，左侧面板显示hosts配置列表
2. 点击"新建"按钮创建新的hosts配置
3. 在右侧编辑区域修改hosts内容
4. 通过开关按钮启用/禁用hosts配置

### 模式切换

- **追加模式**：多个hosts文件可以同时启用，内容会追加到系统hosts文件中
- **独立模式**：每次只能启用一个hosts文件

### 数据存储位置

应用数据默认存储在以下位置：

- Windows: `%APPDATA%\SwitchHosts-GO`
- macOS: `~/Library/Application Support/SwitchHosts-GO`
- Linux: `~/.switchhosts-go`

您可以通过设置菜单更改数据存储位置，所有数据将会复制到新位置并即时生效。

## 开发指南

### 项目结构

- `app.go`：主应用程序逻辑
- `frontend/`：Vue 3 + Ant Design Vue前端代码
- `services/`：后端服务实现
- `models/`：数据模型定义

### 开发模式

运行`wails dev`启动开发模式，这将启动一个Vite开发服务器，提供前端代码的热重载功能。

访问http://localhost:34115可以在浏览器中开发并调用Go方法。

### 构建应用

运行`wails build`构建可分发的生产模式应用包。

## 许可证

[MIT License](LICENSE)

## 贡献指南

欢迎为SwitchHosts-GO贡献代码、报告问题或提出改进建议！

1. Fork本仓库
2. 创建您的特性分支：`git checkout -b feature/amazing-feature`
3. 提交您的更改：`git commit -m 'Add some amazing feature'`
4. 推送到分支：`git push origin feature/amazing-feature`
5. 开启一个Pull Request

如果您发现了bug或有功能建议，也可以直接在[Issues页面](https://github.com/your-username/SwitchHosts-GO/issues)提出。

## 致谢

- [Wails](https://wails.io)：用于构建跨平台桌面应用的Go框架
- [Vue 3](https://vuejs.org)：前端框架
- [Ant Design Vue](https://antdv.com)：UI组件库
- [oldj/SwitchHosts](https://github.com/oldj/SwitchHosts)：提供了灵感
