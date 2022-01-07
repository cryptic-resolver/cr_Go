$version = "1.2.0"

Write-Host "Building for Windows x64"
$GOOS="windows"; $GOARCH="amd64" ; go build -o "build/cr-${version}-windows.exe"

Write-Host "Building for Linux x64"
$GOOS="linux"; $GOARCH="amd64" ; go build -o "build/cr-${version}-linux"

# Write-Host "Building for macOS x64"
# $GOOS="linux"; $GOARCH="arm64" ; go build -o "build/cr-${version}-macOS"
