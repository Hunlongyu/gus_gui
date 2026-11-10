$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $true
$config = Get-Content -Raw mygo.json | ConvertFrom-Json
$tag = "v$($config.version)"
$dir = Join-Path $config.out 'windows-amd64'
$manifestFile = Get-Item (Join-Path $dir 'update-windows-amd64.json')
$manifest = Get-Content -Raw $manifestFile.FullName | ConvertFrom-Json
if ($manifest.version -cne $config.version -or [string]::IsNullOrWhiteSpace($manifest.notes)) {
    throw 'Manifest version or release notes missing'
}
$installer = Get-Item (Join-Path $dir "$($config.slug)-setup-$($config.version)-amd64.exe")
$archive = Get-Item (Join-Path $dir "$($config.slug)-$($config.version)-windows-amd64.tar.gz")
$assets = @($installer.FullName, $archive.FullName, $manifestFile.FullName)
$assets += (Get-Item (Join-Path $dir 'install.ps1')).FullName
$assets += (Get-Item (Join-Path $config.out 'SHA256SUMS')).FullName
$assets += @(Get-ChildItem -Path $dir -Filter '*.delta' -File | ForEach-Object FullName)
$notesFile = Join-Path $env:RUNNER_TEMP 'gus-release-notes.md'
Set-Content -LiteralPath $notesFile -Value $manifest.notes -Encoding utf8NoBOM
# A retry can reuse an incomplete draft. Never mutate a published release.
$releases = gh api --paginate "repos/$env:GITHUB_REPOSITORY/releases" --jq ".[] | select(.tag_name == `"$tag`") | .draft"
if ($releases -eq 'false') { throw "$tag has already been published" }
if ($releases -ne 'true') {
    gh release create $tag --draft --verify-tag --title $tag --notes-file $notesFile
}
gh release upload $tag @assets --clobber
gh release edit $tag --notes-file $notesFile --draft=false --latest
