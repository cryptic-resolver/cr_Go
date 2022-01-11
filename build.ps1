#   ---------------------------------------------------
#   File          : build.ps1
#   Authors       : ccmywish <ccmywish@qq.com>
#   Created on    : <2022-1-7>
#   Last modified : <2022-1-11>
#
#   Build cr on multi-platform via PowerShell
#   ---------------------------------------------------


# Notice the dot sign can only placed at the bottom
$version = 
    (Get-Content .\cr.go | Select-String "const CRYPTIC_VERSION" ).
    ToString().
    TrimStart("const CRYPTIC_VERSION = ").
    Trim("`"")  # should remove double quotes

Write-Host "cr version: $version"

Write-Host "Building for Windows x64"
$env:GOOS="windows"; $env:GOARCH="amd64" ; go build -o "build/cr-${version}-amd64-pc-windows.exe"

Write-Host "Building for Linux x64"
$env:GOOS="linux"; $env:GOARCH="amd64" ; go build -o "build/cr-${version}-amd64-unknown-linux"

Write-Host "Building for macOS x64"
$env:GOOS="darwin"; $env:GOARCH="arm64" ; go build -o "build/cr-${version}-arm64-apple-darwin"

# Auto generate scoop manifest
$windows_bin_sha256 = (Get-FileHash "build/cr-${version}-amd64-pc-windows.exe").hash

Write-Host "Windows bin SHA256: $windows_bin_sha256"


$scoop_manifest = '{
    "version": "' + "$version" + '",
    "description": "Cryptic-Resolver (cr) is a fast command line tool used to record and explain cryptic commands, acronyms and so forth in every field, including your own knowledge base.",
    "homepage": "https://github.com/cryptic-resolver/cr_Go",
    "license": "MIT",
    "architecture": {
        "64bit": {
            "url": "https://github.com/cryptic-resolver/cr_Go/releases/download/v' + "$version" + '/cr-' + "$version" + '-amd64-pc-windows.exe",
            "hash": "' + "$windows_bin_sha256" + '"
        }
    },
    "bin": [ ["cr-' + "$version" + '-amd64-pc-windows.exe","cr"] ] ,
    "checkver": "github",
    "autoupdate": {
        "architecture": {
            "64bit": {
                "url": "https://github.com/cryptic-resolver/cr_Go/releases/download/v$version/cr-$version-amd64-pc-windows.exe"
            }
        }
    }
}

'

Set-Content -Path "install/cryptic-resolver.json" -Value $scoop_manifest
Write-Host "Generate cryptic-resolver.json in ./build/"


$nix_install =  (Get-Content -Path "install/i-template.sh").Replace("cr_ver=`"1.3.1`"","cr_ver=`"${version}`"")
Set-Content -Path "install/i.sh" -Value $nix_install
Write-Host "Generate i.sh in ./build/"
