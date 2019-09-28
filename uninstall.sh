#!/usr/bin/env bash

BIN_PATH="/usr/local/bin"
CMD_NAME="gotrash"
GOTRASH_PATH="/etc/gotrash"
CONFIG_NAME="gotrash-config.json"

sudo rm ${BIN_PATH}/${CMD_NAME}
sudo rm ${BIN_PATH}/${CONFIG_NAME}
sudo rm -r -d ${GOTRASH_PATH}

