# gus

使用 MyGo Native UI 的 Windows 托盘工具。运行环境需要可在 PATH 中找到的 Git for Windows。已在 Git 2.56 上验证。

从 [GitHub Releases](https://github.com/Hunlongyu/gus_gui/releases/latest) 下载 Windows x64 安装程序。在设置中可以查看更新日志，下载经签名校验的新版并重启。以前通过普通 `go build` 分发的版本需手动安装一次。

## 使用

设置页可选择简体中文或 English，以及跟随系统、浅色或深色主题。偏好保存在当前用户的配置中。

1. 运行 `gus.exe`，点击系统托盘图标打开面板。
2. 点击人像图标添加账号，填写标签、提交姓名和邮箱。也可右键已有的 Git 身份，选择「保存当前身份为账号」。
3. 点击文件夹加号，或将一个或多个目录从资源管理器拖到项目面板。拖放前可点击图钉置顶面板，方便切换到资源管理器。子目录会定位到仓库根目录，重复目录只保留一项。
4. 右键顶部的全局身份或项目行，选择要使用的账号。项目行的「继承全局」会移除当前仓库配置文件中的姓名和邮箱覆盖。

账号模板和目录清单保存在 `%APPDATA%\GitUserSwitch\state.json`。每次打开面板，应用重新从 Git 读取身份。编辑账号模板后，可以在目标项目上再次选择它来应用修改。

项目同时使用不同作用域的姓名和邮箱时显示「混合」。将指针停在仓库名上可查看两项配置各自的来源。链接工作树遵循 Git 的配置作用域：仓库本地配置可由多个工作树共享，已有的工作树专用身份写入对应的工作树配置。

## 构建与测试

项目根目录为 MyGo 原生版。仓库中的 Visual Studio 工程为早期 C++ 实现。

```powershell
go test ./...
go build -trimpath -ldflags "-H windowsgui" -o dist/gus.exe .
```

测试使用临时仓库和隔离的全局 Git 配置。

## 发布

修改 `mygo.json` 的 `version`，在 `CHANGELOG.md` 添加对应的 `## x.y.z` 段落。提交并推送后，推送相同版本的 Git 标签，例如 `v1.1.1`。推送到 main 和 PR 会运行 Windows 测试与构建。稳定版标签会触发签名打包，在更新归档通过公钥验证后，将安装器、更新清单和发布说明一起发布到 GitHub Release。

Actions 使用仓库密钥 `MYGO_UPDATER_PRIVATE_KEY`，客户端信任的公钥在 `mygo.json`。请保管原签名私钥用于后续版本。应用启动时检查更新，长时间驻留后再次打开也会检查；用户可以先阅读新版说明，再选择「安装并重启」。
