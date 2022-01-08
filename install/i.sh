#!/usr/bin/env bash
# coding: utf-8
#  ---------------------------------------------------
#  File          : i.sh
#  Authors       : ccmywish <ccmywish@qq.com>
#  Created on    : <2022-1-8>
#  Last modified : <2022-1-8>
#
#  Install 'cr' on Linux or macOS
#  ---------------------------------------------------

echo "Downloading binary from github.com/cryptic-resolver/cr_Go"

ostype=$(uname) 
$cr_ver=1.3.1

dlprefix="https://github.com/cryptic-resolver/cr_Go/releases/download/v${dlprefix}${cr_ver}/cr-${cr_ver}-"

if [[ $ostype == 'linux' ]]; then
   dl="${dlprefix}amd64-unknown-linux"
else then
   dl="${dlprefix}arm64-apple-darwin"
fi

curl -fSL ${dl}

mv ./cr-1.3.1-amd64-unknown-linux ./cr

chmod +x ./cr

echo "move to /usr/local/bin"
sudo mv ./cr /usr/local/bin

echo "Finish! Try 'cr emacs' now"