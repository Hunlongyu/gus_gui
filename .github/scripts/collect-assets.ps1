param([Parameter(Mandatory)][string]$Platform)
$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $true
$config = Get-Content -Raw mygo.json | ConvertFrom-Json
$dir = Join-Path $config.out $Platform
$dest = Join-Path $config.out "release-assets/$Platform"
New-Item -ItemType Directory -Force $dest | Out-Null
$slug, $version = $config.slug, $config.version
$arch = $Platform.Split('-')[1]
$files = @("$slug-$version-$Platform.tar.gz", "update-$Platform.json")
if ($Platform.StartsWith('windows-')) {
    $files += "$slug-setup-$version-$arch.exe", 'install.ps1'
    Copy-Item -LiteralPath "$dir/$slug.exe" -Destination "$dest/$slug-$version-$Platform.exe"
} elseif ($Platform.StartsWith('linux-')) {
    $files += "${slug}_${version}_${arch}.deb", 'install.sh'
    Copy-Item -LiteralPath "$dir/$slug" -Destination "$dest/$slug-$version-$Platform"
    Copy-Item -LiteralPath "$dir/PKGBUILD" -Destination "$dest/PKGBUILD-$Platform"
} elseif ($Platform -eq 'darwin-universal') {
    $files += "$slug-$version.dmg", 'install.sh'
} else {
    throw "Unsupported release platform: $Platform"
}
foreach ($file in $files) {
    $item = Get-Item -LiteralPath (Join-Path $dir $file)
    if ($item.Length -eq 0) { throw "Empty release asset: $file" }
    Copy-Item -LiteralPath $item.FullName -Destination $dest
}
Get-ChildItem -LiteralPath $dir -Filter '*.delta' -File | Copy-Item -Destination $dest
