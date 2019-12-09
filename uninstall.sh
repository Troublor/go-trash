#!/usr/bin/env bash

CMD_NAME="gotrash"
CONFIG_NAME="gotrash-config.json"

BIN_LOCATION=$(which ${CMD_NAME})
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