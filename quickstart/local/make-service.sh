ziti edge controller delete service "${svc}svc"
ziti edge controller delete config "${svc}-client-config"
ziti edge controller delete config "${svc}-server-config"

if [ ! -z "$sproto" ]
then
    echo HOSTED
    ziti edge controller create config "${svc}-server-config" ziti-tunneler-server.v1 "{ \"protocol\":\"${sproto}\", \"hostname\":\"${shost}\", \"port\":${sport} }"
fi
ziti edge controller create config "${svc}-client-config" ziti-tunneler-client.v1 "{ \"hostname\":\"${ihost}\", \"port\":${iport} }"
ziti edge controller create service "${svc}svc" -c "${svc}-server-config","${svc}-client-config"

if [ ! -z "$to" ]
then
    echo ROUTER TERMINATED
    ziti edge controller create terminator "${svc}svc" "${ZITI_EDGE_ROUTER_NAME}" "${to}"
fi



