# 更新日志

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
