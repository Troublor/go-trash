#!/usr/bin/env bash

USER_HOME=$(eval echo ~${USER})

BIN_PATH="/usr/local/bin"
CMD_NAME="gotrash"
GOTRASH_PATH="/etc/gotrash"
CONFIG_NAME="gotrash-config.json"

if [[ ! -d ${GOTRASH_PATH} ]]; then
    sudo mkdir -p ${GOTRASH_PATH}
fi
if [[ -f ${BIN_PATH}/${CMD_NAME} ]]; then
    sudo rm ${BIN_PATH}/${CMD_NAME}
fi
if [[ -f ${BIN_PATH}/${CONFIG_NAME} ]]; then
    sudo rm ${BIN_PATH}/${CONFIG_NAME}
fi

echo "{\"trashDir\":\"${GOTRASH_PATH}\"}" > ${GOTRASH_PATH}/${CONFIG_NAME}
go build -o ${GOTRASH_PATH}/${CMD_NAME}
sudo ln -s ${GOTRASH_PATH}/${CMD_NAME} ${BIN_PATH}/${CMD_NAME}
sudo chmod 777 ${BIN_PATH}/${CMD_NAME}
sudo ln -s ${GOTRASH_PATH}/${CONFIG_NAME} ${BIN_PATH}/${CONFIG_NAME}
sudo chmod 777 ${BIN_PATH}/${CONFIG_NAME}
sudo chmod 777 -R ${GOTRASH_PATH}

