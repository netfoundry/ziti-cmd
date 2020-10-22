suffix=$(date +"%b-%d-%H%M%S")
idname="User${suffix}"

ziti edge controller login "${ZITI_EDGE_API_HOSTNAME}" -u "${ZITI_USER}" -p "${ZITI_PWD}" -c "${ZITI_PKI}/${ZITI_EDGE_ROOTCA_NAME}/certs/${ZITI_EDGE_INTERMEDIATE_NAME}.cert"

ziti edge controller delete identity "${idname}"
ziti edge controller create identity device "${idname}" -o "${ZITI_HOME}/${idname}.jwt"

cp "${ZITI_HOME}/${idname}.jwt" /mnt/v/temp/ziti-windows-tunneler/${idname}.jwt
echo "jwt written to: /mnt/v/temp/ziti-windows-tunneler/${idname}.jwt"
export NEW_ID_FILE="/mnt/v/temp/ziti-windows-tunneler/${idname}.jwt"
echo "              : $(wslpath -w "/mnt/v/temp/ziti-windows-tunneler/${idname}.jwt")"
