param([switch]$ValidateOnly)
$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $true
$config = Get-Content -Raw mygo.json | ConvertFrom-Json
$tag = "v$($config.version)"
$files = @{}
Get-ChildItem -LiteralPath (Join-Path $config.out 'release-assets') -File -Recurse | ForEach-Object {
    if ($_.Length -eq 0) { throw "Empty asset: $($_.Name)" }
    if ($files.ContainsKey($_.Name)) {
        # MyGo generates architecture-independent install scripts in each job.
        if ((Get-FileHash -LiteralPath $files[$_.Name]).Hash -cne (Get-FileHash -LiteralPath $_.FullName).Hash) {
            throw "Conflicting assets: $($_.Name)"
        }
    } else {
        $files[$_.Name] = $_.FullName
    }
}
$slug, $version = $config.slug, $config.version
$notes = $null
$required = @('install.ps1', 'install.sh', "$slug-$version.dmg")
foreach ($platform in @('windows-amd64', 'windows-arm64', 'darwin-universal', 'linux-amd64', 'linux-arm64')) {
    $manifestName = "update-$platform.json"
    if (!$files.ContainsKey($manifestName)) { throw "Missing $manifestName" }
    $manifest = Get-Content -Raw -LiteralPath $files[$manifestName] | ConvertFrom-Json
    if ($manifest.version -cne $version -or [string]::IsNullOrWhiteSpace($manifest.notes) -or !$manifest.signature) {
        throw "Invalid manifest: $manifestName"
    }
    if ($notes -and $notes -cne $manifest.notes) { throw 'Platform release notes differ' }
    $notes = $manifest.notes
    $archive = "$slug-$version-$platform.tar.gz"
    if ($manifest.url -cne "https://github.com/$env:GITHUB_REPOSITORY/releases/download/$tag/$archive") {
        throw "Unexpected update URL in $manifestName"
    }
    if (!$files.ContainsKey($archive) -or (Get-Item -LiteralPath $files[$archive]).Length -ne $manifest.size) {
        throw "Missing or truncated archive: $archive"
    }
    $arch = $platform.Split('-')[1]
    if ($platform.StartsWith('windows-')) {
        $required += "$slug-setup-$version-$arch.exe", "$slug-$version-$platform.exe"
    } elseif ($platform.StartsWith('linux-')) {
        $required += "${slug}_${version}_${arch}.deb", "$slug-$version-$platform", "PKGBUILD-$platform"
    }
}
foreach ($name in $required) {
    if (!$files.ContainsKey($name)) { throw "Missing asset: $name" }
}
$checksumFile = Join-Path $config.out 'SHA256SUMS'
$checksums = foreach ($name in ($files.Keys | Sort-Object -CaseSensitive)) {
    "{0}  {1}" -f (Get-FileHash -LiteralPath $files[$name] -Algorithm SHA256).Hash.ToLowerInvariant(), $name
}
Set-Content -LiteralPath $checksumFile -Value $checksums -Encoding utf8NoBOM
$assets = @($files.Values) + @((Get-Item -LiteralPath $checksumFile).FullName)
if ($ValidateOnly) {
    Write-Host "Verified $($assets.Count) assets across all five release targets"
    return
}
$notesFile = Join-Path $env:RUNNER_TEMP 'gus-release-notes.md'
$notes += @"


### 下载 / Downloads

| 平台 / Platform | 单文件或应用 / Standalone app | 安装包 / Installer |
| --- | --- | --- |
| Windows x64 / ARM64 | ``$slug-$version-windows-<arch>.exe`` | ``$slug-setup-$version-<arch>.exe`` |
| macOS Intel + Apple Silicon | ``$slug-$version-darwin-universal.tar.gz``（内含 .app） | ``$slug-$version.dmg`` |
| Linux x64 / ARM64 | ``$slug-$version-linux-<arch>`` | ``${slug}_${version}_<arch>.deb`` |

单文件版无需旁边的资源文件。Linux 下载后需添加执行权限，并安装 GTK 3 和 Git；托盘需要 AppIndicator 支持。macOS 使用完整 .app 包，当前为 ad-hoc 签名，未经过 Apple 公证。所有平台均需要 Git。

Standalone Windows/Linux executables need no resource sidecars. Linux requires executable permission, GTK 3 and Git; tray support additionally needs AppIndicator. macOS requires the complete .app bundle; it is ad-hoc signed, not Apple-notarized. Git is required on every platform.

保留下载文件名的单文件版会引导到发布页下载更新；安装版支持应用内签名更新。/ Standalone downloads keep updates manual; installed apps support signed in-app updates.
"@
Set-Content -LiteralPath $notesFile -Value $notes -Encoding utf8NoBOM
# A retry can reuse an incomplete draft. Never mutate a published release.
$releases = gh api --paginate "repos/$env:GITHUB_REPOSITORY/releases" --jq ".[] | select(.tag_name == `"$tag`") | .draft"
if ($releases -eq 'false') { throw "$tag has already been published" }
if ($releases -ne 'true') {
    gh release create $tag --draft --verify-tag --title $tag --notes-file $notesFile
}
gh release upload $tag @assets --clobber
gh release edit $tag --notes-file $notesFile --draft=false --latest
