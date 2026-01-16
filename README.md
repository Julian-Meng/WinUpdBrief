# WinUpdBrief

## 简介 / Introduction

WinUpdBrief 是一个轻量级的 Windows 系统工具，用于快速获取和展示当前系统的版本信息及最新的Windows更新情况。

WinUpdBrief is a lightweight Windows utility tool designed to quickly fetch and display the current system version information and the latest Windows update status.

## 功能特性 / Features

- **系统版本查询** - 从注册表读取Windows系统版本、Build号等详细信息
- **更新检索** - 从Microsoft官方渠道获取针对当前Build的最新更新记录
- **KB文章获取** - 自动获取更新对应的KB文章内容和摘要
- **简洁展示** - 以易读的文本格式呈现所有信息

## 文件结构 / Project Structure

```
WinUpdBrief/
├── main.go              # 主程序入口
├── go.mod               # Go模块配置
├── winver/
│   └── winver.go        # 系统版本信息读取（从注册表）
├── updates/
│   ├── fetch.go         # 更新信息抓取
│   ├── history.go       # 更新历史记录处理
│   └── kb.go            # KB文章内容处理
├── render/
│   └── text.go          # 信息展示渲染
└── README.md            # 项目说明
```

## 工作原理 / How it Works

1. **读取本地系统信息** - 通过Windows注册表获取当前系统版本和Build号
2. **抓取更新信息** - 查询Microsoft官方更新页面，获取对应Build的最新补丁
3. **获取KB详情** - 提取更新对应的Knowledge Base文章内容
4. **格式化输出** - 整合所有信息并以友好的文本格式呈现

## 技术栈 / Tech Stack

- **Go 1.25.5** - 主要开发语言
- **goquery** - HTML解析库
- **golang.org/x/sys** - Windows系统调用接口
