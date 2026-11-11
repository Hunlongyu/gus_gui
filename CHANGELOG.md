# 更新日志

## 1.3.0

### 简体中文

- 项目列表支持多个远程仓库和推送地址，以堆叠平台图标显示，悬浮可查看完整平台列表。
- 双击项目可查看账号、分支及拉取和推送地址，支持用默认浏览器打开仓库链接。
- 账号管理新增“设为全局”按钮，设置成功后显示“全局默认”。
- 列宽根据内容和窗口宽度自动分配；窗口尺寸在调整后自动保存，重启后恢复。
- 主面板操作提示改为悬浮通知；“提交姓名”统一改为“名字”，移除名字和邮箱各自的悬浮提示。
- “关于”新增开发者、仓库和反馈文字链接。
- 补全托管平台图标，去除应用图标外圈透明留白，并在 Windows 构建中嵌入多尺寸图标。

### English

- Show multiple Git remotes and push destinations as stacked platform icons, with a full platform list on hover.
- Double-click a project to inspect its account, branch, fetch URLs and push URLs, and open repository links in the default browser.
- Set a saved account as the global Git identity directly from Accounts, with a Global default indicator.
- Allocate column widths according to their content and window width. Save window dimensions after resizing and restore them on restart.
- Show main-panel operation feedback as floating notifications. Rename Git name to Name and remove separate name and email tooltips.
- Add developer, repository and feedback text links to About.
- Complete hosting-platform icons, remove excess transparent padding from the app icon and embed multiple Windows icon sizes in builds.

## 1.2.3

### 简体中文

- 应用统一命名为 gus，Windows 主程序为 `gus.exe`。托盘提示、窗口标题、关于页和新安装程序使用同一名称。
- 从旧版更新后会自动启动 gus，并继续使用已保存的账号和项目目录。

### English

- The app is now named gus. Its Windows executable is `gus.exe`, with matching tray, window, About and installer names.
- Upgrading an existing installation relaunches gus and keeps saved accounts and project folders.

## 1.2.2

### 简体中文

- 拖动面板边缘或四角即可调整窗口大小。初始尺寸为 480 × 490，最小尺寸为 440 × 320，正常退出后会记住上次的大小。
- 小窗口中的账号表单支持滚动，保存和取消按钮始终可见。

### English

- Resize the tray panel by dragging its edges or corners. It opens at 480 × 490 by default, has a minimum size of 440 × 320 and remembers its size after a normal exit.
- Account fields scroll in smaller windows so Save and Cancel remain accessible.

## 1.2.1

### 简体中文

- 同一 Windows 登录会话中，安装版和独立 EXE 共用一个托盘实例。再次启动会打开已有面板，并发启动也只保留一个进程。
- 正常退出、进程异常结束和自动更新重启后均可重新启动。

### English

- The installer and standalone executable share one tray instance per Windows sign-in session. Launching again opens the existing panel, including when several copies start together.
- The instance lock is released on exit, process termination and update relaunch.

## 1.2.0

### 简体中文

- 设置中可即时切换简体中文和 English，以及跟随系统、浅色和深色外观。重启和升级后保留选择。
- 支持 Windows 路径的 8.3 短名称，避免将同一 Git 仓库识别为不同项目。
- Windows x64 安装版支持后台检查、更新日志、签名校验和安装后自动重启。

### English

- Switch between English and Simplified Chinese and choose System, Light or Dark appearance in Settings. Preferences survive restarts and updates.
- Recognize Windows short and long paths as the same Git repository.
- The Windows x64 installer supports background update checks, release notes, signed downloads and automatic restart after installation.

## 1.1.0

### 简体中文

- Windows 托盘版支持添加和拖入真实仓库、按项目名称搜索、切换全局与项目提交身份。
- 账号模板和项目清单可在重启后恢复，支持识别姓名与邮箱来自不同 Git 配置层的情况。
- 新增后台版本检查、应用内更新日志和签名校验。发现新版本后，可在设置中下载安装并重启。
- 首次从本地开发版升级，请使用 GitHub Release 中的 Windows 安装程序。

### English

- Manage Git workspaces from the Windows tray, add folders by drag and drop, search projects and switch commit identities.
- Account profiles and project folders survive restarts; the panel identifies mixed configuration scopes.
- Check for updates in the background, review release notes and restart after a verified download.
- Install the GitHub Release package once when upgrading from a local development build.
