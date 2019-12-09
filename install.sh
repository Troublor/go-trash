#!/usr/bin/env bash

USER_HOME=$(eval echo ~${USER})

CMD_NAME="gotrash"
CONFIG_NAME="gotrash-config.json"
GOTRASH_PATH=""
BIN_PATH=""

GLOBAL=1

while getopts "e:d:l" arg
do
    case $arg in
        l)
        GLOBAL=0
        ;;

        e)
        BIN_PATH=$OPTARG
        ;;

        d)
        GOTRASH_PATH=$OPTARG
        ;;

        \?)
        echo "Invalid option: -$OPTARG"
        ;;
    esac
done

if [[ ${GLOBAL} -eq 1 ]]; then
    if [[ ${GOTRASH_PATH} = "" ]]; then
        GOTRASH_PATH="/etc/gotrash"
    fi
    if [[ ${BIN_PATH} = "" ]]; then
        BIN_PATH="/usr/local/bin"
    fi
else
    if [[ ${GOTRASH_PATH} = "" ]]; then
        GOTRASH_PATH="${USER_HOME}/.gotrash"
    fi
    if [[ ${BIN_PATH} = "" ]]; then
        BIN_PATH="${USER_HOME}/bin"
    fi
fi

# check GOTRASH_PATH
if [[ -d ${GOTRASH_PATH} ]]; then
    echo "[WARN] ${GOTRASH_PATH} exists"
elif [[ -f ${GOTRASH_PATH} ]]; then
    echo "[ERROR] ${GOTRASH_PATH} is a file"
    exit
else
    if [[ ${GLOBAL} -eq 1 ]]; then
        sudo mkdir -p ${GOTRASH_PATH}
    else
        mkdir -p ${BIN_PATH}
    fi
fi

# set GOTRASH_PATH privilege for global installation
if [[ ${GLOBAL} -eq 1 ]]; then
    sudo chmod 777 -R ${GOTRASH_PATH}
fi

# check BIN_PATH
if [[ -d ${BIN_PATH} ]]; then
    echo "[WARN] ${BIN_PATH} exists"
elif [[ -f ${BIN_PATH} ]]; then
    echo "[ERROR] ${BIN_PATH} is a file"
    exit
else
    if [[ ${GLOBAL} -eq 1 ]]; then
        sudo mkdir -p ${BIN_PATH}
    else
        mkdir -p ${BIN_PATH}
    fi
fi

# remove previous installation (if any)
if [[ ${GLOBAL} -eq 1 ]]; then
    if [[ -L ${BIN_PATH}/${CMD_NAME} ]]; then
        sudo rm ${BIN_PATH}/${CMD_NAME}
    fi
    if [[ -L ${BIN_PATH}/${CONFIG_NAME} ]]; then
        sudo rm ${BIN_PATH}/${CONFIG_NAME}
    fi
else
    if [[ -L ${BIN_PATH}/${CMD_NAME} ]]; then
        rm ${BIN_PATH}/${CMD_NAME}
    fi
    if [[ -L ${BIN_PATH}/${CONFIG_NAME} ]]; then
        rm ${BIN_PATH}/${CONFIG_NAME}
    fi
fi

# generate config file
# TODO need refinement
if [[ ${GLOBAL} -eq 1 ]]; then
    sudo echo "{\"trashDir\":\"${GOTRASH_PATH}\"}" > ${GOTRASH_PATH}/${CONFIG_NAME}
else
    echo "{\"trashDir\":\"${GOTRASH_PATH}\"}" > ${GOTRASH_PATH}/${CONFIG_NAME}
fi

# build
echo "[INFO] Start building gotrash..."
go build -o ${GOTRASH_PATH}/${CMD_NAME}
echo "[INFO] Building finished"

if [[ ${GLOBAL} -eq 1 ]]; then
    sudo ln -s ${GOTRASH_PATH}/${CMD_NAME} ${BIN_PATH}/${CMD_NAME}
    sudo ln -s ${GOTRASH_PATH}/${CONFIG_NAME} ${BIN_PATH}/${CONFIG_NAME}
    sudo chmod 666 -R ${GOTRASH_PATH}
    sudo chmod +x ${GOTRASH_PATH}
    sudo chmod 777 ${BIN_PATH}/${CMD_NAME}
    sudo chmod 666 ${BIN_PATH}/${CONFIG_NAME}
else
    ln -s ${GOTRASH_PATH}/${CMD_NAME} ${BIN_PATH}/${CMD_NAME}
    ln -s ${GOTRASH_PATH}/${CONFIG_NAME} ${BIN_PATH}/${CONFIG_NAME}
    chmod +x ${BIN_PATH}/${CMD_NAME}
    chmod -x ${BIN_PATH}/${CONFIG_NAME}
fi

if [[ ${GLOBAL} -eq 1 ]]; then
    echo "[SUCCESS] Global installation finished"
    echo "[INFO] ${CMD_NAME} is available for all users"
else
    echo "[SUCCESS] Local installation finished"
    echo "[INFO] ${CMD_NAME} is only available for current user"
fi
echo "[INFO] Trash database is located at ${GOTRASH_PATH}"
echo "[INFO] ${CMD_NAME} command has been compiled to ${GOTRASH_PATH} and a soft link to ${CMD_NAME} has been created at ${BIN_PATH}"