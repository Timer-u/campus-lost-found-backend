# 校园失物招领平台后端

校园失物招领系统后端服务：学生可发布失物/招领信息并提交认领申请，失物招领管理员负责审核与物品状态管理，系统管理员负责账号、公告与数据统计。

- 接口契约：[docs/openapi.yaml](docs/openapi.yaml)（可导入 Apifox/Postman 调试）

## 在线环境

| 项 | 地址 |
|---|---|
| 接口基地址 | `http://your-server-ip/api/v1` |
| 健康检查 | `http://your-server-ip/health` |
| 图片访问 | `http://your-server-ip/uploads/<文件名>` |

## 技术栈

Go 1.27 · Gin · GORM · MySQL 8 · JWT（golang-jwt/v5）· bcrypt · GitHub Actions（CI/CD，阿里云 ECS 部署）

## 目录结构

```text
campus-lost-found-backend/
├── main.go                 # 入口：配置 → 数据库 → 路由
├── config/                 # Viper 配置加载（config.yaml 不入库，见 config.example.yaml）
├── model/                  # GORM 模型与建表迁移（User/Item/Claim/Announcement）
├── controller/             # HTTP 层：参数绑定与响应
├── service/                # 业务层：用户/物品/认领/公告/统计/上传
├── middleware/             # panic 恢复、JWT 鉴权、角色校验
├── pkg/
│   ├── auth/               # JWT 签发与解析
│   ├── response/           # 统一响应 {code,msg,data} 与业务错误码
│   └── util/               # 分页等工具
├── router/                 # 路由注册
├── docs/                   # OpenAPI 接口契约
├── scripts/                # seed.sql 演示数据
└── .github/workflows/      # CI/CD
```

## 快速开始

1. 环境要求：Go 1.27+、MySQL 8.0
2. 准备配置：`cp config/config.example.yaml config/config.yaml`，填写数据库账号与 JWT secret
3. 启动：`go run .`（首次启动 AutoMigrate 自动建表）
4. （可选）导入演示数据：`mysql -u<用户> -p <库名> < scripts/seed.sql`

演示账号（执行 seed.sql 后可用，密码均为 `password123`）：

| 学号 | 姓名 | 角色 |
|---|---|---|
| 302026000001 | 李四 | 失物招领管理员（lost_admin） |
| 302026000002 | 王五 | 系统管理员（system_admin） |
| 302026000003 | 赵六 | 学生（student） |

## 接口概览（29 个，详见 docs/openapi.yaml）

| 模块 | 接口 |
|---|---|
| 认证 | 注册（201）、登录（签发 JWT）、获取当前用户、退出登录 |
| 物品（公开） | 列表（keyword/类型/分类/状态/地点筛选 + 排序 + 分页）、详情 |
| 物品（登录） | 发布（默认待审核）、编辑（修改后回到待审核）、删除、我的发布 |
| 图片 | 上传（单张 ≤5MB，每个物品 ≤6 张），`/uploads/` 静态访问 |
| 认领 | 提交申请（仅开放中的招领）、物品的申请列表（发布者/管理员）、我的申请、取消 |
| 管理员审核 | 物品列表（待审核优先）、物品审核、物品状态管理、认领审核（事务防一物多通过） |
| 用户管理 | 用户列表/搜索、角色变更、启用禁用（system_admin） |
| 公告 | 公开列表、管理端增删改查（system_admin） |
| 统计 | 用户/物品/待审核/认领/已解决数量与类型分布（两级管理员） |

响应统一为 `{code, msg, data}`：`code=0` 成功；业务错误码一码一 msg（10001 参数错误、10007 用户名已存在、16001 图片超限等，完整定义见 docs/openapi.yaml）。

## 部署

服务器采用直连部署：二进制 + systemd（`campus-backend.service`）+ Nginx 反代（`/api/`、`/uploads/`、`/health`），部署目录 `/opt/campus/backend`。

CI/CD：push 到 main 自动执行 `go build` + `go test`，通过后交叉编译并 SSH 部署重启（`.github/workflows/deploy.yml`，需要在仓库 Secrets 中配置 `ECS_HOST`、`ECS_SSH_KEY`）。PR 只跑构建与测试，作为合并门禁。

## 联调约定

- 前端请求前缀：`/api/v1`；登录响应中的 `accessToken` 放入 `Authorization: Bearer <token>`，退出登录为无状态设计，前端删除本地 Token 即可
- 公开列表/详情仅展示审核通过（approved）的物品；新发布物品默认待审核，不会立即公开
- 认领申请仅针对审核通过且处于开放状态的招领信息；发布者本人与重复申请会被拒绝
- 中文查询参数（keyword 等）需 URL 编码（Apifox/axios 会自动处理）
