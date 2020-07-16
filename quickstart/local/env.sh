#!/bin/bash
export PFXLOG_NO_JSON=true

# make the quickstart home folder where all the config files, logs, etc will go
export ZITI_HOME=~/.ziti/quickstart/${network_name}
export ZITI_POSTGRES_HOST="localhost"
export ZITI_NETWORK=${network_name}
export ZITI_USER="admin"
export ZITI_PWD="admin"
export ZITI_DOMAIN_SUFFIX=".ziti.netfoundry.io"
export ZITI_DOMAIN_SUFFIX=""
export ZITI_ID="${ZITI_HOME}/identities.yml"
export ZITI_FAB_MGMT_PORT="10000"
export ZITI_FAB_CTRL_PORT="6262"
export ZITI_CONTROLLER_NAME="${ZITI_NETWORK}-controller"
export ZITI_EDGE_NAME="${ZITI_NETWORK}-edge-controller"
export ZITI_EDGE_PORT="1280"
export ZITI_EDGE_ROUTER_NAME="${ZITI_NETWORK}-edge-router"
export ZITI_ROUTER_BR_NAME="${ZITI_NETWORK}-fabric-router-br"
export ZITI_ROUTER_BLUE_NAME="${ZITI_NETWORK}-fabric-router-blue"
export ZITI_ROUTER_RED_NAME="${ZITI_NETWORK}-fabric-router-red"

export ZITI_PKI="${ZITI_HOME}/pki"
export ZITI_CONTROLLER_HOSTNAME="${ZITI_CONTROLLER_NAME}${ZITI_DOMAIN_SUFFIX}"
export ZITI_EDGE_HOSTNAME="${ZITI_EDGE_NAME}${ZITI_DOMAIN_SUFFIX}"
export ZITI_EDGE_ROUTER_HOSTNAME="${ZITI_EDGE_ROUTER_NAME}${ZITI_DOMAIN_SUFFIX}"
export ZITI_SIGNING_CERT_NAME="${ZITI_NETWORK}-signing"
export ZITI_ROUTER_BR_HOSTNAME="${ZITI_ROUTER_BR_NAME}${ZITI_DOMAIN_SUFFIX}"
export ZITI_ROUTER_BLUE_HOSTNAME="${ZITI_ROUTER_BLUE_NAME}${ZITI_DOMAIN_SUFFIX}"
export ZITI_ROUTER_RED_HOSTNAME="${ZITI_ROUTER_RED_NAME}${ZITI_DOMAIN_SUFFIX}"

export ZITI_EDGE_API_HOSTNAME="${ZITI_EDGE_HOSTNAME}:${ZITI_EDGE_PORT}"

export ZITI_CONTROLLER_ROOTCA_NAME="${ZITI_CONTROLLER_HOSTNAME}-root-ca"
export ZITI_EDGE_ROOTCA_NAME="${ZITI_EDGE_HOSTNAME}-root-ca"
export ZITI_SIGNING_ROOTCA_NAME="${ZITI_SIGNING_CERT_NAME}-root-ca"

export ZITI_CONTROLLER_INTERMEDIATE_NAME="${ZITI_CONTROLLER_HOSTNAME}-intermediate"
export ZITI_EDGE_INTERMEDIATE_NAME="${ZITI_EDGE_HOSTNAME}-intermediate"
export ZITI_SIGNING_INTERMEDIATE_NAME="${ZITI_SIGNING_CERT_NAME}-intermediate"
export ZITI_SIGNING_SPURIOUS_NAME="${ZITI_SIGNING_INTERMEDIATE_NAME}_spurious_intermediate"

mkdir -p ${ZITI_HOME}/db
mkdir -p ${ZITI_PKI}

for zEnvVar in $(env | grep ZITI_ | sort); do echo "export "${zEnvVar}"" >> ${ZITI_HOME}/env; done
echo "export PFXLOG_NO_JSON=true" >> ${ZITI_HOME}/env