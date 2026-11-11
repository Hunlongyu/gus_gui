### 下载 / Downloads

**根据你的系统下载 / Download based on your OS:**

<table>
  <thead>
    <tr><th align="left">OS</th><th align="left">Download</th></tr>
  </thead>
  <tbody>
    <tr>
      <td>Windows</td>
      <td>
        <a href="https://github.com/{{repository}}/releases/download/v{{version}}/{{slug}}-setup-{{version}}-amd64.exe"><img alt="Windows Setup x64" src="https://img.shields.io/badge/Setup-x64-2d7d9a.svg"></a><br>
        <a href="https://github.com/{{repository}}/releases/download/v{{version}}/{{slug}}-setup-{{version}}-arm64.exe"><img alt="Windows Setup ARM64" src="https://img.shields.io/badge/Setup-ARM64-2d7d9a.svg"></a><br>
        <a href="https://github.com/{{repository}}/releases/download/v{{version}}/{{slug}}-{{version}}-windows-amd64.exe"><img alt="Windows Portable x64" src="https://img.shields.io/badge/Portable-x64-67b7d1.svg"></a><br>
        <a href="https://github.com/{{repository}}/releases/download/v{{version}}/{{slug}}-{{version}}-windows-arm64.exe"><img alt="Windows Portable ARM64" src="https://img.shields.io/badge/Portable-ARM64-67b7d1.svg"></a>
      </td>
    </tr>
    <tr>
      <td>macOS</td>
      <td>
        <a href="https://github.com/{{repository}}/releases/download/v{{version}}/{{slug}}-{{version}}.dmg"><img alt="macOS DMG Universal" src="https://img.shields.io/badge/DMG-Universal-000000.svg?logo=apple&amp;logoColor=white"></a><br>
        <a href="https://github.com/{{repository}}/releases/download/v{{version}}/{{slug}}-{{version}}-darwin-universal.tar.gz"><img alt="macOS APP Universal" src="https://img.shields.io/badge/APP-Universal-00A9E0.svg?logo=apple&amp;logoColor=white"></a>
      </td>
    </tr>
    <tr>
      <td>Linux</td>
      <td>
        <a href="https://github.com/{{repository}}/releases/download/v{{version}}/{{slug}}_{{version}}_amd64.deb"><img alt="Linux DebPackage x64" src="https://img.shields.io/badge/DebPackage-x64-FF9966.svg?logo=debian&amp;logoColor=white"></a><br>
        <a href="https://github.com/{{repository}}/releases/download/v{{version}}/{{slug}}_{{version}}_arm64.deb"><img alt="Linux DebPackage ARM64" src="https://img.shields.io/badge/DebPackage-ARM64-FF9966.svg?logo=debian&amp;logoColor=white"></a><br>
        <a href="https://github.com/{{repository}}/releases/download/v{{version}}/{{slug}}-{{version}}-linux-amd64"><img alt="Linux Portable x64" src="https://img.shields.io/badge/Portable-x64-f84e29.svg?logo=linux&amp;logoColor=white"></a><br>
        <a href="https://github.com/{{repository}}/releases/download/v{{version}}/{{slug}}-{{version}}-linux-arm64"><img alt="Linux Portable ARM64" src="https://img.shields.io/badge/Portable-ARM64-f84e29.svg?logo=linux&amp;logoColor=white"></a><br>
        <a href="https://github.com/{{repository}}/releases/download/v{{version}}/{{slug}}-{{version}}-linux-amd64.tar.gz"><img alt="Linux Tar.gz x64" src="https://img.shields.io/badge/Tar.gz-x64-F1B42F.svg?logo=linux&amp;logoColor=white"></a><br>
        <a href="https://github.com/{{repository}}/releases/download/v{{version}}/{{slug}}-{{version}}-linux-arm64.tar.gz"><img alt="Linux Tar.gz ARM64" src="https://img.shields.io/badge/Tar.gz-ARM64-F1B42F.svg?logo=linux&amp;logoColor=white"></a>
      </td>
    </tr>
  </tbody>
</table>

**Portable**：单文件 / standalone executable。**Universal**：同时支持 Apple Silicon 和 Intel x64 / supports both Apple Silicon and Intel x64。

[SHA-256 校验值 / Checksums](https://github.com/{{repository}}/releases/download/v{{version}}/SHA256SUMS)

单文件版无需旁边的资源文件。Linux 下载后需添加执行权限，并安装 GTK 3 和 Git；托盘需要 AppIndicator 支持。macOS 使用完整 .app 包，当前为 ad-hoc 签名，未经过 Apple 公证。所有平台均需要 Git。

Standalone Windows/Linux executables need no resource sidecars. Linux requires executable permission, GTK 3 and Git; tray support additionally needs AppIndicator. macOS requires the complete .app bundle; it is ad-hoc signed, not Apple-notarized. Git is required on every platform.

保留下载文件名的单文件版会引导到发布页下载更新；安装版支持应用内签名更新。/ Standalone downloads keep updates manual; installed apps support signed in-app updates.
