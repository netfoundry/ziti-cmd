SVC_NAME=$1
SVC_HOST=$2
SVC_PORT=$3
TCP_HOST_PORT=$4

ziti edge controller delete service "${SVC_NAME}svc"
ziti edge controller delete config "${SVC_NAME}svcconfig"
ziti edge controller create config "${SVC_NAME}svcconfig" ziti-tunneler-client.v1 '{ "hostname" : "'"${SVC_HOST}"'", "port" : '"${SVC_PORT}"' }'
ziti edge controller create service "${SVC_NAME}svc" --configs "${SVC_NAME}svcconfig"
ziti edge controller create terminator "${SVC_NAME}svc" "${ZITI_EDGE_ROUTER_NAME}" $TCP_HOST_PORT
