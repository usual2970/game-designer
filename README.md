# Game Designer Plugin

Game Designer 是一个面向代码 Agent 的 H5 游戏开发插件。它把 Phaser 前端模板、Go 游戏后端模板、TypeScript SDK、OpenAPI 合约、部署 CLI 和 Agent skills 打包在一起，让 Agent 可以从零创建、联调、测试并部署一个服务端判定结果的老虎机 H5 游戏。

插件适用于 Claude Code 和 Codex。仓库根目录就是插件根目录，安装时必须从这里导入，确保 `skills/`、`server-template/`、`frontend-template-phaser/`、`sdk-js/`、`contracts/`、`cli/` 等资源能被 Agent 发现。

## 你可以用它做什么

| 能力 | 说明 |
|------|------|
| 创建 H5 前端 | 从 Phaser + TypeScript + Vite 模板生成浏览器可玩的游戏前端 |
| 创建 Go 后端 | 从模板生成带 session、profile、balance、spin、history、leaderboard 的游戏服务 |
| 接入 SDK | 在 H5 游戏中接入 `@game-designer/sdk`，完成登录、余额、配置、spin 等调用 |
| 生成老虎机玩法 | 使用服务端权威 spin 结果，不在客户端计算赔付 |
| 定制主题 | 修改标题、颜色、符号素材、音效和移动端布局 |
| 本地验证 | 校验插件包、合约、后端、SDK、CLI、前端构建和浏览器表现 |
| 部署游戏 | 通过 Go CLI 发布 backend、socket、frontend 三个 surface |
| 定位问题 | 拆分前端、SDK、服务端、合约和部署问题，给出诊断路径 |

## 使用流程总览

推荐按下面顺序使用插件：

```text
安装插件
  -> 首次构建部署 CLI
  -> 创建后端
  -> 创建或接入 H5 前端
  -> 添加老虎机玩法与主题
  -> 本地测试和预部署检查
  -> 打包前端
  -> 部署到 PaaS
  -> 验证线上地址
```

对应的 Agent skills：

| 阶段 | Skill | 何时使用 |
|------|-------|----------|
| 首次准备 | `gd-setup-cli` | 安装插件后第一次使用，或 CLI 源码变化后 |
| 后端 | `gd-create-server` | 需要在目标项目中创建 Go 游戏后端 |
| SDK | `gd-connect-sdk` | 已有 H5 项目，需要接入 Game Designer SDK |
| 前端 | `gd-create-h5-game` | 需要创建 Phaser H5 前端 |
| 玩法 | `gd-create-slot-game` | 需要生成或补全老虎机玩法 |
| 主题 | `gd-theme-h5-game` | 需要定制视觉、素材、音效、文案 |
| 前端验证 | `gd-test-h5-game` | 前端改动后，或打包前 |
| 前端排错 | `gd-debug-h5-game` | 白屏、素材 404、canvas、音频、构建等问题 |
| 前端打包 | `gd-package-frontend` | 部署前确认 `frontend/dist/` 可发布 |
| 预部署 | `gd-prepare-deploy` | 部署前跑 CLI preflight、本地构建和测试 |
| 部署 | `gd-deploy-game` | 发布到 fake provider 或 3os PaaS |
| 集成排错 | `gd-debug-integration` | SDK、后端、合约、部署、线上验证失败 |

## 1. 安装插件

### Claude Code 本地调试

```bash
git clone <repo-url> game-designer-backend
cd game-designer-backend
claude --plugin-dir .
```

插件更新后可在 Claude Code 会话中执行 `/reload-plugins`。

### Codex 本地导入

1. 打开 Codex Plugins UI
2. 添加本地插件源，路径选择本仓库根目录
3. 启用 `game-designer`
4. 在会话中用 `$gd-setup-cli`、`$gd-create-server` 等方式触发技能

安装后验证插件包结构：

```bash
./scripts/verify-plugin-package.sh
```

更完整的安装说明见 [Plugin installation](docs/integration/plugin-installation.md)。

## 2. 首次准备：构建部署 CLI

插件安装不会自动编译 CLI。首次部署相关操作前，先让 Agent 执行 `gd-setup-cli`，或手动构建：

```bash
cd cli
GOWORK=off go build -o game-designer ./cmd/game-designer
./game-designer version
```

后续每次更新 `cli/` 源码后，都应重新执行这一步。

## 3. 选择使用路径

### 路径 A：从零创建完整游戏

适合没有现成项目，需要 Agent 创建完整的后端、前端、玩法和部署流程。

```text
gd-setup-cli
  -> gd-create-server
  -> gd-create-h5-game
  -> gd-create-slot-game
  -> gd-theme-h5-game (可选)
  -> gd-test-h5-game
  -> gd-package-frontend
  -> gd-prepare-deploy
  -> gd-deploy-game
```

本地运行时，后端默认监听 `http://localhost:8080`，前端默认使用 Vite dev server：

```bash
cd server-template
GOWORK=off go run ./cmd/server

cd ../frontend-template-phaser
npm install
npm run dev
```

打开 `http://localhost:3000`，前端会通过 SDK 访问 `http://localhost:8080`。

### 路径 B：已有 H5 前端，只接入后端和 SDK

适合已有 Phaser、Canvas、React、原生 H5 或其他前端项目，只需要接入 Game Designer 的服务端能力。

```text
gd-setup-cli
  -> gd-create-server
  -> gd-connect-sdk
  -> gd-prepare-deploy
  -> gd-deploy-game
```

SDK 基础用法：

```typescript
import { GameDesignerClient } from "@game-designer/sdk";

const client = new GameDesignerClient({ baseUrl: "http://localhost:8080" });

await client.createOrResumeSession({
  playerId: "player-123",
  nickname: "Alice",
});

const config = await client.getSlotConfig();
const balance = await client.getBalance();
const result = await client.spin({ wager: 10 });
const history = await client.getSpinHistory({ limit: 20 });
const leaderboard = await client.getSlotLeaderboard({ limit: 10 });
```

关键约束：客户端只负责表现和交互，spin 结果、赔付、余额变化都以服务端返回为准。

### 路径 C：只验证或维护插件本身

适合开发本插件、修改模板、SDK、CLI 或 skill 文档。

```bash
# 插件包结构
./scripts/verify-plugin-package.sh

# 合约、服务端、SDK、CLI、本地接口循环
./scripts/verify-local.sh

# 单独测试
cd server-template && GOWORK=off go test ./... -v
cd sdk-js && npm test
cd cli && GOWORK=off go test ./... -v
cd frontend-template-phaser && npm test
```

`./scripts/verify-local.sh` 的 live slot loop 需要后端已运行；如果没有运行，会跳过 live 接口验证并提示启动命令。

## 4. 部署流程

部署前推荐执行：

```text
gd-test-h5-game
  -> gd-package-frontend
  -> gd-prepare-deploy
```

本地 dry run 可用 fake provider：

```bash
cd cli
GOWORK=off go run ./cmd/game-designer deploy \
  --server-path ../server-template \
  --app-name my-game \
  --env production \
  --provider fake
```

生产部署使用 3os provider，并通过环境变量提供凭据：

```bash
export GD_IDENTIFIER="<identifier>"
export GD_PASSWORD="<password>"

cd cli
./game-designer deploy \
  --provider 3os \
  --mode create \
  --identifier "$GD_IDENTIFIER" \
  --password "$GD_PASSWORD" \
  --game-name "My Game" \
  --package-path ./game.zip \
  --version 1.0.0 \
  --change-log "Initial release" \
  --backend-dir "admin" \
  --backend-cmd "./server -type admin" \
  --frontend-dir "h5" \
  --socket-dir "logic" \
  --socket-cmd "./server -type logic"
```

线上验证：

```bash
./scripts/verify-deployed.sh https://<deployed-url>
```

部署细节见 [PaaS provider](docs/deployment/paas-provider.md)、[Deployed verification](docs/deployment/deployed-verification.md) 和 [CLI README](cli/README.md)。

## 项目结构

```text
contracts/                OpenAPI contract，后端和 SDK 的单一接口事实来源
server-template/          Go slot machine backend template
frontend-template-phaser/ Phaser + TypeScript + Vite H5 frontend template
sdk-js/                   TypeScript H5 SDK
cli/                      Go deploy CLI
skills/                   Agent-facing plugin skills
examples/                 Example H5 slot machine game
scripts/                  Verification scripts
docs/                     Integration, deployment, planning docs
```

## 排错入口

| 问题 | 优先入口 |
|------|----------|
| 插件安装后看不到 skills | [Plugin installation](docs/integration/plugin-installation.md) |
| CLI 不存在或版本不对 | `gd-setup-cli` |
| 后端构建、接口、session、balance、spin 异常 | `gd-debug-integration` |
| SDK 类型、请求、错误码异常 | `gd-debug-integration` |
| Phaser 白屏、canvas 尺寸、素材 404、移动端显示异常 | `gd-debug-h5-game` |
| 前端构建产物不可部署 | `gd-package-frontend` |
| 部署命令失败或线上验证失败 | `gd-debug-integration` |

更多说明：

- [Agent golden path](docs/integration/agent-golden-path.md)
- [Contract-first workflow](docs/integration/contract-first-workflow.md)
- [Local verification](docs/integration/local-verification.md)
- [SDK usage](docs/integration/sdk-usage.md)
- [Troubleshooting](docs/deployment/troubleshooting.md)

## License

MIT
