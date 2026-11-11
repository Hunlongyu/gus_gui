param([Parameter(Mandatory)][string]$Platform, [string]$BuildRoot = 'dist')
$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $true
# Run the packaged executable in an otherwise empty directory, away from
# source resources and real user configuration. macOS requires its app bundle.
$root = Join-Path $env:RUNNER_TEMP "gus-smoke-$Platform"
New-Item -ItemType Directory -Force $root | Out-Null
$env:HOME = Join-Path $root 'home'
$env:APPDATA = Join-Path $root 'config'
$env:XDG_CONFIG_HOME = $env:APPDATA
$env:XDG_DATA_HOME = Join-Path $root 'data'
$env:GIT_CONFIG_NOSYSTEM = '1'
$env:GIT_CONFIG_GLOBAL = Join-Path $root 'empty.gitconfig'
New-Item -ItemType Directory -Force $env:HOME, $env:APPDATA, $env:XDG_DATA_HOME | Out-Null
Set-Content -LiteralPath $env:GIT_CONFIG_GLOBAL -Value ''
if ($IsMacOS) {
    Copy-Item -LiteralPath "$BuildRoot/$Platform/gus.app" -Destination $root -Recurse
    $exe = Join-Path $root 'gus.app/Contents/MacOS/gus'
    codesign --verify --deep --strict "$root/gus.app"
} else {
    $name = if ($IsWindows) { 'gus.exe' } else { 'gus' }
    Copy-Item -LiteralPath "$BuildRoot/$Platform/$name" -Destination $root
    $exe = Join-Path $root $name
    if (!$IsWindows) { chmod +x $exe }
}
$log = Join-Path $root 'stderr.log'
$options = @{FilePath = $exe; WorkingDirectory = $root; PassThru = $true; RedirectStandardError = $log}
if ($IsWindows) { $options.WindowStyle = 'Hidden' }
$process = Start-Process @options
try {
    if ($process.WaitForExit(8000)) {
        Get-Content -LiteralPath $log
        throw "Application exited during startup: $($process.ExitCode)"
    }
    if ($IsWindows) {
        # A second invocation should focus the existing tray panel and exit.
        $second = Start-Process -FilePath $exe -WorkingDirectory $root -PassThru -WindowStyle Hidden
        if (!$second.WaitForExit(8000) -or $second.ExitCode -ne 0) {
            if (!$second.HasExited) { Stop-Process -Id $second.Id -Force }
            throw 'Single-instance activation failed'
        }
    }
    Write-Host "Native startup passed: $Platform"
} finally {
    if (!$process.HasExited) { Stop-Process -Id $process.Id -Force }
    if (Test-Path -LiteralPath $log) { Get-Content -LiteralPath $log }
}
