param(
    [string]$BuildDir = "build",
    [string]$BinaryName = "dashboard"
)

# Clean previous build
Write-Host "Cleaning previous build..."
@("windows_amd64", "linux_amd64") | ForEach-Object {
    $path = "$BuildDir/$_"
    if (Test-Path $path) {
        Write-Host "  Removing $path"
        Remove-Item -Path $path -Recurse -Force
    }
}

# Create release structure
$releasePath = "$BuildDir/windows_amd64"
$binPath = "$releasePath/bin"
$docPath = "$releasePath/docs"

Write-Host "Creating release structure..."
New-Item -ItemType Directory -Force -Path $binPath
New-Item -ItemType Directory -Force -Path $docPath
New-Item -ItemType Directory -Force -Path "$BuildDir/linux_amd64"

# Set up Windows build environment
$env:CGO_ENABLED = '1'
$env:CC = 'gcc'
$env:GOOS = 'windows'
$env:GOARCH = 'amd64'
$env:PATH = "c:\mingw64\bin;$env:PATH"

# Build Windows binary
Write-Host "Building Windows binary..."
go build -o "$binPath/$BinaryName.exe" ./cmd/dashboard

# Copy release documentation
Write-Host "Copying release documentation..."

# Core files (stay in root)
Copy-Item -Path "README.md", "LICENSE" -Destination $releasePath -Force

# Documentation (goes in docs/)
Get-ChildItem -Path "docs/*.md" | ForEach-Object {
    Write-Host "  Copying $($_.Name)"
    Copy-Item -Path $_.FullName -Destination $docPath -Force
}

# Documentation assets
if (Test-Path "docs/images") {
    Write-Host "  Copying documentation images"
    Copy-Item -Path "docs/images" -Destination "$docPath/images" -Force -Recurse
}

Write-Host "Release structure created at $releasePath"
Write-Host "  bin/  - Binary files"
Write-Host "  docs/ - Documentation" 