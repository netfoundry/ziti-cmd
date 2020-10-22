CURDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

ziti edge controller login "${ZITI_EDGE_API_HOSTNAME}" -u "${ZITI_USER}" -p "${ZITI_PWD}" -c "${ZITI_PKI}/${ZITI_EDGE_ROOTCA_NAME}/certs/${ZITI_EDGE_INTERMEDIATE_NAME}.cert"

svc=clint-ssh-hosted-dns ihost=clint.ssh.hosted iport=22 sproto=tcp shost=wsl2-controller sport=2200 $CURDIR/make-service.sh
svc=clint-ssh-hosted-ip ihost=192.168.15.3 iport=22 sproto=tcp shost=wsl2-controller sport=2200 $CURDIR/make-service.sh

svc=clint-ssh-router-dns ihost=clint.ssh.router iport=22 to=tcp:localhost:2200 $CURDIR/make-service.sh
svc=clint-ssh-router-ip ihost=192.168.15.4 iport=22 to=tcp:localhost:2200 $CURDIR/make-service.sh


ziti edge controller delete service-policy dial-all
ziti edge controller create service-policy dial-all Dial --service-roles '#all' --identity-roles '#all'

ziti edge controller delete service-policy bind-all
ziti edge controller create service-policy bind-all Bind --service-roles '#all' --identity-roles '#all'
#ziti-enroller --jwt "${ZITI_HOME}/test_identity.jwt" -o "${ZITI_HOME}/test_identity".json

#ziti-tunnel proxy netcatsvc:8145 -i "${ZITI_HOME}/test_identity".json > "${ZITI_HOME}/ziti-test_identity.log" 2>&1 &
cp "${ZITI_HOME}/test_identity.jwt" /mnt/v/temp/identities/_new_id.jwt



