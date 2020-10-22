#!/bin/bash

#need to enable glob or the rm below fails
shopt -s extglob

nw=wsl2
echo rm -r "~/.ziti/quickstart/${nw}/!(pki)"
rm -r ~/.ziti/quickstart/${nw}/!(pki)
./init.sh ~/git/github/openziti/nf/ziti/linux-build ${nw}
killall ziti-router ziti-controller
echo "==============================="
