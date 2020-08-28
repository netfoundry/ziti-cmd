nw=wsl2
echo rm -r "~/.ziti/quickstart/${nw}/!(pki)"
rm -r "~/.ziti/quickstart/${nw}/!(pki)"
./init.sh ~/git/github/openziti/nf/ziti/build ${nw}
