SVC_NAME=$1
SVC_HOST=$2
SVC_PORT=$3
SVR_PROTO=$4
SVR_HOST=$5
SVR_PORT=$6

ziti edge controller delete service "${SVC_NAME}svc"
ziti edge controller delete config "${SVC_NAME}-clientconfig"
ziti edge controller delete config "${SVC_NAME}-serverconfig"

ziti edge controller create config "${SVC_NAME}-clientconfig" ziti-tunneler-client.v1 '{ "hostname" : "'"${SVC_HOST}"'", "port" : '"${SVC_PORT}"' }'
ziti edge controller create service "${SVC_NAME}svc" --configs "${SVC_NAME}-clientconfig"
ziti edge controller create config netcat-udp-server ziti-tunneler-server.v1 '{"protocol":"'"${SVR_PROTO}"'","hostname":"'"${SVR_HOST}"'","port":'"${SVR_PORT}"'}'
ziti edge controller create service netcat-udp -c "${SVC_NAME}-clientconfig,{SVC_NAME}-serverconfig"
