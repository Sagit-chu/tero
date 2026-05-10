# Flvx Monitor

Flvx Monitor 是一个自动化节点监控与替换系统。该系统用于实时监控 Linux 服务器节点的连通性（检测是否宕机或被 GFW 封锁），并在检测到节点失效时，自动通过 [flvx 面板 API](https://github.com/Sagit-chu/flvx) 进行节点替换，随后自动更新 Cloudflare 的 DNS 解析。

## 🌟 核心特性

- **两阶段连通性验证 (Two-Stage Verification)**:
  - **阶段 1 (全局检查)**: 使用本地 Ping 测试节点是否存活（检测宕机）。
  - **阶段 2 (被墙检测)**: 当阶段 1 成功时，调用配置好的第三方 API（如 ITDog/Ping.pe）从中国大陆探测节点连通性。
- **自动化无缝替换**:
  - 当确认节点被墙或宕机时，系统会自动从预设的“备用节点池”中提取一个可用节点（包含 IP、SSH 端口、密码）。
  - 自动调用 `flvx` 接口完成节点替换。
  - 替换完成后，自动调用 Cloudflare API 刷新域名的 A 记录。
- **现代化 Web 面板**:
  - 基于 React 18 + Vite + Tailwind CSS 构建的现代化单页应用 (SPA)。
  - 实时查看当前节点状态（存活、宕机、被墙、替换中）。
  - 可视化管理备用节点池（添加、编辑、删除）。
- **极简部署**:
  - 后端采用 Go 语言编写，使用 SQLite 嵌入式数据库。
  - 最终可编译为单一二进制文件运行，无需复杂的环境依赖。

## 🏗️ 系统架构

- **后端 (Backend)**: Go 1.22+ 
  - 负责执行定时监控任务。
  - 提供 RESTful API 供前端面板调用。
  - 使用 `mattn/go-sqlite3` 管理数据库（存储备用节点和系统配置）。
- **前端 (Frontend)**: React + Vite + Tailwind CSS
  - 提供响应式的监控与管理仪表盘。
- **测试 (Testing)**:
  - **单元测试**: 后端 Go `testing`，前端 `Vitest` + `React Testing Library`。
  - **端到端测试 (E2E)**: `Playwright` 自动化全链路测试。

## 🚀 快速开始

### 1. 运行后端服务

确保已安装 Go 1.22 或更高版本。

```bash
# 获取依赖并运行 API 服务器 (默认运行在 8080 端口)
go mod tidy
go run backend/api/cmd/main.go
```

### 2. 运行前端面板

确保已安装 Node.js 和 npm。

```bash
cd frontend
npm install

# 启动开发服务器 (默认运行在 5173 端口，会自动代理 /api 请求到后端)
npm run dev
```

打开浏览器访问 `http://localhost:5173` 即可查看 Flvx Monitor 仪表盘。

## 🧪 运行测试

本项目包含完整的自动化测试套件。

**后端单元测试 (Go)**
```bash
go test ./backend/... -v
```

**前端组件测试 (Vitest)**
```bash
cd frontend
npm test
```

**端到端自动化测试 (Playwright E2E)**
```bash
cd e2e
npm install
npx playwright install --with-deps chromium

# 运行 E2E 测试 (会自动启动后端和前端服务)
npx playwright test
```

## 📂 目录结构

```text
.
├── backend/            # Go 后端代码 (API、监控逻辑、数据库、存储库)
│   ├── api/            # RESTful API 控制器及入口
│   ├── db/             # SQLite 数据库初始化
│   ├── monitor/        # 监控循环与两阶段验证逻辑
│   └── repository/     # 数据库 CRUD 操作
├── frontend/           # React 前端代码
│   ├── src/            # 面板组件及页面
│   └── vite.config.ts  # Vite 配置及后端代理
├── e2e/                # Playwright 端到端测试配置及用例
├── docs/               # 设计规范与实施计划文档
└── README.md           # 项目说明文档
```

## 📝 待办事项 (TODO)

- [x] 实现完整的监控循环 (Monitor Loop) 及两阶段验证逻辑。
- [x] 编写全量测试以及端到端测试 (Playwright)。
- [ ] 接入 `flvx` API 实现备用节点的自动化替换。
- [ ] 接入 Cloudflare API 实现 DNS 记录更新。
- [ ] 完善前端面板的备用节点池管理功能。
