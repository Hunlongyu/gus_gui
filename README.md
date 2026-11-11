<p align="center">
  <img width="96" src="resources/icon.png" alt="gus Logo">
</p>

<h1 align="center">gus</h1>

<p align="center">
  Manage Git commit identities across all your repositories.<br>
  <b>Native Desktop UI · Standalone Portable Binary · Windows, macOS & Linux</b>
</p>

<p align="center">
  <a href="https://github.com/Hunlongyu/gus_gui/releases/latest"><img alt="GitHub Release" src="https://img.shields.io/github/v/release/Hunlongyu/gus_gui?include_prereleases&style=flat-square&logo=github"></a>
  <a href="https://github.com/Hunlongyu/gus_gui/releases"><img alt="GitHub Downloads" src="https://img.shields.io/github/downloads/Hunlongyu/gus_gui/total?style=flat-square&color=E1384B"></a>
  <a href="https://github.com/Hunlongyu/gus_gui/releases/latest"><img alt="Portable" src="https://img.shields.io/badge/portable-Single%20Binary%20%2F%20Installer-2ea44f?style=flat-square"></a>
  <img alt="Platform" src="https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-0078D4?style=flat-square">
  <a href="LICENSE.txt"><img alt="License" src="https://img.shields.io/github/license/Hunlongyu/gus_gui?style=flat-square"></a>
</p>

<p align="center">
  <b>English</b> · <a href="README.zh-CN.md">简体中文</a>
</p>

<p align="center">
  <img src="assets/images/gus-panel-en.png" width="480" alt="gus Main Interface">
</p>

<details>
<summary>📸 <b>More Screenshots</b> (Accounts, Identity Switch & Settings)</summary>
<br>
<table align="center">
  <tr>
    <th>Account Profiles</th>
    <th>Right-Click Identity Switch</th>
  </tr>
  <tr>
    <td><img src="assets/images/gus-accounts-en.png" width="400" alt="gus Accounts"></td>
    <td><img src="assets/images/gus-switch-en.png" width="400" alt="gus Identity Switch Menu"></td>
  </tr>
  <tr>
    <th colspan="2">Settings</th>
  </tr>
  <tr>
    <td colspan="2" align="center"><img src="assets/images/gus-settings-en.png" width="400" alt="gus Settings"></td>
  </tr>
</table>
</details>

## ✨ Features

- **Multi-Account Profiles** — Save Git identities (`user.name` & `user.email`) for Personal, Work, or Open Source projects and track how many repositories use each profile.
- **One-Click Identity Switching** — Right-click any repository to assign a dedicated identity, switch accounts, or revert to inheriting your global default.
- **Visual Audit at a Glance** — Inspect repository name, active author, email, push hosting platforms (GitHub, GitLab, Gitea, Codeberg, Bitbucket, Gitee, etc.), and identity scope in one unified panel.
- **Instant Startup & Background Refresh** — Displays cached repository states immediately on launch and refreshes projects concurrently in the background.
- **Tray-First Workflow** — Runs quietly in the system tray with Always-on-Top pin mode, drag-and-drop folder import, instant search, Light/Dark themes, and English / Simplified Chinese UI.
- **Native & Portable** — Built with native desktop UI — no Electron or WebView2 bloat. Available as a standalone single-file executable (no installation required) or native installer.

## 📦 Download

Get the latest release for your platform from **[GitHub Releases](https://github.com/Hunlongyu/gus_gui/releases/latest)**:

| Platform | Standalone / Portable | Installer / Package |
| :--- | :--- | :--- |
| **Windows** *(x64 / ARM64)* | `gus-<version>-windows-<arch>.exe` *(Single EXE, no install)* | `gus-setup-<version>-<arch>.exe` |
| **macOS** *(Universal)* | `gus-<version>-darwin-universal.tar.gz` | `gus-<version>.dmg` |
| **Linux** *(x64 / ARM64)* | `gus-<version>-linux-<arch>` *(Single binary)* | `gus_<version>_<arch>.deb` |

> **Note:** Requires `git` in your system `PATH`. Windows and Linux standalone binaries have all icons and assets embedded — just download and run.

## ⌨️ Controls & Indicators

| Action / Indicator | Description |
| :--- | :--- |
| **Right-click Project** | Switch repository identity, inherit global default, or remove from list |
| **Right-click Global Bar** | Switch global default Git identity or save current global config as an account |
| **Double-click Project** | Inspect current branch, fetch/push remote URLs, and open repository in browser |
| **Drag & Drop Folders** | Add one or more Git repositories (subfolders automatically resolve to repo root) |
| **Global** *(Green Icon)* | Both `user.name` and `user.email` are inherited from global Git config |
| **Local** *(Purple Icon)* | Identity is configured specifically for this repository |
| **Mixed** *(Amber Icon)* | Name and email come from different Git config scopes (e.g., global name + local email) |

## 📄 License

Released under the [MIT License](LICENSE.txt).
