<p align="center">
  <img src="resources/icon.png" width="112" height="112" alt="gus logo">
</p>

<h1 align="center">gus</h1>

<p align="center">在 Windows 托盘里，管理每个 Git 项目的提交身份。</p>

<p align="center">
  <a href="https://github.com/Hunlongyu/gus_gui/releases">下载发布版</a> ·
  <a href="https://github.com/Hunlongyu/gus_gui/actions/workflows/build.yml">开发构建</a> ·
  <a href="CHANGELOG.md">更新日志</a> ·
  <a href="LICENSE.txt">MIT License</a>
</p>

个人项目用私人邮箱，工作项目用公司邮箱。gus 将全局默认身份、项目实际使用的姓名和邮箱放在同一个面板中，右键即可切换，提交前就能确认署名。

<p align="center">
  <img src="assets/images/gus-panel.png" width="480" alt="gus 主面板：全局身份、八个示例仓库，以及继承、专用和混合三种身份来源">
</p>

## 功能

- **保存多套身份**：为个人、工作、开源等用途保存姓名与邮箱，查看各身份的使用情况。
- **按项目切换**：设置全局默认身份，或给仓库单独指定身份，也可以恢复继承全局。
- **一眼核对配置**：并排显示仓库、名字、邮箱和配置来源，每次打开面板重新读取 Git 配置。
- **查看推送平台和详情**：项目列表用堆叠图标显示推送平台，用符号标记身份来源，悬浮可展开查看完整信息，支持多个 remote 和多个推送地址。双击项目行可查看已保存账号、当前分支，以及每个远程仓库的拉取和推送 URL；点击 URL 可用默认浏览器打开仓库，常见 SSH 地址自动转换为网页地址。
- **快速找到项目**：支持多选目录、拖入目录、自动定位仓库根目录、去重和按项目名搜索。
- **常驻托盘**：支持面板置顶、调整尺寸、开机自启、简体中文 / English，以及跟随系统、浅色和深色外观。

gus 管理 Git 的 `user.name` 和 `user.email`。GitHub / GitLab 的登录凭据、访问令牌和 SSH 密钥仍由你原有的工具管理。

## 下载与启动

需要 **Windows x64** 和已加入 `PATH` 的 **Git for Windows**。在 PowerShell 运行 `git --version`，确认 Git 可用。

| 获取方式 | 操作 |
| --- | --- |
| 安装包 | 打开 [Releases](https://github.com/Hunlongyu/gus_gui/releases)，下载已发布的 Windows x64 安装程序。尚无发布版时，可使用开发构建。 |
| 开发构建 | 登录 GitHub，打开最近一次成功的 [Windows build](https://github.com/Hunlongyu/gus_gui/actions/workflows/build.yml)，在 **Artifacts** 下载 `gus-windows-amd64-dev`，解压后运行 `gus.exe`。 |

启动后，点击通知区域中的红色 gus 图标打开面板。图标被收起时，展开任务栏的“显示隐藏的图标”。再次运行程序会唤起已有面板。

## 三步开始使用

### 1. 保存身份

点击主面板右上角的人像按钮，进入 **账号管理 → 新增账号**，填写标签、名字和邮箱。例如，保存“个人”和“工作”两套身份。

也可以右键顶部的全局身份，选择 **保存当前身份为账号…**。账号列表会显示“全局默认”或使用该身份的项目数量。

### 2. 添加项目

点击文件夹加号按钮，选择一个或多个 Git 项目目录。也可以先点击图钉置顶面板，再从资源管理器拖入目录。选择仓库中的子目录时，gus 会自动找到根目录，重复添加不会生成重复记录。

### 3. 应用身份

右键项目行，在菜单中选择已保存的账号。在 **账号管理** 中点击账号右侧的 **设为全局**，或右键顶部的 **全局** 区域，可以更改默认身份；为项目选择 **继承全局**，可以撤销它的单独配置。

例如，将“个人”设为全局默认，再为 `atlas-api` 和 `atlas-web` 选择“工作”，这些项目之后的提交就会使用工作邮箱。

<table>
  <tr>
    <th>保存和管理身份</th>
    <th>右键切换项目身份</th>
  </tr>
  <tr>
    <td><img src="assets/images/gus-accounts.png" width="400" alt="账号管理：个人身份为全局默认，工作身份用于三个项目，开源身份用于一个项目"></td>
    <td><img src="assets/images/gus-switch.png" width="400" alt="项目右键菜单：继承全局，或选择个人、工作、开源身份"></td>
  </tr>
</table>

保存或编辑账号只是更新身份模板，需要在项目上选择该账号才会写入 Git 配置。切换身份不会改写已有提交的作者信息。

## 看懂身份来源

| 标记 | 含义 |
| --- | --- |
| **继承** | 姓名和邮箱均来自全局 Git 配置。 |
| **专用** | 姓名和邮箱均来自仓库本地配置，或均来自工作树专用配置。 |
| **混合** | 两个字段来源不同，例如姓名继承全局、邮箱由项目单独配置。 |
| **未设置** | 没有读取到姓名和邮箱配置。 |
| **不可用** | 项目目录无法访问，或无法读取其 Git 配置。 |

悬停仓库名可查看完整路径与两个字段各自的配置作用域。使用 `git worktree` 时，同一仓库的多个工作树可能共享本地配置。

要在命令行核对身份，可进入对应项目目录运行：

```powershell
git config --show-scope --show-origin --get user.name
git config --show-scope --show-origin --get user.email
```

## 外观与日常使用

点击齿轮进入设置，选择语言和外观，或启用开机自启。启用后，登录 Windows 时会启动到托盘。

<p align="center">
  <img src="assets/images/gus-settings.png" width="480" alt="gus 设置页：语言、外观、开机自启、应用更新和关于">
</p>

截图使用示例账号和仓库；设置页来自本地构建。

- **图钉**：让面板保持置顶，方便拖入项目或连续切换。
- **最小化**：隐藏面板，程序继续留在托盘。
- **关闭**：退出程序，也可以在托盘右键菜单中选择“退出”。
- **调整大小**：拖动边缘或四角，正常退出后会记住面板尺寸。
- **应用更新**：设置页可查看更新日志；启用在线更新的构建可检查新版本，校验签名后安装并重启。

## 数据保存在哪里

账号模板、项目目录列表、语言和外观偏好保存在：

```text
%APPDATA%\GitUserSwitch\state.json
```

实际提交身份保存在 Git 配置中。右键项目选择 **移出目录列表** 只会移除 gus 中的记录，项目文件和 Git 配置会保留。备份应用数据时，可复制上面的 `state.json`。

## 从源码构建

当前版本使用 **Go + MyGo Native UI**，Go 工具链版本见 [go.mod](go.mod)。仓库中的 Visual Studio 工程和 C++ 源码为早期实现。

```powershell
git clone https://github.com/Hunlongyu/gus_gui.git
cd gus_gui

go test ./...
go vet ./...
go run github.com/egoist/mygo/cmd/mygo build -platform windows/amd64 -skip-installer -o dist

.\dist\windows-amd64\gus.exe
```

测试使用临时仓库和隔离的全局 Git 配置。`-skip-installer` 生成可执行程序，适合本地运行。

启动图标、窗口图标和托盘图标共用 [resources/icon.png](resources/icon.png)。MyGo 构建命令会将多尺寸 Windows 图标嵌入 EXE；直接使用 `go build` 不会生成这些图标资源。

<details>
<summary>维护者：发布流程</summary>

日常推送到 `main` 和提交 PR 会触发 [Windows build](.github/workflows/build.yml)，执行测试、静态检查并生成开发构建。

发布正式版本时：

1. 更新 [mygo.json](mygo.json) 中的版本号和 [CHANGELOG.md](CHANGELOG.md)。
2. 创建与版本号一致的稳定版标签，例如 `v1.2.3`。
3. 从 `main` 手动运行 **Publish Windows release**，填写该标签。

[发布工作流](.github/workflows/release.yml) 使用仓库密钥 `MYGO_UPDATER_PRIVATE_KEY` 签名，校验通过后发布安装包及更新清单。

</details>

## 许可证

[MIT](LICENSE.txt)
