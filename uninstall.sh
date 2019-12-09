#!/usr/bin/env bash

CMD_NAME="gotrash"
CONFIG_NAME="config-gotrash.json"

BIN_LOCATION=$(command -v ${CMD_NAME})
if [[ ! -x ${BIN_LOCATION} ]]; then
    echo "command ${CMD_NAME} not found"
    exit
fi
BIN_PATH=$(dirname ${BIN_LOCATION})
CMD_LOCATION=$(readlink ${BIN_LOCATION})
GOTRASH_PATH=$(dirname ${CMD_LOCATION})
owner=$(stat -c '%U' ${GOTRASH_PATH})
if [[ ${owner} = ${USER} ]]; then
    rm ${BIN_LOCATION}
    rm ${BIN_PATH}/${CONFIG_NAME}
    rm -r -d ${GOTRASH_PATH}
else
    sudo rm ${BIN_LOCATION}
    sudo rm ${BIN_PATH}/${CONFIG_NAME}
    sudo rm -r -d ${GOTRASH_PATH}
fi