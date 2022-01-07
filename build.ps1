$version = "1.3.0"

Write-Host "Building for Windows x64"
$GOOS="windows"; $GOARCH="amd64" ; go build -o "build/cr-${version}-amd64-pc-windows.exe"

Write-Host "Building for Linux x64"
$GOOS="linux"; $GOARCH="amd64" ; go build -o "build/cr-${version}-amd64-unknown-linux"

# Write-Host "Building for macOS x64"
# $GOOS="darwin"; $GOARCH="arm64" ; go build -o "build/cr-${version}-arm64-apple-darwin"
