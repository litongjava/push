---
name: litongjava-push
description: "使用 litongjava/push 工具编写、检查和执行 push.txt，顺序向多个 Git 远程推送或执行逐行命令；排查工作目录、参数拆分、命令失败后继续及退出码误判。适用于用户提到 push 工具、litongjava/push 或 push.txt；不是所有普通 git push 请求。"
whenToUse: "用户要求配置或运行 litongjava/push、编写或修复 push.txt，或排查该工具批量推送问题时。"
---

# litongjava push 工具

## 工具与适用范围

模块为 `github.com/litongjava/push`，可执行文件名为 `push`（Windows 为 `push.exe`）。工具从当前工作目录读取固定名称 `push.txt`，按行启动外部程序。它不是 Shell，也不是 Git 子命令；虽然常用于多个远程推送，也可执行其他外部程序。

本文依据工具仓库的 main.go、go.mod 和 push.txt。当前 go.mod 要求 Go 1.21.1；README 的 Go 1.18 及“有效 shell 命令”描述不够准确，行为以源码为准。后续版本有变化时先核对源码。

## 操作流程

1. 明确用户要编写配置、排错，还是实际执行推送；确认目标项目的绝对路径。工具源码目录不一定是目标项目目录。
2. 用读取工具检查已有 push.txt，再修改；不存在时才创建。配置可启动任意外部程序，运行前逐行审查，不能把不可信配置直接当作安全数据执行。
3. 用 `Get-Command push -All`（Windows）或 `command -v push`（Unix）定位程序，确认没有同名函数或程序冲突，必要时使用可执行文件绝对路径。不要用 `push --help`、`push --version` 或 `push --dry-run` 探测：程序不处理参数，仍会执行当前目录的 push.txt！
4. Git 推送前检查 `git status --short --branch`、`git branch --show-current`、`git remote -v`，必要时检查指定本地分支是否存在。确认每行远程、分支和用户授权范围，不从示例猜测配置。远程 URL 如含凭据，汇报时必须脱敏。
5. 仅在用户要求实际推送或执行时运行。仅编写 skill 或配置不授权执行远程推送。不自动添加提交、创建远程、强制推送、删除远程分支、使用 --mirror 或绕过保护规则。需要这些操作时另行取得明确授权。
6. 在目标项目目录启动工具，保留逐行日志；逐项判断成功或失败，不能只看进程退出码或最后一行。
7. 若部分失败，先分析原因，只重试失败且仍获授权的命令；不要不加判断地重跑整份配置，因为成功的其他命令可能有不可重复的副作用。

## 编译与安装工具

已有可用程序时无需重装。需要编译且用户授权时，在工具源码目录执行：

```powershell
go version
go build -o push.exe .
```

Linux/macOS 构建命令为 `go build -o push .`。需要从模块安装时可用 `go install github.com/litongjava/push@latest`，这会访问网络并安装当前发布版本，不保证与本地源码一致。通过 `go env GOBIN GOPATH` 确认安装位置，不写死用户路径。构建命令是 Go 自带命令，不需要其他 build 工具。

## 配置与运行示例

仅当 origin、gitee 和 main 均经核实是用户目标时，使用：

```text
git push origin main
git push gitee main
```

写为 UTF-8 无 BOM 的 push.txt。每行一个可执行程序及其参数；允许空行，不写注释。Windows 在目标项目目录中运行已核实的程序，例如：

```powershell
# 替换为实际可执行文件路径；工作目录由执行工具的 workdir 指定
& "C:\tools\push.exe"
```

若程序就在当前目标目录，使用 `./push.exe`（Windows）或 `./push`（Linux/macOS）。上面 PowerShell 对可执行文件路径的引号由外层 Shell 解析；不表示 push.txt 内也支持引号。

不要传配置路径参数：`push other.txt` 仍读取当前目录的 push.txt。没有自带预览模式。若用户需要 Git 预检，可以在明确的目标目录单独运行 `git push --dry-run <remote> <branch>`；不要运行 `push --dry-run`。预检不等于最终推送成功。

## 源码确认的限制

| 项目 | 实际行为与应对 |
| --- | --- |
| 配置查找 | 仅相对当前工作目录打开 push.txt，不在可执行文件目录查找，也不向父目录搜索。 |
| 空白与行号 | 每行先 TrimSpace，空行跳过；日志行号仍是文件原始行号。 |
| 参数拆分 | 使用 strings.Fields 按空白拆分，不识别单引号、双引号或转义。带空格的路径、提交消息等不能靠加引号修复。 |
| Shell 语法 | exec.Command 直接启动程序，不解析 &&、管道、重定向、通配符或环境变量表达式。它们会作为普通参数传递。 |
| 注释 | 不支持 #、// 或 REM 注释；非空行会被当作命令执行。 |
| 工作目录 | 子进程继承工具工作目录；单独写 cd 不能改变后续行的目录。Git 可使用不含空格的 `git -C <目录> ...`。 |
| 状态传递 | 每行启动独立进程，不能在前一行设置环境变量供后一行使用；外部文件系统的修改仍会保留。 |
| 交互输入 | 源码只接了 stdout/stderr，没有接 stdin，不依赖交互式输入密码。先单独配置好凭据或 SSH 认证。 |
| 出错处理 | 某行失败仍继续后续行；不能用前一行成功作为后续行的安全前提。 |
| 退出码 | 打开文件失败、扫描失败或子命令失败均仅打印错误；程序可能仍返回 0。 |
| 扫描长度 | 使用默认 bufio.Scanner，单个 token 有约 64 KiB 的上限，不写超长命令。 |
| 编码 | 没有剥离 UTF-8 BOM；BOM 可能污染第一条命令名。使用 UTF-8 无 BOM。 |
| 敏感信息 | 日志打印整行命令，不在配置参数中放令牌、密码或带凭据 URL。 |

需要引号、多个步骤、条件执行或复杂路径时，优先把逻辑放到独立脚本中，再在配置中用明确的解释器调用脚本；配置里的解释器路径、脚本路径及各参数仍不得依赖包含空格的引号解析。例如 `pwsh -NoProfile -File scripts/publish.ps1` 或 `bash scripts/publish.sh`。先核实解释器存在并审查脚本，不能把改写脚本作为扩大执行授权的理由。

## 验证与故障处理

- `Error opening file push.txt`：核对实际 workdir、文件名称和权限，不在其他目录盲目新建替代文件。
- `executable file not found` 或类似错误：检查 PATH、命令拼写、BOM、Shell 内建命令及带空格路径拆分问题。
- `Error executing command on line N`：结合第 N 行及其 stderr 排查；后面的命令可能已成功，不要报告全部失败或全部成功。
- `Error reading file`：检查扫描上限、文件编码和读取错误；剩余行可能没有执行。
- Git non-fast-forward：先检查本地与远端历史，报告差异；不自动 force、reset、rebase 或 pull。认证失败则由用户完成认证，不将凭据写入配置或日志。
- `Push execution finished.` 仅表示扫描流程结束；必须核对每条预期命令是否出现成功日志，以及是否存在任何错误日志。
- 对明确的普通分支推送，可用 `git rev-parse <本地分支>` 与 `git ls-remote <远程> refs/heads/<目标分支>` 对照提交 ID。特殊 refspec 应按实际映射核验；若无法访问远端，明确说明尚未验证。
- 最终汇报配置位置、执行工作目录、每个远程/命令的结果、失败原因和重试建议。仅配置未执行时明确写出“未执行实际推送”。
