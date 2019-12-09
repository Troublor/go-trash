#!/usr/bin/env bash

CMD_NAME="gotrash"
BIN_LOCATION=$(which ${CMD_NAME})
CMD_LOCATION=$(readlink ${BIN_LOCATION})
GOTRASH_PATH=$(dirname ${CMD_LOCATION})
owner=$(stat -c '%U' ${GOTRASH_PATH})
if [[ ${owner} -eq ${USER} ]]; then
    sudo rm ${BIN_LOCATION}
    sudo rm -r -d ${GOTRASH_PATH}
else
    rm ${BIN_LOCATION}
    rm -r -d ${GOTRASH_PATH}
fi