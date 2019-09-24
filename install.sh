#!/usr/bin/env bash

USER_HOME=$(eval echo ~${USER})

INSTALL_PATH="/usr/local/bin"
CMD_NAME="gotrash"

TRASH_DIR="${USER_HOME}/.gotrash"

echo "{\"trashDir\":\"${TRASH_DIR}\"}" > /tmp/gotrash-config.json
go build -o /tmp/${CMD_NAME}
sudo mv /tmp/${CMD_NAME} ${INSTALL_PATH}/${CMD_NAME}
sudo mv /tmp/gotrash-config.json ${INSTALL_PATH}/gotrash-config.json

